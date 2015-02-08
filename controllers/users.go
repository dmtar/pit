package controllers

import (
	"errors"
	"net/http"

	"github.com/dmtar/pit/models"
	"github.com/zenazn/goji/web"
	gojiMiddleware "github.com/zenazn/goji/web/middleware"
	"gopkg.in/mgo.v2/bson"
)

var Users = NewUsersController()

type UsersController struct {
	BaseController
	M *models.UserModel
}

func NewUsersController() *UsersController {
	return &UsersController{
		M: models.User,
	}
}

func (controller *UsersController) Routes() (root *web.Mux) {
	root = web.New()
	root.Use(gojiMiddleware.SubRouter)
	root.Post("/", Users.New)
	root.Put("/:objectid", Users.Edit)
	root.Get("/logout", Users.Logout)
	root.Get("/me", Users.CurrentUser)
	root.Post("/auth", Users.Auth)
	root.Get("/search/username/:username", Users.SearchByUsername)
	root.Get("/:objectId", Users.Find)
	root.Get("/:objectId/albums", Users.GetAlbums)
	return
}

func (controller *UsersController) CurrentUser(c web.C, w http.ResponseWriter, r *http.Request) {
	if user := controller.GetCurrentUser(c); user == nil {
		controller.Error(w, errors.New("You are not logged in!"))
	} else {
		controller.Write(w, user)
	}
}
func (controller *UsersController) Find(c web.C, w http.ResponseWriter, r *http.Request) {
	if user, err := controller.M.Find(c.URLParams["objectId"]); err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, user)
	}
}

func (controller *UsersController) GetAlbums(c web.C, w http.ResponseWriter, r *http.Request) {
	currentUser := controller.GetCurrentUser(c)
	objectId := c.URLParams["objectId"]
	public := currentUser == nil || currentUser.Id != bson.ObjectIdHex(objectId)
	if albums, err := controller.M.GetAlbums(objectId, public); err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, albums)
	}
}

func (controller *UsersController) Auth(c web.C, w http.ResponseWriter, r *http.Request) {
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
		session.Values["UserId"] = user.Id.Hex()
		session.Save(r, w)
		controller.Write(w, user)
	}
}

func (controller *UsersController) SearchByUsername(c web.C, w http.ResponseWriter, r *http.Request) {
	if users, err := controller.M.SearchByUsername(c.URLParams["username"]); err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, users)
	}
}

func (controller *UsersController) Logout(c web.C, w http.ResponseWriter, r *http.Request) {
	session := controller.GetSession(c)
	currentUser := controller.GetCurrentUser(c)

	if currentUser == nil {
		controller.Error(w, errors.New("You are not logged it!"))
		return
	}

	session.Values["UserId"] = nil
	session.Save(r, w)
}

func (controller *UsersController) New(c web.C, w http.ResponseWriter, r *http.Request) {
	params := controller.GetParams(c)
	session := controller.GetSession(c)
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
		session.Values["UserId"] = user.Id.Hex()
		session.Save(r, w)
		controller.Write(w, user)
	}
}

func (controller *UsersController) Edit(c web.C, w http.ResponseWriter, r *http.Request) {
	currentUser := controller.GetCurrentUser(c)

	if currentUser == nil {
		controller.Error(w, errors.New("You are not logged it!"))
		return
	}
	params := controller.GetParams(c)
	params.Add("user", currentUser)

	if user, err := controller.M.Edit(params); err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, user)
	}
}
