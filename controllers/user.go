package controllers

import (
	"errors"
	"fmt"
	"net/http"

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

func (controller *UserController) Routes() (root *web.Mux) {
	root = web.New()
	root.Use(gojiMiddleware.SubRouter)
	root.Put("/new", User.New)
	root.Get("/logout", User.Logout)
	root.Post("/auth", User.Auth)
	root.Post("/:objectId/edit", User.Edit)
	root.Get("/search/username/:username", User.SearchByUsername)
	root.Get("/:objectId", User.Find)
	return
}

func (controller *UserController) Find(c web.C, w http.ResponseWriter, r *http.Request) {
	if user, err := controller.M.Find(c.URLParams["objectId"]); err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, user)
	}
}

func (controller *UserController) Auth(c web.C, w http.ResponseWriter, r *http.Request) {
	currentUser := controller.GetCurrentUser(c)
	session := controller.GetSession(c)

	if currentUser != nil {
		controller.Error(w, errors.New("You must log out first!"))
		return
	}

	params := controller.GetParams(c)

	requiredParams := []string{"email", "password"}

	if err := params.Required(requiredParams...); err != nil {
		controller.Error(w, err)
		return
	}

	if err := params.ShouldBeEmail("email"); err != nil {
		controller.Error(w, err)
		return
	}

	if user, err := controller.M.Auth(params); err != nil {
		controller.Error(w, err)
	} else {
		fmt.Println(session.Values)
		session.Values["UserId"] = user.Id.Hex()
		session.Save(r, w)
		fmt.Println(session.Values)
		controller.Write(w, user)
	}
}

func (controller *UserController) SearchByUsername(c web.C, w http.ResponseWriter, r *http.Request) {
	if users, err := controller.M.SearchByUsername(c.URLParams["username"]); err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, users)
	}
}

func (controller *UserController) Logout(c web.C, w http.ResponseWriter, r *http.Request) {
	session := controller.GetSession(c)
	currentUser := controller.GetCurrentUser(c)

	if currentUser == nil {
		controller.Error(w, errors.New("You are not logged it!"))
		return
	}

	session.Values["UserId"] = nil
	session.Save(r, w)
}

func (controller *UserController) New(c web.C, w http.ResponseWriter, r *http.Request) {
	params := controller.GetParams(c)
	requiredParams := []string{"email", "username", "display_name", "password"}

	if err := params.Required(requiredParams...); err != nil {
		controller.Error(w, err)
		return
	}

	if err := params.ShouldBeEmail("email"); err != nil {
		controller.Error(w, err)
		return
	}

	if user, err := controller.M.Create(params); err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, user)
	}
}

func (controller *UserController) Edit(c web.C, w http.ResponseWriter, r *http.Request) {
	user, err := controller.M.Find(c.URLParams["objectId"])
	params := controller.GetParams(c)

	if err != nil {
		controller.Error(w, err)
		return
	}

	user, err = controller.M.Edit(user, params)
	if err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, user)
	}
}
