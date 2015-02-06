package controllers

import (
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
	if album, err := controller.M.Find(c.URLParams["objectId"]); err != nil {
		controller.Error(w, err)
	} else {
		if album.Public {
			controller.Write(w, album)
		}
	}
}

func (controller *AlbumController) New(c web.C, w http.ResponseWriter, r *http.Request) {
	params := controller.GetParams(c)
	requiredParams := []string{
		"name",
		"location", "location.lat", "location.lng", "location.name",
		"tags",
		"date_range", "date_range.start", "date_range.end",
	}

	if err := params.Required(requiredParams...); err != nil {
		controller.Error(w, err)
		return
	}

	if album, err := controller.M.Create(params); err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, album)
	}
}

func (controller *AlbumController) Edit(c web.C, w http.ResponseWriter, r *http.Request) {
	params := controller.GetParams(c)
	album, err := controller.M.Find(c.URLParams["objectId"])

	if err != nil {
		controller.Error(w, err)
		return
	}

	if album, err = controller.M.Edit(album, params); err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, album)
	}
}