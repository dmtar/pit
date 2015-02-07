package controllers

import (
	"fmt"
	"net/http"

	"github.com/dmtar/pit/models"
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
	params := controller.GetParams(c)
	params = system.Params{
		//TODO: Get the parameters from the post form and pass them to picture model
		r.FormValue("")
	}

	file, _, err := r.FormFile("picture")

	if err != nil {
        fmt.Fprintln(w, err)
        return
    }
	
	defer file.Close()

	if picture, err := controller.M.Create(params, file); err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, picture)
	}
}