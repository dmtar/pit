package controllers

import (
	"github.com/dmtar/pit/system"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

func Root() (root *web.Mux) {
	root = web.New()
	root.Use(middleware.SubRouter)
	root.Use(system.JSON)
	root.Handle("/users/*", Users.Routes())
	root.Handle("/albums/*", Albums.Routes())
	root.Handle("/notifications/*", Notifications.Routes())
	return
}
