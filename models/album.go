package models

import (
	"errors"
	"fmt"
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
	album = NewAlbumData()
	err = model.MgoFind(objectId, album)

	return
}

func (model *AlbumModel) Create(params system.Params) (album *AlbumData, err error) {
	fmt.Println("creating album")
	err = model.Connect()

	if err != nil {
		return nil, err
	}

	user, ok := params.GetI("user").(*UserData)
	if !ok || user == nil {
		return nil, errors.New("We are missing a user here!")
	}

	dateRange := DateRange{
		Start: ParseDate(params.GetP("date_range").Get("start")),
		End:   ParseDate(params.GetP("date_range").Get("end")),
	}

	tags := tagit.NewTags(params.GetAString("tags")...)

	location := Location{
		Longitude: ParseFloat64(params.GetP("location").Get("lng")),
		Latitude:  ParseFloat64(params.GetP("location").Get("lat")),
		Name:      params.GetP("location").Get("name"),
	}

	existing, err := model.FindByUserAndFilters(system.Params{
		"tags":      tags,
		"dateRange": dateRange,
		"user":      user,
		"location":  location,
	})

	if existing != nil && existing.User == user.Id {
		return nil, fmt.Errorf("You already have an album with tags: %s and within data_range: %s for location: %s", tags, dateRange, location)
	}

	album = &AlbumData{
		Id:        bson.NewObjectId(),
		Name:      params.Get("name"),
		Tags:      tags,
		Location:  location,
		DateRange: dateRange,
		Public:    ParseBool(params.Get("public")),
		NumPhotos: 0,
		User:      user.Id,
	}

	err = model.C.Insert(album)

	return
}

func (model *AlbumModel) GetForUser(params system.Params) (albums []*AlbumData, err error) {
	if err := model.Connect(); err != nil {
		return nil, err
	}

	albums = make([]*AlbumData, 0)

	user, ok := params.GetI("user").(*UserData)
	if !ok || user == nil {
		return nil, errors.New("We are missing a user here!")
	}

	query := bson.M{"user": user.Id}

	if ParseBool(params.Get("public")) {
		query["public"] = true
	}

	err = model.C.Find(query).Sort("-_id").All(&albums)

	return
}

func (model *AlbumModel) FindByUserAndFilters(params system.Params) (*AlbumData, error) {
	tags, ok := params.GetI("tags").(*tagit.Tags)
	user, ok := params.GetI("user").(*UserData)
	dateRange, ok := params.GetI("dateRange").(DateRange)
	location, ok := params.GetI("location").(Location)

	if !ok {
		return nil, errors.New("Wrong input parameters for FindByUserAndFIlters")
	}

	var err error
	album := NewAlbumData()

	if err := model.Connect(); err != nil {
		return nil, err
	}

	query := bson.M{
		"user":             user.Id,
		"tags":             bson.M{"$all": tags.All()},
		"location.name":    location.Name,
		"date_range.start": bson.M{"$gte": dateRange.Start},
		"date_range.end":   bson.M{"$lte": dateRange.End},
	}

	err = model.C.Find(query).One(album)

	return album, err
}

func (model *AlbumModel) Edit(params system.Params) (*AlbumData, error) {
	if err := model.Connect(); err != nil {
		return nil, err
	}

	album, ok := params.GetI("album").(*AlbumData)
	if !ok {
		return nil, errors.New("Missing album!")
	}

	album.Name = params.Get("Name")
	album.Public = ParseBool(params.Get("public"))

	err := model.C.UpdateId(album.Id, album)

	return album, err
}

func (model *AlbumModel) FindForPicture(picture *PictureMeta) *AlbumData {
	if err := model.Connect(); err != nil {
		return nil
	}

	album := NewAlbumData()

	query := bson.M{
		"user":             picture.User,
		"tags":             bson.M{"$all": picture.Tags.All()},
		"location.name":    picture.Location.Name,
		"date_range.start": bson.M{"$lte": picture.Date},
		"date_range.end":   bson.M{"$gte": picture.Date},
	}

	if err := model.C.Find(query).One(album); err != nil {
		return nil
	} else {
		return album
	}

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
