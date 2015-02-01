package models

import (
	"fmt"

	"github.com/dmtar/pit/lib"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

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

func (um *UserModel) Create(params lib.Params) (user *UserData, err error) {
	err = um.Connect()

	if err != nil {
		return nil, err
	}

	email := params.Get("email")

	user, err = um.FindByEmail(email)
	if user.Email == email {
		err = fmt.Errorf("The email %s is already taken!", email)
		return nil, err
	}

	password, err := um.generatePassword(params.Get("password"))

	if err != nil {
		return nil, err
	}

	user = &UserData{
		Id:          bson.NewObjectId(),
		Username:    params.Get("username"),
		DisplayName: params.Get("display_name"),
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

	err = um.C.Find(bson.M{"username": username}).All(&users)
	return
}

func (um *UserModel) Edit(user *UserData, params lib.Params) (*UserData, error) {
	err := um.Connect()

	if err != nil {
		return nil, err
	}

	user.DisplayName = params.Get("display_name")
	err = um.C.UpdateId(user.Id, user)

	return user, err
}

func (um *UserModel) FindByEmail(email string) (user *UserData, err error) {
	user = new(UserData)
	err = um.Connect()

	if err != nil {
		return nil, err
	}

	err = um.C.Find(bson.M{"email": email}).Limit(1).One(&user)
	return
}

func (um *UserModel) generatePassword(password string) (hash string, err error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return password, err
	}

	hash = string(hashBytes)
	return hash, nil
}
