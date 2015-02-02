package controllers

import (
	"net/http"

	"github.com/dmtar/pit/lib"
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

func (ac *AlbumController) Routes() (root *web.Mux) {
	root = web.New()
	root.Use(gojiMiddleware.SubRouter)
	root.Put("/new", Album.New)
	root.Get("/:objectId", Album.Find)
	root.Post("/:objectId/edit", Album.Edit)
	return
}

func (ac *AlbumController) Find(c web.C, w http.ResponseWriter, r *http.Request) {
	album, err := ac.M.Find(c.URLParams["objectId"])
	if err != nil {
		ac.Error(w, err)
	} else {
		ac.Write(w, album)
	}
}

func (ac *AlbumController) New(c web.C, w http.ResponseWriter, r *http.Request) {
	params := c.Env["params"].(lib.Params)

	err := params.Required("name", "location", "tags", "date_range")

	if err != nil {
		ac.Error(w, err)
		return
	}

	album, err := ac.M.Create(params)
	if err != nil {
		ac.Error(w, err)
	} else {
		ac.Write(w, album)
	}
}

func (ac *AlbumController) Edit(c web.C, w http.ResponseWriter, r *http.Request) {
	album, err := ac.M.Find(c.URLParams["objectId"])

	if err != nil {
		ac.Error(w, err)
		return
	}

	album, err = ac.M.Edit(album, c.Env["params"].(lib.Params))
	if err != nil {
		ac.Error(w, err)
	} else {
		ac.Write(w, album)
	}
}
