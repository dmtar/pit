package controllers

import (
	"errors"
	"net/http"

	"github.com/dmtar/pit/models"
	"github.com/dmtar/pit/system"
	"github.com/zenazn/goji/web"
	gojiMiddleware "github.com/zenazn/goji/web/middleware"
)

var Albums = NewAlbumsController()

type AlbumsController struct {
	BaseController
	M *models.AlbumModel
}

func NewAlbumsController() *AlbumsController {
	return &AlbumsController{
		M: models.Album,
	}
}

func (controller *AlbumsController) Routes() (root *web.Mux) {
	root = web.New()
	root.Use(gojiMiddleware.SubRouter)
	root.Get("/", Albums.FindByUser)
	root.Get("/public", Albums.Public)
	root.Get("/:objectId", Albums.Find)
	root.Get("/:objectId/pictures", Albums.GetPictures)
	root.Delete("/:objectId", Albums.Remove)
	root.Post("/", Albums.New)
	root.Put("/:objectId", Albums.Edit)
	return
}

func (controller *AlbumsController) Find(c web.C, w http.ResponseWriter, r *http.Request) {
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

func (controller *AlbumsController) GetPictures(c web.C, w http.ResponseWriter, r *http.Request) {
	if album, err := controller.M.Find(c.URLParams["objectId"]); err != nil {
		controller.Error(w, err)
	} else {
		currentUser := controller.GetCurrentUser(c)
		if album.Public || currentUser != nil && currentUser.Id == album.User {
			if pictures, err := models.Picture.FindByAlbum(album.Id.Hex()); err != nil {
				controller.Error(w, err)
			} else {
				controller.Write(w, pictures)
			}
		} else {
			controller.Error(w, errors.New("This album is private!"))
		}
	}
}

func (controller *AlbumsController) FindByUser(c web.C, w http.ResponseWriter, r *http.Request) {
	currentUser := controller.GetCurrentUser(c)

	if currentUser == nil {
		controller.Error(w, errors.New("You must be logged in to get your albums!"))
		return
	}

	albums, err := controller.M.FindByUser(system.Params{
		"user": currentUser,
	})

	if err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, albums)
	}
}

func (controller *AlbumsController) Public(c web.C, w http.ResponseWriter, r *http.Request) {
	currentUser := controller.GetCurrentUser(c)

	if currentUser == nil {
		controller.Error(w, errors.New("You must be logged in to view albums!"))
		return
	}

	if albums, err := controller.M.Public(); err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, albums)
	}
}

func (controller *AlbumsController) New(c web.C, w http.ResponseWriter, r *http.Request) {
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

func (controller *AlbumsController) Edit(c web.C, w http.ResponseWriter, r *http.Request) {
	params := controller.GetParams(c)
	currentUser := controller.GetCurrentUser(c)
	requiredParams := []string{"name"}
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

func (controller *AlbumsController) Remove(c web.C, w http.ResponseWriter, r *http.Request) {
	currentUser := controller.GetCurrentUser(c)
	if currentUser != nil {
		if err := controller.M.Remove(c.URLParams["objectId"]); err != nil {
			controller.Error(w, errors.New("Sorry something went wrong. The album cannot be removed right now!"))
			return
		} else {
			controller.Write(w, "The album is removed!")
			return
		}
	}
}
