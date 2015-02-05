package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dmtar/pit/common"
)

type BaseController struct{}

func (bc *BaseController) Write(w http.ResponseWriter, data interface{}) {
	result, err := json.Marshal(data)

	if err != nil {
		bc.Error(w, common.ServerError{err})
	} else {
		fmt.Fprint(w, string(result))
	}
}

func (bc *BaseController) Error(w http.ResponseWriter, err error) {

	code := 400
	if err.Error() == "not found" {
		code = 404
	} else if _, ok := err.(common.ServerError); ok {
		code = 500
	}

	w.WriteHeader(code)
	message, _ := json.Marshal(ErrorResponse{err.Error()})

	fmt.Fprintf(w, string(message))
}

func (bc *BaseController) NotFound(w http.ResponseWriter) {
	bc.Error(w, errors.New("not found"))
}

type ErrorResponse struct {
	Error string `json:"error"`
}
