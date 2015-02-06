package main

import (
	"net/http"

	"github.com/lidashuang/goji-gzip"
	"github.com/zenazn/goji"

	"github.com/dmtar/pit/config"
	"github.com/dmtar/pit/controllers"
	"github.com/dmtar/pit/system"
)

func init() {
	goji.Use(gzip.GzipHandler)
}

func main() {
	defer goji.Serve()
	app := system.Application{}
	app.Init()
	goji.Use(app.Session)
	HandleIndex()
	HandleAssets()
	goji.Handle("/*", controllers.Root())
}

func HandleIndex() {
	goji.Get("/", http.FileServer(http.Dir(config.RootPath)))
	goji.Get("/index.html", http.FileServer(http.Dir(config.RootPath)))
}

func HandleAssets() {
	goji.Get("/assets/*", http.FileServer(http.Dir(config.RootPath)))
	goji.Get("/js/*", http.FileServer(http.Dir(config.RootPath)))
	goji.Get("/tpl/*", http.FileServer(http.Dir(config.RootPath)))
}
