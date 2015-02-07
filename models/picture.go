package models

import (
	"fmt"
	"io"
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
	pd := new(PictureMeta)
	pd.Tags = tagit.NewTags()
	return pd
}

type PictureModel struct {
	MgoModel
}

var Picture = NewPictureModel("pit")

func NewPictureModel(prefix string) *PictureModel {
	model := new(PictureModel)
	model.SetGridFSPrefix(prefix)
	return model
}

func (model *PictureModel) Create(params system.Params, formFile multipart.File) (picture *PictureMeta, err error) {
	if err := model.Connect(); err != nil {
		return nil, err
	}

	picture = &PictureMeta{
		Name: params.Get("name"),
		Tags: tagit.NewTags(params.GetAString("tags")...),
		Location: Location{
			Longitude: ParseFloat64(params.GetP("location").Get("lng")),
			Latitude:  ParseFloat64(params.GetP("location").Get("lat")),
			Name:      params.GetP("location").Get("name"),
		},
		Date: ParseDate(params.GetP("date"))
	}

	// picture.Album = model.FindAlbumForPicture(picture)

	file, err := model.Grid.Create(picture.Name)
	
	if err != nil {
		return nil, err
	}

	n, err := io.Copy(file, formFile)

	fmt.Println(n)

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