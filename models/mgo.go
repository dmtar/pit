package models

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

type MgoModel struct {
	session *mgo.Session
	db      *mgo.Database
	C       *mgo.Collection
}

func (mm *MgoModel) Connect(collection string) (err error) {
	mm.session, err = mgo.Dial("localhost")
	if err != nil {
		return fmt.Errorf("Can't connect to mongo, go error %v\n", err)
	}
	mm.db = mm.session.DB("pit")
	mm.C = mm.db.C(collection)
	return nil
}

func (mm *MgoModel) Close() {
	mm.session.Close()
}
