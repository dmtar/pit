package controllers

import (
	"net/http"

	"github.com/dmtar/pit/models"
	"github.com/zenazn/goji/web"
	gojiMiddleware "github.com/zenazn/goji/web/middleware"
)

var User = NewUserController()

type UserController struct {
	ApplicationController
	Model *models.UserModel
}

func NewUserController() *UserController {
	return &UserController{
		Model: models.User,
	}
}

func (uc *UserController) Routes() (root *web.Mux) {
	root = web.New()
	root.Use(gojiMiddleware.SubRouter)
	root.Get("/:objectId", User.Find)
	root.Get("/search/username/:username", User.SearchByUsername)
	return
}

func (uc *UserController) Find(c web.C, w http.ResponseWriter, r *http.Request) {
	user, err := uc.Model.Find(c.URLParams["objectId"])
	if err != nil {
		uc.Error(w, err)
	} else {
		uc.Write(w, user)
	}
}

func (uc *UserController) SearchByUsername(c web.C, w http.ResponseWriter, r *http.Request) {
	users, err := uc.Model.SearchByUsername(c.URLParams["username"])
	if err != nil {
		uc.Error(w, err)
	} else {
		uc.Write(w, users)
	}
}
