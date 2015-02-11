package models

import (
	"errors"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/dmtar/pit/system"
	tagit "github.com/ndyakov/tagit/bson"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type PictureMeta struct {
	Name     string        `bson:"name" json:"name"`
	Tags     *tagit.Tags   `bson:"tags" json:"tags"`
	Location Location      `bson:"location" json:"location"`
	Date     time.Time     `bson:"date" json:"date"`
	Album    bson.ObjectId `bson:"album,omitempty" json:"album,omitempty"`
	User     bson.ObjectId `bson:"user" json:"user"`
}

func NewPictureMeta() *PictureMeta {
	pm := new(PictureMeta)
	pm.Tags = tagit.NewTags()
	return pm
}

type PictureModel struct {
	MgoModel
}

var Picture = NewPictureModel("fs")

func NewPictureModel(prefix string) *PictureModel {
	model := new(PictureModel)
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
	if err != nil {
		return nil, errors.New("Something went wrong!")
	}

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
	}

	file, err := model.Grid.Create(picture.Name)

	if err != nil {
		return nil, err
	}

	_, err = io.Copy(file, formFile)

	if err != nil {
		return nil, err
	}

	file.SetMeta(picture)
	err = file.Close()

	return
}

// func (model *PictureModel) CreateThumbnail(file *GridFile) {
// 	//TODO: Create a thumb for the picture, save it somewhere
// }
