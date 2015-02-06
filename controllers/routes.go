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
	root.Handle("/user/*", User.Routes())
	root.Handle("/album/*", Album.Routes())
	return
}
