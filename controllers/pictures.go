package controllers

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/dmtar/pit/models"
	"github.com/dmtar/pit/system"
	"github.com/zenazn/goji/web"
	gojiMiddleware "github.com/zenazn/goji/web/middleware"
)

var Pictures = NewPicturesController()

type PicturesController struct {
	BaseController
	M *models.PictureModel
}

func NewPicturesController() *PicturesController {
	return &PicturesController{
		M: models.Picture,
	}
}

func (controller *PicturesController) Routes() (root *web.Mux) {
	root = web.New()
	root.Use(gojiMiddleware.SubRouter)
	root.Post("/new", Pictures.New)
	root.Get("/:objectId", Pictures.Find)
	root.Get("/file/:objectId", Pictures.GetFile)
	root.Delete("/remove/:objectId", Pictures.Remove)
	return
}

func (controller *PicturesController) GetFile(c web.C, w http.ResponseWriter, r *http.Request) {
	if file, err := controller.M.GetFile(c.URLParams["objectId"]); err != nil {
		controller.Error(w, err)
	} else {
		w.Header().Set("Content-Type", file.ContentType())
		_, err := io.Copy(w, file)
		if err != nil {
			controller.Error(w, err)
		}
	}
}

func (controller *PicturesController) Find(c web.C, w http.ResponseWriter, r *http.Request) {
	if picture, err := controller.M.Find(c.URLParams["objectId"]); err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, picture)
	}
}

func (controller *PicturesController) New(c web.C, w http.ResponseWriter, r *http.Request) {
	//TODO: Check the uploaded file for size, validity, existence and stuff.
	currentUser := controller.GetCurrentUser(c)
	if currentUser == nil {
		controller.Error(w, errors.New("You must be logged in to upload picture!"))
		return
	}

	r.ParseMultipartForm(0)

	file, header, err := r.FormFile("picture")
	if err != nil {
		controller.Error(w, errors.New("Missing file!"))
		return
	}

	contentType := header.Header.Get("Content-Type")
	ct := strings.Split(contentType, "/")
	if ct[0] != "image" {
		controller.Error(w, fmt.Errorf("%s doesn't seem like an image!", contentType))
		return
	}

	defer file.Close()

	params := system.Params{}
	for k, v := range r.Form {
		params.Add(k, v[0])
	}

	part := multipart.Part{Header: header.Header}
	params.Add("fileName", part.FileName())
	params.Add("contentType", contentType)

	if err = params.Required(
		"name",
		"tags",
		"date",
		"location[name]",
		"fileName",
		"contentType",
	); err != nil {
		controller.Error(w, err)
		return
	}

	params.Add("user", currentUser)

	picture, err := controller.M.Create(params, file)

	if err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, picture)
	}
}

func (controller *PicturesController) Remove(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Println(c.URLParams["objectId"])
	if err := controller.M.Remove(c.URLParams["objectId"]); err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, err)
	}
}

func (controller *PicturesController) CheckPictureName(name string) interface{} {
	//TODO: Check the picture name for special chars, length and stuff.
	return name
}
