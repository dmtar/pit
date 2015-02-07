package controllers

import (
	"fmt"
	"errors"
	"net/http"

	"github.com/dmtar/pit/models"
	"github.com/dmtar/pit/system"
	"github.com/zenazn/goji/web"
	gojiMiddleware "github.com/zenazn/goji/web/middleware"
)

var Picture = NewPictureController()

type PictureController struct {
	BaseController
	M *models.PictureModel
}

func NewPictureController() *PictureController {
	return &PictureController{
		M: models.Picture,
	}
}

func (controller *PictureController) Routes() (root *web.Mux) {
	root = web.New()
	root.Use(gojiMiddleware.SubRouter)
	root.Post("/new", Picture.New)
	return
}

func (controller *PictureController) New(c web.C, w http.ResponseWriter, r *http.Request) {

	file, _, err := r.FormFile("picture")

	if err != nil {
        fmt.Fprintln(w, err)
        return
    }

	currentUser := controller.GetCurrentUser(c)
	if currentUser == nil {
		 controller.Error(w, errors.New("You must be logged in to upload picture!"))
		 return
	}

	defer file.Close()

	picture, err := controller.M.Create(system.Params{
		"name": r.FormValue("name"),
		"tags": "tag1, tag2",
		"user_id": currentUser.Id,
	}, file);

	if err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, picture)
	}
}