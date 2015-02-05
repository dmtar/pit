package controllers

import (
	"net/http"

	"github.com/dmtar/pit/lib"
	"github.com/dmtar/pit/models"
	"github.com/zenazn/goji/web"
	gojiMiddleware "github.com/zenazn/goji/web/middleware"
)

var User = NewUserController()

type UserController struct {
	BaseController
	M *models.UserModel
}

func NewUserController() *UserController {
	return &UserController{
		M: models.User,
	}
}

func (uc *UserController) Routes() (root *web.Mux) {
	root = web.New()
	root.Use(gojiMiddleware.SubRouter)
	root.Put("/new", User.New)
	root.Get("/:objectId", User.Find)
	root.Post("/:objectId/edit", User.Edit)
	root.Get("/search/username/:username", User.SearchByUsername)
	return
}

func (uc *UserController) Find(c web.C, w http.ResponseWriter, r *http.Request) {
	if user, err := uc.M.Find(c.URLParams["objectId"]); err != nil {
		uc.Error(w, err)
	} else {
		uc.Write(w, user)
	}
}

func (uc *UserController) SearchByUsername(c web.C, w http.ResponseWriter, r *http.Request) {
	if users, err := uc.M.SearchByUsername(c.URLParams["username"]); err != nil {
		uc.Error(w, err)
	} else {
		uc.Write(w, users)
	}
}

func (uc *UserController) New(c web.C, w http.ResponseWriter, r *http.Request) {
	params := c.Env["params"].(lib.Params)
	requiredParams := []string{"email", "username", "display_name", "password"}

	if err := params.Required(requiredParams...); err != nil {
		uc.Error(w, err)
		return
	}

	if err := params.ShouldBeEmail("email"); err != nil {
		uc.Error(w, err)
		return
	}

	if user, err := uc.M.Create(params); err != nil {
		uc.Error(w, err)
	} else {
		uc.Write(w, user)
	}
}

func (uc *UserController) Edit(c web.C, w http.ResponseWriter, r *http.Request) {
	user, err := uc.M.Find(c.URLParams["objectId"])
	params := c.Env["params"].(lib.Params)

	if err != nil {
		uc.Error(w, err)
		return
	}

	user, err = uc.M.Edit(user, params)
	if err != nil {
		uc.Error(w, err)
	} else {
		uc.Write(w, user)
	}
}
