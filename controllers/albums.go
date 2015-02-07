package controllers

import (
	"errors"
	"net/http"

	"github.com/dmtar/pit/models"
	"github.com/zenazn/goji/web"
	gojiMiddleware "github.com/zenazn/goji/web/middleware"
)

var Album = NewAlbumController()

type AlbumController struct {
	BaseController
	M *models.AlbumModel
}

func NewAlbumController() *AlbumController {
	return &AlbumController{
		M: models.Album,
	}
}

func (controller *AlbumController) Routes() (root *web.Mux) {
	root = web.New()
	root.Use(gojiMiddleware.SubRouter)
	root.Put("/new", Album.New)
	root.Get("/:objectId", Album.Find)
	root.Post("/:objectId/edit", Album.Edit)
	return
}

func (controller *AlbumController) Find(c web.C, w http.ResponseWriter, r *http.Request) {
	currentUser := controller.GetCurrentUser(c)
	if album, err := controller.M.Find(c.URLParams["objectId"]); err != nil {
		controller.Error(w, err)
	} else {
		if album.Public || currentUser != nil && currentUser.Id == album.User {
			controller.Write(w, album)
		} else {
			controller.Error(w, errors.New("This album is private!"))
		}
	}
}

func (controller *AlbumController) New(c web.C, w http.ResponseWriter, r *http.Request) {
	params := controller.GetParams(c)
	currentUser := controller.GetCurrentUser(c)
	requiredParams := []string{
		"name",
		"location", "location.lat", "location.lng", "location.name",
		"tags",
		"date_range", "date_range.start", "date_range.end",
	}

	if currentUser == nil {
		controller.Error(w, errors.New("You must be logged in to create an album!"))
		return
	}

	if err := params.Required(requiredParams...); err != nil {
		controller.Error(w, err)
		return
	}

	params.Add("user", currentUser)

	if album, err := controller.M.Create(params); err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, album)
	}
}

func (controller *AlbumController) Edit(c web.C, w http.ResponseWriter, r *http.Request) {
	params := controller.GetParams(c)
	currentUser := controller.GetCurrentUser(c)
	requiredParams := []string{"name", "public"}

	if err := params.Required(requiredParams...); err != nil {
		controller.Error(w, err)
		return
	}

	if currentUser == nil {
		controller.Error(w, errors.New("You must be logged in to edit an album!"))
		return
	}

	album, err := controller.M.Find(c.URLParams["objectId"])

	if err != nil {
		controller.Error(w, err)
		return
	}

	if album.User != currentUser.Id {
		controller.Error(w, errors.New("This album is not owned by you!"))
		return
	}

	params.Add("user", currentUser)
	params.Add("album", album)

	if album, err = controller.M.Edit(params); err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, album)
	}
}
