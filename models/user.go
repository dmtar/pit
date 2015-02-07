package models

import (
	"errors"
	"fmt"

	"github.com/dmtar/pit/system"
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

func NewUserData() *UserData {
	return &UserData{}
}

type UserModel struct {
	MgoModel
}

var User = NewUserModel("users")

func NewUserModel(collection string) *UserModel {
	model := new(UserModel)
	model.SetCollectionName(collection)
	return model
}

func (model *UserModel) Find(objectId string) (user *UserData, err error) {
	user = NewUserData()
	err = model.MgoFind(objectId, user)

	return
}

func (model *UserModel) Auth(params system.Params) (user *UserData, err error) {
	if err := model.Connect(); err != nil {
		return nil, err
	}

	email := params.Get("email")
	user, err = model.FindByEmail(email)

	if err == nil {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Get("password")))
		if err == nil {
			return user, nil
		}
	}

	return nil, errors.New("User not found or password incorrect!")
}

func (model *UserModel) Create(params system.Params) (user *UserData, err error) {
	if err := model.Connect(); err != nil {
		return nil, err
	}

	email := params.Get("email")

	user, err = model.FindByEmail(email)
	if user.Email == email {
		err = fmt.Errorf("The email %s is already taken!", email)
		return nil, err
	}

	password, err := model.generatePassword(params.Get("password"))

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

	err = model.C.Insert(user)

	return
}

func (model *UserModel) SearchByUsername(username string) (users []*UserData, err error) {
	if err := model.Connect(); err != nil {
		return nil, err
	}

	users = make([]*UserData, 0)
	err = model.C.Find(bson.M{"username": username}).All(&users)
	return
}

func (model *UserModel) Edit(params system.Params) (*UserData, error) {
	if err := model.Connect(); err != nil {
		return nil, err
	}

	user, ok := params.GetI("user").(*UserData)
	if !ok || user == nil {
		return nil, errors.New("We are missing a user here!")
	}

	user.DisplayName = params.Get("display_name")
	err := model.C.UpdateId(user.Id, user)

	return user, err
}

func (model *UserModel) FindByEmail(email string) (user *UserData, err error) {
	if err := model.Connect(); err != nil {
		return nil, err
	}

	user = NewUserData()

	err = model.C.Find(bson.M{"email": email}).Limit(1).One(&user)
	return
}

func (model *UserModel) GetAlbums(objectId string, public bool) (albums []*AlbumData, err error) {
	user, err := model.Find(objectId)
	if err != nil {
		return nil, err
	}

	return Album.GetForUser(system.Params{
		"user":   user,
		"public": public,
	})

}

func (model *UserModel) generatePassword(password string) (hash string, err error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return password, err
	}

	hash = string(hashBytes)
	return hash, nil
}
