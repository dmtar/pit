package models

import (
	"strconv"

	"github.com/dmtar/pit/system"
	tagit "github.com/ndyakov/tagit/bson"
	"gopkg.in/mgo.v2/bson"
)

type AlbumData struct {
	Id        bson.ObjectId `bson:"_id" json:"id"`
	Name      string        `bson:"name" json:"name"`
	Tags      *tagit.Tags   `bson:"tags" json:"tags"`
	Location  Location      `bson:"location" json:"location"`
	DateRange DateRange     `bson:"date_range" json:"date_range"`
	Public    bool          `bson:"public" json:"public"`
	NumPhotos int           `bson:"num_photos" json:"num_photos"`
	User      bson.ObjectId `bson:"user", json:"user"`
}

func NewAlbumData() *AlbumData {
	ad := new(AlbumData)
	ad.Tags = tagit.NewTags()
	return ad
}

type AlbumModel struct {
	MgoModel
}

var Album = NewAlbumModel("albums")

func NewAlbumModel(collection string) *AlbumModel {
	model := new(AlbumModel)
	model.SetCollectionName(collection)
	return model
}

func (model *AlbumModel) Find(objectId string) (album *AlbumData, err error) {
	album = new(AlbumData)
	err = model.MgoFind(objectId, album)

	return
}

func (model *AlbumModel) Create(params system.Params) (album *AlbumData, err error) {
	err = model.Connect()

	if err != nil {
		return nil, err
	}

	album = &AlbumData{
		Id:   bson.NewObjectId(),
		Name: params.Get("name"),
		Tags: tagit.NewTags(params.GetAString("tags")...),
		Location: Location{
			Longitude: ParseFloat64(params.GetP("location").Get("lng")),
			Latitude:  ParseFloat64(params.GetP("location").Get("lat")),
			Name:      params.GetP("location").Get("name"),
		},
		DateRange: DateRange{
			Start: ParseDate(params.GetP("date_range").Get("start")),
			End:   ParseDate(params.GetP("date_range").Get("end")),
		},
		Public:    ParseBool(params.Get("public")),
		NumPhotos: 0,
		User:      bson.NewObjectId(),
	}

	err = model.C.Insert(album)

	return
}

func (model *AlbumModel) Edit(album *AlbumData, params system.Params) (*AlbumData, error) {
	err := model.Connect()

	//if err != nil {
	//return nil, err
	//}

	//album.DisplayName = params.Get("display_name")
	//err = am.C.UpdateId(album.Id, album)

	return album, err
}

func ParseBool(input string) bool {
	if input == "" {
		return false
	}

	result, err := strconv.ParseBool(input)

	if err != nil {
		panic(err)
	}

	return result
}
