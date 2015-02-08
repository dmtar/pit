package models

import (
	"strings"
	"errors"
	"io"
	"time"
	"mime/multipart"
	"github.com/dmtar/pit/system"
	tagit "github.com/ndyakov/tagit/bson"
	"gopkg.in/mgo.v2/bson"
)

type PictureMeta struct {
	Name      string        	`bson:"name" json:"name"`
	Tags      *tagit.Tags   	`bson:"tags" json:"tags"`
	Location  Location      	`bson:"location" json:"location"`
	Date 	  time.Time     	`bson:"date" json:"date"`
	Album	  bson.ObjectId		`bson:"album" json:"album"`
	User	  bson.ObjectId		`bson:"user" json:"user"`
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

func (model *PictureModel) Find(objectId string) (picture *PictureMeta, err error) {
	//TODO: Find a way to provide the actual file also.
	picture = NewPictureMeta()
	err = model.MgoFind(objectId, picture)

	return
}

func (model *PictureModel) Remove(objectId string) (err error) {
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

	picture = &PictureMeta{
		Name: params.Get("name"),
		Tags: tagit.NewTags(strings.Split(params.Get("tags"), ",")...),
		Location: params.GetI("location").(Location),
		Date: ParseDate(params.Get("date")),
		User: bson.ObjectIdHex(params.Get("user_id")),
	}

	// pictureAlbum, err := model.FindAlbumForPicture(picture)

	if err != nil {
		return nil, err
	}

	// picture.Album = pictureAlbum.Id
	picture.Album = bson.NewObjectId()

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

	if err != nil {
		return nil, err
	}

	return
}

// func (model *PictureModel) CreateThumbnail(file *GridFile) {
// 	//TODO: Create a thumb for the picture, save it somewhere
// }

func (model *PictureModel) FindAlbumForPicture(picture *PictureMeta) (album *AlbumData, err error){
	//TODO: Find approriate album for the picture based on tags, date and stuff. If not insert it in the default album.

	if err := model.Connect(); err != nil {
		return nil, err
	}

	//TODO: Create the query for searching relevant album
	query := bson.M{}

	err = model.C.Find(query).One(album)

	if err != nil {
		return nil, err
	}

	//TODO: Check if there are no albums return the default one.
	return album, nil
}