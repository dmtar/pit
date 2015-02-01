package main

import (
	"net/http"

	"github.com/zenazn/goji"

	"github.com/dmtar/pit/configuration"
	"github.com/dmtar/pit/controllers"
)

func main() {
	defer goji.Serve()
	goji.Get("/", http.FileServer(http.Dir(configuration.RootPath)))
	goji.Get("/index.html", http.FileServer(http.Dir(configuration.RootPath)))
	goji.Get("/assets/*", http.FileServer(http.Dir(configuration.RootPath)))
	goji.Handle("/*", controllers.Root())
}
