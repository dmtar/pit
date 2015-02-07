package models

import (
	"errors"
	"fmt"

	"github.com/dmtar/pit/common"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MgoModel struct {
	session    *mgo.Session
	db         *mgo.Database
	C          *mgo.Collection
	Grid       *mgo.GridFS
	collection string
	gridPrefix string
}

func (model *MgoModel) SetCollectionName(collection string) {
	model.collection = collection
}

func (model *MgoModel) SetGridFSPrefix(prefix string) {
	model.gridPrefix = prefix
}

func (model *MgoModel) Connect() (err error) {
	if model.session != nil && model.session.Ping() == nil {
		return nil
	}

	model.session, err = mgo.Dial("localhost")
	if err != nil {
		return common.ServerError{
			fmt.Errorf("Can't connect to mongo, go error: %v", err),
		}
	}
	model.db = model.session.DB("pit")
	model.C = model.db.C(model.collection)
	model.Grid = model.db.GridFS(model.gridPrefix)
	return nil
}

func (model *MgoModel) Close() {
	model.session.Close()
}

func (model *MgoModel) MgoFind(objectId string, reciever interface{}) (err error) {
	if !bson.IsObjectIdHex(objectId) {
		return errors.New("The provided objectID is not valid!")
	}

	err = model.Connect()

	if err != nil {
		return err
	}

	err = model.C.Find(bson.M{"_id": bson.ObjectIdHex(objectId)}).One(reciever)

	return
}
