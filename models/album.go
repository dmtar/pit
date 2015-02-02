package models

import (
	"strconv"

	"github.com/dmtar/pit/lib"
	tagit "github.com/ndyakov/tagit/bson"
	"gopkg.in/mgo.v2/bson"
)

//albums = {
//_id: 1,
//name: 1,
//location: {
//lon: 1,
//lat: 1,
//name: 1
//},
//date_range: {
//start_date: 1,
//end_date: 1
//},
//public: 1,
//num_photos: 1,
//user: @user,
//tags: [a,b]
//}

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
	am := new(AlbumModel)
	am.SetCollectionName(collection)
	return am
}

func (am *AlbumModel) Find(objectId string) (album *AlbumData, err error) {
	album = new(AlbumData)
	err = am.MgoFind(objectId, album)

	return
}

func (am *AlbumModel) Create(params lib.Params) (album *AlbumData, err error) {
	err = am.Connect()

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

	err = am.C.Insert(album)

	return
}

func (am *AlbumModel) Edit(album *AlbumData, params lib.Params) (*AlbumData, error) {
	err := am.Connect()

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
