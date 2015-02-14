package system

import (
	"github.com/dmtar/pit/config"
	"github.com/gorilla/sessions"
	"gopkg.in/mgo.v2"
)

var App = NewApplication()

type Application struct {
	Store     *sessions.CookieStore
	DB        *mgo.Database
	dbSession *mgo.Session
}

func NewApplication() *Application {
	app := &Application{}
	app.Init()
	return app
}

func (application *Application) Init() {
	var err error
	secret := []byte(config.Secret)
	application.Store = sessions.NewCookieStore(secret)
	application.dbSession, err = mgo.Dial("localhost")
	check(err)
	application.DB = application.dbSession.DB("pit")
	application.DB.Login("pit", "pitter")
	go application.CreateIndexes()
}

func (application *Application) DBSession() *mgo.Session {
	var err error
	if application.dbSession.Ping() != nil {
		application.dbSession, err = mgo.Dial("localhost")
		application.DB = application.dbSession.DB("pit")
		application.DB.Login("pit", "pitter")
		check(err)
	}

	return application.dbSession
}
func (application *Application) Close() {
	application.dbSession.Close()
}

func (application *Application) CreateIndexes() {
	application.createIndexesOnPictures()
	application.createIndexesOnAlbums()
	application.createIndexesOnNotifications()
	application.createIndexesOnUsers()

}

func (application *Application) createIndexesOnPictures() {
	var err error
	collection := application.DB.C("fs.files")
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
	collection := application.DB.C("albums")
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

func (application *Application) createIndexesOnUsers() {
	var err error
	collection := application.DB.C("users")
	err = collection.EnsureIndex(mgo.Index{
		Key:        []string{"email"},
		Background: true,
		Unique:     true,
	})
	check(err)
}

func (application *Application) createIndexesOnNotifications() {
	var err error
	collection := application.DB.C("notifications")
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
