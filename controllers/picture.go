package controllers

import (
	"fmt"
	"strings"
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
	//TODO: Check the uploaded file for size, validity, existence and stuff.
	file, _, err := r.FormFile("picture")

	if err != nil {
        controller.Error(w, errors.New("Something is not ok with the uploaded file!"))
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
		return
	}

	defer file.Close()

	pictureLocation := models.Location{
		Longitude: models.ParseFloat64(r.FormValue("location_lng")),
		Latitude: models.ParseFloat64(r.FormValue("location_lat")),
		Name: r.FormValue("location_name"),
	}

	formValues := system.Params{}

	formValues.Add("name", pictureName)
	formValues.Add("tags", strings.TrimSpace(r.FormValue("tags")))
	formValues.Add("user_id", currentUser.Id.Hex())
	formValues.Add("location", pictureLocation)
	formValues.Add("date", "2015-02-08T10:35:33.648Z")

	picture, err := controller.M.Create(formValues, file);

	if err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, picture)
	}
}

func (controller *PictureController) Remove(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Println(c.URLParams["objectId"])
	if err := controller.M.Remove(c.URLParams["objectId"]); err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, err)
	}
}

func (controller *PictureController) CheckPictureName(name string) interface{} {
	//TODO: Check the picture name for special chars, length and stuff.
	return name
}