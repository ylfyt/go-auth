package controllers

import (
	"go-auth/src/controllers/auth"
	"go-auth/src/controllers/home"
	"go-auth/src/ctx"
)

var appRoutes = [][]ctx.Route{
	home.Routes,
	auth.Routes,
}