package main

import (
	"github.com/zenazn/goji"

	"github.com/dmtar/pit/controllers"
	"github.com/dmtar/pit/middlewares"
)

func init() {
	goji.Use(middlewares.JSON)
}

func main() {
	defer goji.Serve()
	goji.Handle("/*", controllers.Root())
}
