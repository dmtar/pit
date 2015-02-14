package models

import (
	"errors"

	"github.com/dmtar/pit/system"
	"gopkg.in/mgo.v2/bson"
)

type NotificationData struct {
	Id   bson.ObjectId `bson:"_id" json:"id"`
	Text string        `bson:"text" json:"text"`
	User bson.ObjectId `bson:"user" json:"user"`
}

func NewNotificationData() *NotificationData {
	return &NotificationData{}
}

type NotificationModel struct {
	MgoModel
}

var Notification = NewNotificationModel("notifications")

func NewNotificationModel(collection string) *NotificationModel {
	model := new(NotificationModel)
	model.SetCollectionName(collection)
	return model
}

func (model *NotificationModel) Find(objectId string) (notification *NotificationData, err error) {
	notification = NewNotificationData()
	err = model.MgoFind(objectId, notification)

	return
}

func (model *NotificationModel) Create(params system.Params) (notification *NotificationData, err error) {
	if err := model.Connect(); err != nil {
		return nil, err
	}

	user, ok := params.GetI("user").(*UserData)
	if !ok || user == nil {
		return nil, errors.New("We are missing a user here!")
	}

	notification = &NotificationData{
		Id:   bson.NewObjectId(),
		Text: params.Get("text"),
		User: user.Id,
	}

	err = model.C.Insert(notification)

	return
}

func (model *NotificationModel) Pop(params system.Params) (notification *NotificationData, err error) {
	if err := model.Connect(); err != nil {
		return nil, err
	}

	notification = &NotificationData{}

	user, ok := params.GetI("user").(*UserData)
	if !ok || user == nil {
		return nil, errors.New("We are missing a user here!")
	}

	err = model.C.Find(bson.M{
		"user": user.Id,
	}).Sort("_id").One(notification)

	if err == nil && notification.Id.Hex() != "" {
		model.C.RemoveId(notification.Id)
	}

	return
}

func (model *NotificationModel) Delete(objectId string) error {
	if err := model.Connect(); err != nil {
		return err
	}

	if !bson.IsObjectIdHex(objectId) {
		return errors.New("The provided objectId is not valid!")
	}

	return model.C.RemoveId(objectId)
}
