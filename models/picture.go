package models

import (
	"fmt"
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
	fmt.Println("Creating picture meta")
	picture = &PictureMeta{
		Name: params.Get("name"),
		Tags: tagit.NewTags(strings.Split(params.Get("tags"), ",")...),
		Location: params.GetI("location").(Location),
		Date: ParseDate(params.Get("date")),
		User: bson.ObjectIdHex(params.Get("user_id")),
	}

	pictureAlbum, err := model.FindAlbumForPicture(picture)

	if err != nil {
		return nil, err
	}

	picture.Album = pictureAlbum.Id

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

	q := bson.M{
		"tags":             bson.M{"$all": picture.Tags.All()},
		"date_range.start": bson.M{"$lte": picture.Date},
		"date_range.end":   bson.M{"$gte": picture.Date},
	}

	model.SetCollectionName("albums")
	query := model.C.Find(q)
	n, err := query.Count()

	if err != nil {
		return nil, errors.New("Something went wrong!")
	}

	if n == 0 {
		//TODO: Create default album
		return nil, errors.New("Cannot find album for this picture!")
	} else {
		err = query.One(&album)
	}

	if err != nil {
		return nil, errors.New("Cannot find album for this picture!")
	}

	return album, nil
}