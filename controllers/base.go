package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dmtar/pit/common"
	"github.com/dmtar/pit/lib"
	"github.com/zenazn/goji/web"
)

type BaseController struct{}

func (controller *BaseController) Write(w http.ResponseWriter, data interface{}) {
	result, err := json.Marshal(data)

	if err != nil {
		controller.Error(w, common.ServerError{err})
	} else {
		fmt.Fprint(w, string(result))
	}
}

func (controller *BaseController) Error(w http.ResponseWriter, err error) {

	code := 400
	if err.Error() == "not found" {
		code = 404
	} else if _, ok := err.(common.ServerError); ok {
		code = 500
	}

	w.WriteHeader(code)
	message, _ := json.Marshal(common.ErrorResponse{err.Error()})

	fmt.Fprintf(w, string(message))
}

func (controller *BaseController) NotFound(w http.ResponseWriter) {
	controller.Error(w, errors.New("not found"))
}

func (controller *BaseController) GetParams(c web.C) (p lib.Params) {
	p, ok := c.Env["Params"].(lib.Params)

	if !ok {
		p = lib.Params{}
	}

	return
}
