package controllers

import (
	"go-auth/src/controllers/auth"
	"go-auth/src/controllers/home"
	"go-auth/src/controllers/product"

	"github.com/ylfyt/meta/meta"
)

func New() *meta.App {
	app := meta.New(&meta.Config{
		BaseUrl: "/api",
	})

	app.AddController(&home.HomeController{})
	app.AddController(&auth.AuthController{})
	app.AddController(&product.ProductController{})

	return app
}
