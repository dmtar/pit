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
	root.Get("/:objectId", Picture.Find)
	root.Delete("/remove/:objectId", Picture.Remove)
	return
}

func (controller *PictureController) Find(c web.C, w http.ResponseWriter, r *http.Request) {
	if user, err := controller.M.Find(c.URLParams["objectId"]); err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, user)
	}
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

	pictureName := controller.CheckPictureName(r.FormValue("name"))
	if pictureName == nil {
		controller.Error(w, errors.New("Not correct picture name!"))
	}

	defer file.Close()

	picture, err := controller.M.Create(system.Params{
		"name": pictureName,
		"tags": "tag1, tag2",
		"user_id": currentUser.Id,
	}, file);

	if err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, picture)
	}
}

func (controller *PictureController) Remove(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Println(c.URLParams["objectId"])
	if user, err := controller.M.Remove(c.URLParams["objectId"]); err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, user)
	}
}

func (controller *PictureController) CheckPictureName(name string) interface{} {
	//TODO: Check the picture name for special chars, length and stuff.
	return name
}