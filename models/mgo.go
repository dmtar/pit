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
	collection string
}

func (mm *MgoModel) SetCollectionName(collection string) {
	mm.collection = collection
}

func (mm *MgoModel) Connect() (err error) {
	if mm.session != nil && mm.session.Ping() == nil {
		return nil
	}

	mm.session, err = mgo.Dial("localhost")
	if err != nil {
		return common.ServerError{
			fmt.Errorf("Can't connect to mongo, go error: %v", err),
		}
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

	err = mm.C.Find(bson.M{"_id": bson.ObjectIdHex(objectId)}).One(reciever)

	return
}
