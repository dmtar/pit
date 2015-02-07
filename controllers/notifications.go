package controllers

import (
	"errors"
	"net/http"

	"github.com/dmtar/pit/models"
	"github.com/dmtar/pit/system"
	"github.com/zenazn/goji/web"
	gojiMiddleware "github.com/zenazn/goji/web/middleware"
)

var Notifications = NewNotificationsController()

type NotificationsController struct {
	BaseController
	M *models.NotificationModel
}

func NewNotificationsController() *NotificationsController {
	return &NotificationsController{
		M: models.Notification,
	}
}

func (controller *NotificationsController) Routes() (root *web.Mux) {
	root = web.New()
	root.Use(gojiMiddleware.SubRouter)
	root.Get("/", Notifications.GetForUser)
	root.Get("/:objectId", Notifications.Find)
	root.Delete("/:objectId", Notifications.Delete)
	return
}

func (controller *NotificationsController) Find(c web.C, w http.ResponseWriter, r *http.Request) {
	currentUser := controller.GetCurrentUser(c)
	if notification, err := controller.M.Find(c.URLParams["objectId"]); err != nil {
		controller.Error(w, err)
	} else {
		if currentUser != nil && currentUser.Id == notification.User {
			controller.Write(w, notification)
		} else {
			controller.Error(w, errors.New("This notification is not owned by you!"))
		}
	}
}

func (controller *NotificationsController) GetForUser(c web.C, w http.ResponseWriter, r *http.Request) {
	currentUser := controller.GetCurrentUser(c)

	if currentUser == nil {
		controller.Error(w, errors.New("You must be logged in to get your notifications!"))
		return
	}

	notifications, err := controller.M.GetForUser(system.Params{"user": currentUser})

	if err != nil {
		controller.Error(w, err)
	} else {
		controller.Write(w, notifications)
	}
}

func (controller *NotificationsController) Delete(c web.C, w http.ResponseWriter, r *http.Request) {
	var err error
	currentUser := controller.GetCurrentUser(c)
	objectId := c.URLParams["objectId"]

	if notification, err := controller.M.Find(objectId); err == nil {
		if currentUser != nil && currentUser.Id == notification.User {
			if err = controller.M.Delete(objectId); err == nil {
				controller.Write(w, notification)
				return
			}
		} else {
			err = errors.New("This notification is not owned by you!")
		}
	}

	controller.Error(w, err)
}
