package main

import (
	"fmt"
	"net/http"

	"github.com/zenazn/goji/web"

	"github.com/dmtar/pit/models"
)

func hello(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", c.URLParams["name"])
}

func main() {
	User := models.InitUserModel()

	// mongo test
	//user, err := User.New("pesho", "pesho", "pesho", "pesho")
	//fmt.Printf("\n%v, %v\n", user, err)

	// goji test
	//goji.Get("/hello/:name", hello)
	//goji.Get("asdasd")
	//goji.Serve()
}
