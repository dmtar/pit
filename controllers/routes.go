package controllers

import (
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

func Root() (root *web.Mux) {
	root = web.New()
	root.Use(middleware.SubRouter)
	root.Handle("/user/*", User.Routes())
	return
}
