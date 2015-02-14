package system

import (
	"fmt"

	"github.com/dmtar/pit/config"
	"github.com/gorilla/sessions"
	"gopkg.in/mgo.v2"
)

type Application struct {
	Store   *sessions.CookieStore
	db      *mgo.Database
	session *mgo.Session
}

func (application *Application) Init() {
	secret := []byte(config.Secret)
	application.Store = sessions.NewCookieStore(secret)
	go application.CreateIndexes()
}

func (application *Application) CreateIndexes() {
	var err error
	application.session, err = mgo.Dial("localhost")
	if err != nil {
		check(fmt.Errorf("Can't connect to mongo, go error: %v", err))
	}
	application.db = application.session.DB("pit")

	application.createIndexesOnPictures()
	application.createIndexesOnAlbums()
	application.createIndexesOnNotifications()

}

func (application *Application) createIndexesOnPictures() {
	var err error
	collection := application.db.C("fs.files")
	err = collection.EnsureIndex(mgo.Index{
		Key:        []string{"$2dsphere:metadata.location"},
		Background: true,
	})
	check(err)

	err = collection.EnsureIndex(mgo.Index{
		Key:        []string{"metadata.tags"},
		Background: true,
	})
	check(err)

	err = collection.EnsureIndex(mgo.Index{
		Key:        []string{"metadata.album"},
		Background: true,
	})
	check(err)

	err = collection.EnsureIndex(mgo.Index{
		Key:        []string{"metadata.user"},
		Background: true,
	})
	check(err)
}

func (application *Application) createIndexesOnAlbums() {
	var err error
	collection := application.db.C("albums")
	err = collection.EnsureIndex(mgo.Index{
		Key:        []string{"$2dsphere:location"},
		Background: true,
	})
	check(err)

	err = collection.EnsureIndex(mgo.Index{
		Key:        []string{"metadata.tags"},
		Background: true,
	})
	check(err)

	err = collection.EnsureIndex(mgo.Index{
		Key:        []string{"data_range.start"},
		Background: true,
	})
	check(err)

	err = collection.EnsureIndex(mgo.Index{
		Key:        []string{"data_range.end"},
		Background: true,
	})
	check(err)

	err = collection.EnsureIndex(mgo.Index{
		Key:        []string{"user"},
		Background: true,
	})
	check(err)
}

func (application *Application) createIndexesOnNotifications() {
	var err error
	collection := application.db.C("notifications")
	err = collection.EnsureIndex(mgo.Index{
		Key:        []string{"user"},
		Background: true,
	})
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
