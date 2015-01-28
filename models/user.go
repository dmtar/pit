package models

import (
	"errors"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id          bson.ObjectId `bson:"_id"`
	Username    string        `bson:"username"`
	DisplayName string        `bson:"display_name"`
	Email       string        `bson:"email"`
	Password    string        `bson:"password"`
}

type UserModel struct {
	MgoModel
}

func InitUserModel() *UserModel {
	return &UserModel{}
}

func (um *UserModel) Find(objectId string) (user *User, err error) {
	if !bson.IsObjectIdHex(objectId) {
		err = errors.New("The provided objectID is not valid!")
		return nil, err
	}

	err = um.Connect("users")

	if err != nil {
		return nil, err
	}

	defer um.Close()

	user = new(User)
	err = um.C.Find(bson.M{"_id": bson.ObjectIdHex(objectId)}).One(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (um *UserModel) New(username, displayname, email, password string) (user *User, err error) {
	err = um.Connect("users")

	if err != nil {
		return nil, err
	}

	defer um.Close()
	user = new(User)
	user.Id = bson.NewObjectId()
	user.Username = username
	user.DisplayName = displayname
	user.Email = email
	user.Password = password
	err = um.C.Insert(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}
