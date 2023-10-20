package controllers

import (
	"go-auth/src/controllers/auth"
	"go-auth/src/controllers/home"
	"go-auth/src/controllers/product"

	"github.com/ylfyt/meta/meta"
)

var appRoutes = [][]meta.EndPoint{
	home.Routes,
	auth.Routes,
	product.Routes,
}