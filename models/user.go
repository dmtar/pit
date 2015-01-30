package models

import "gopkg.in/mgo.v2/bson"

type UserData struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	Username    string        `bson:"username" json:"username"`
	DisplayName string        `bson:"display_name" json:"display_name"`
	Email       string        `bson:"email" json:"email"`
	Password    string        `bson:"password" json:"-"`
}

type UserModel struct {
	MgoModel
}

var User = NewUserModel("users")

func NewUserModel(collection string) *UserModel {
	um := new(UserModel)
	um.SetCollectionName(collection)
	return um
}

func (um *UserModel) Find(objectId string) (user *UserData, err error) {
	user = new(UserData)
	err = um.MgoFind(objectId, user)

	if err != nil {
		return nil, err
	}

	return
}

func (um *UserModel) Create(username, displayname, email, password string) (user *UserData, err error) {
	err = um.Connect()

	if err != nil {
		return nil, err
	}

	defer um.Close()
	user = &UserData{
		Id:          bson.NewObjectId(),
		Username:    username,
		DisplayName: displayname,
		Email:       email,
		Password:    password,
	}

	err = um.C.Insert(user)

	return
}

func (um *UserModel) SearchByUsername(username string) (users []*UserData, err error) {
	users = make([]*UserData, 0)
	err = um.Connect()

	if err != nil {
		return nil, err
	}

	defer um.Close()

	err = um.C.Find(bson.M{"username": username}).All(&users)
	return
}
