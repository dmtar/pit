package models

import (
	"errors"
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MgoModel struct {
	session    *mgo.Session
	db         *mgo.Database
	C          *mgo.Collection
	collection string
}

func (mm *MgoModel) SetCollectionName(collection string) {
	mm.collection = collection
}

func (mm *MgoModel) Connect() (err error) {
	mm.session, err = mgo.Dial("localhost")
	if err != nil {
		return fmt.Errorf("Can't connect to mongo, go error %v\n", err)
	}
	mm.db = mm.session.DB("pit")
	mm.C = mm.db.C(mm.collection)
	return nil
}

func (mm *MgoModel) Close() {
	mm.session.Close()
}

func (mm *MgoModel) MgoFind(objectId string, reciever interface{}) (err error) {
	if !bson.IsObjectIdHex(objectId) {
		return errors.New("The provided objectID is not valid!")
	}

	err = mm.Connect()

	if err != nil {
		return err
	}

	defer mm.Close()

	err = mm.C.Find(bson.M{"_id": bson.ObjectIdHex(objectId)}).One(reciever)

	return
}
