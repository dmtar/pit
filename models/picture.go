package models

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/dmtar/pit/system"
	tagit "github.com/ndyakov/tagit/bson"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type PictureFile struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	UploadDate  time.Time     `bson:"uploadDate" json:"upload_date"`
	Filename    string        `bson:"filename" json:"filename"`
	ContentType string        `bson:"contentType" json:"content_type"`
	Metadata    PictureMeta   `bson:"metadata" json:"metadata"`
}

type PictureMeta struct {
	Id       bson.ObjectId   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string          `bson:"name" json:"name"`
	Tags     *tagit.Tags     `bson:"tags" json:"tags"`
	Location Location        `bson:"location" json:"location"`
	Date     time.Time       `bson:"date" json:"date"`
	Album    bson.ObjectId   `bson:"album,omitempty" json:"album,omitempty"`
	User     bson.ObjectId   `bson:"user" json:"user"`
	Likes    []bson.ObjectId `bson:"likes" json:"likes"`
}

func NewPictureMeta() *PictureMeta {
	pm := new(PictureMeta)
	pm.Tags = tagit.NewTags()
	return pm
}

func (pm PictureMeta) CanBeViewedBy(user *UserData) error {
	public := false
	if pm.Album.Hex() != "" {
		if album, err := Album.Find(pm.Album.Hex()); err == nil {
			public = album.Public
		}
	}

	if public || user != nil && pm.User == user.Id {
		return nil
	}

	return errors.New("You don't have permissions to view this picture!")
}

func (pm PictureMeta) CanBeEditedBy(user *UserData) error {
	if user != nil && pm.User == user.Id {
		return nil
	}
	return errors.New("You don't have permissions to edit this picture!")
}

type PictureModel struct {
	MgoModel
}

var Picture = NewPictureModel("fs", "fs.files")

func NewPictureModel(prefix, collection string) *PictureModel {
	model := new(PictureModel)
	model.SetCollectionName(collection)
	model.SetGridFSPrefix(prefix)
	return model
}

func (model *PictureModel) GetFile(objectId string) (picture *mgo.GridFile, err error) {

	if err := model.Connect(); err != nil {
		return nil, err
	}

	if !bson.IsObjectIdHex(objectId) {
		return nil, errors.New("The provided objectID is not valid!")
	}

	return model.Grid.OpenId(bson.ObjectIdHex(objectId))
}

func (model *PictureModel) Find(objectId string) (picture *PictureMeta, err error) {
	if err := model.Connect(); err != nil {
		return nil, err
	}

	if !bson.IsObjectIdHex(objectId) {
		return nil, errors.New("The provided objectID is not valid!")
	}

	file, err := model.Grid.OpenId(bson.ObjectIdHex(objectId))
	if err != nil {
		return nil, errors.New("Cannot find picture with id: " + objectId)
	}

	err = file.GetMeta(&picture)
	picture.Id = bson.ObjectIdHex(objectId)

	return
}

func (model *PictureModel) Remove(objectId string) (err error) {
	if err := model.Connect(); err != nil {
		return err
	}

	if !bson.IsObjectIdHex(objectId) {
		return errors.New("The provided objectID is not valid!")
	}

	err = model.Grid.RemoveId(bson.ObjectIdHex(objectId))
	if err != nil {
		return errors.New("Something went wrong!")
	}

	return nil
}

func (model *PictureModel) Edit(params system.Params) (*PictureMeta, error) {
	if err := model.Connect(); err != nil {
		return nil, err
	}

	picture, ok := params.GetI("picture").(*PictureMeta)
	if !ok {
		return nil, errors.New("Missing picture!")
	}

	picture.Name = params.Get("name")

	err := model.C.UpdateId(picture.Id, bson.M{"$set": bson.M{"metadata": picture}})

	return picture, err
}

func (model *PictureModel) Create(params system.Params, formFile multipart.File) (picture *PictureMeta, err error) {
	if err := model.Connect(); err != nil {
		return nil, err
	}

	user := params.GetI("user").(*UserData)

	location := Location{
		Name:      params.Get("location[name]"),
		Longitude: ParseFloat64(params.Get("location[lng]")),
		Latitude:  ParseFloat64(params.Get("location[lat]")),
	}

	picture = &PictureMeta{
		Name:     params.Get("name"),
		Tags:     tagit.NewTags(strings.Split(params.Get("tags"), ",")...),
		Location: location,
		Date:     ParseDate(params.Get("date")),
		User:     user.Id,
	}

	album := Album.FindForPicture(picture)

	if album != nil {
		picture.Album = album.Id
		Album.IncreaseNumPictures(album.Id, 1)
	}

	file, err := model.Grid.Create(params.Get("fileName"))
	file.SetContentType(params.Get("contentType"))

	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(file, formFile); err != nil {
		return nil, err
	}

	file.SetMeta(picture)
	err = file.Close()

	if err == nil {
		picture.Id = file.Id().(bson.ObjectId)
	}

	return
}

func (model *PictureModel) FindByAlbum(albumId string) (pictures []*PictureMeta, err error) {
	if !bson.IsObjectIdHex(albumId) {
		return nil, errors.New("The provided objectID is not valid!")
	}

	if err := model.Connect(); err != nil {
		return nil, err
	}

	files := make([]*PictureFile, 0)
	pictures = make([]*PictureMeta, 0)

	err = model.Grid.Find(bson.M{"metadata.album": bson.ObjectIdHex(albumId)}).All(&files)

	for _, file := range files {
		file.Metadata.Id = file.Id
		pictures = append(pictures, &file.Metadata)
	}

	return pictures, err

}

func (model *PictureModel) FindByUser(objectId string) (pictures []*PictureMeta, err error) {
	if !bson.IsObjectIdHex(objectId) {
		return nil, errors.New("The provided objectID is not valid!")
	}

	if err := model.Connect(); err != nil {
		return nil, err
	}

	files := make([]*PictureFile, 0)
	pictures = make([]*PictureMeta, 0)

	err = model.Grid.Find(bson.M{"metadata.user": bson.ObjectIdHex(objectId)}).All(&files)

	for _, file := range files {
		file.Metadata.Id = file.Id
		pictures = append(pictures, &file.Metadata)
	}

	return pictures, err
}

func (model *PictureModel) Like(params system.Params) error {

	user := params.GetI("user").(*UserData)
	picture := params.GetI("picture").(*PictureMeta)

	query := bson.M{
		"_id":            picture.Id,
		"metadata.user":  bson.M{"$ne": user.Id},
		"metadata.likes": bson.M{"$nin": []bson.ObjectId{user.Id}},
	}

	err := model.C.Update(query, bson.M{
		"$addToSet": bson.M{"metadata.likes": user.Id},
	})

	if err == nil {
		Notification.Send(
			picture.User,
			fmt.Sprintf("%s, liked your picture %s!", user.DisplayName, picture.Name),
		)
	}

	return err
}

// func (model *PictureModel) CreateThumbnail(file *GridFile) {
// 	//TODO: Create a thumb for the picture, save it somewhere
// }
