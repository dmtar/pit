package main

import (
	"fmt"
	"net/http"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"

	"github.com/dmtar/pit/controllers"
	"github.com/dmtar/pit/middlewares"
)

func hello(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", c.URLParams["name"])
}

func init() {
	goji.Use(middlewares.JSON)
}

func main() {
	defer goji.Serve()
	goji.Handle("/*", controllers.Root())
}
