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
	"gopkg.in/mgo.v2/bson"
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
	root.Get("/", Pictures.FindByUser)
	root.Post("/", Pictures.New)
	root.Post("/like/:objectId", Pictures.Like)
	root.Get("/:objectId", Pictures.Find)
	root.Get("/canview/:objectId", Pictures.CanBeViewed)
	root.Get("/file/:objectId", Pictures.GetFile)
	root.Put("/:objectId", Pictures.Edit)
	root.Delete("/remove/:objectId", Pictures.Remove)
	return
}

func (controller *PicturesController) Like(c web.C, w http.ResponseWriter, r *http.Request) {
	currentUser := controller.GetCurrentUser(c)
	if currentUser == nil {
		controller.Error(w, errors.New("You must be logged in to like a picture!"))
		return
	}
	picture, err := controller.M.Find(c.URLParams["objectId"])
	if err != nil {
		controller.Error(w, err)
	}
	params := system.Params{"user": currentUser, "picture": picture}

	if err := controller.M.Like(params); err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, bson.M{"success": true})
	}

}

func (controller *PicturesController) FindByUser(c web.C, w http.ResponseWriter, r *http.Request) {
	currentUser := controller.GetCurrentUser(c)
	if currentUser == nil {
		controller.Error(w, errors.New("You must be logged in to view your picture!"))
		return
	}
	if pictures, err := controller.M.FindByUser(currentUser.Id.Hex()); err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, pictures)
	}
}

func (controller *PicturesController) CanBeViewed(c web.C, w http.ResponseWriter, r *http.Request) {
	currentUser := controller.GetCurrentUser(c)
	picture, err := controller.M.Find(c.URLParams["objectId"])
	if err != nil {
		controller.Error(w, err)
	} else {
		if err := picture.CanBeViewedBy(currentUser); err == nil {
			controller.Write(w, bson.M{"success": true})
		} else {
			controller.Error(w, err)
		}
	}
}

func (controller *PicturesController) GetFile(c web.C, w http.ResponseWriter, r *http.Request) {
	currentUser := controller.GetCurrentUser(c)
	picture, err := controller.M.Find(c.URLParams["objectId"])

	if err != nil {
		controller.Error(w, err)
		return
	}

	if err := picture.CanBeViewedBy(currentUser); err == nil {
		file, err := controller.M.GetFile(c.URLParams["objectId"])
		if err != nil {
			controller.Error(w, err)
			return
		}

		w.Header().Set("Content-Type", file.ContentType())
		_, err = io.Copy(w, file)

		if err != nil {
			controller.Error(w, err)
			return
		}

	} else {
		controller.Error(w, err)
	}
}

func (controller *PicturesController) Find(c web.C, w http.ResponseWriter, r *http.Request) {
	if picture, err := controller.M.Find(c.URLParams["objectId"]); err != nil {
		controller.Error(w, err)
	} else {
		currentUser := controller.GetCurrentUser(c)
		if err := picture.CanBeViewedBy(currentUser); err == nil {
			controller.Write(w, picture)
		} else {
			controller.Error(w, err)
		}
	}
}

func (controller *PicturesController) New(c web.C, w http.ResponseWriter, r *http.Request) {
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
	currentUser := controller.GetCurrentUser(c)
	picture, err := controller.M.Find(c.URLParams["objectId"])

	if err != nil {
		controller.Error(w, err)
		return
	}

	if err := picture.CanBeEditedBy(currentUser); err == nil {
		controller.Write(w, bson.M{"success": true})
	} else {
		controller.Error(w, err)
	}
	if err := controller.M.Remove(c.URLParams["objectId"]); err != nil {
		controller.Error(w, err)
	} else {
	}
}

func (controller *PicturesController) Edit(c web.C, w http.ResponseWriter, r *http.Request) {
	params := controller.GetParams(c)
	currentUser := controller.GetCurrentUser(c)
	requiredParams := []string{"name"}
	if err := params.Required(requiredParams...); err != nil {
		controller.Error(w, err)
		return
	}

	if currentUser == nil {
		controller.Error(w, errors.New("You must be logged in to edit picture!"))
		return
	}

	picture, err := controller.M.Find(c.URLParams["objectId"])

	if err != nil {
		controller.Error(w, err)
		return
	}

	if picture.User != currentUser.Id {
		controller.Error(w, errors.New("This picture is not yours!"))
		return
	}

	params.Add("user", currentUser)
	params.Add("picture", picture)

	if picture, err = controller.M.Edit(params); err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, picture)
	}
}

func (controller *PicturesController) CheckPictureName(name string) interface{} {
	//TODO: Check the picture name for special chars, length and stuff.
	return name
}
