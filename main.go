package main

import (
	"net/http"

	"github.com/lidashuang/goji-gzip"
	"github.com/zenazn/goji"

	"github.com/dmtar/pit/config"
	"github.com/dmtar/pit/controllers"
)

func init() {
	goji.Use(gzip.GzipHandler)
}

func main() {
	defer goji.Serve()
	goji.Get("/", http.FileServer(http.Dir(config.RootPath)))
	goji.Get("/index.html", http.FileServer(http.Dir(config.RootPath)))
	goji.Get("/assets/*", http.FileServer(http.Dir(config.RootPath)))
	goji.Get("/js/*", http.FileServer(http.Dir(config.RootPath)))
	goji.Handle("/*", controllers.Root())
}
