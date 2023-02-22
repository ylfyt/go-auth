package controllers

import (
	"go-auth/src/controllers/auth"
	"go-auth/src/controllers/home"
	"go-auth/src/controllers/product"
	"go-auth/src/meta"
)

var appRoutes = [][]meta.EndPoint{
	home.Routes,
	auth.Routes,
	product.Routes,
}