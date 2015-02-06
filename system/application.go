package system

import (
	"github.com/dmtar/pit/config"
	"github.com/gorilla/sessions"
)

type Application struct {
	Store *sessions.CookieStore
}

func (application *Application) Init() {
	secret := []byte(config.Secret)
	application.Store = sessions.NewCookieStore(secret)
}
