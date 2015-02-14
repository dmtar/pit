package main

import (
	"net/http"
	"runtime"

	"github.com/lidashuang/goji-gzip"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/graceful"

	"github.com/dmtar/pit/config"
	"github.com/dmtar/pit/controllers"
	"github.com/dmtar/pit/system"
)

func init() {
	goji.Use(gzip.GzipHandler)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	defer goji.Serve()
	goji.Use(system.App.Session)
	HandleIndex()
	HandleAssets()
	graceful.PostHook(func() {
		system.App.Close()
	})

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
