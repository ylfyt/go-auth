package controllers

import (
	"go-auth/src/controllers/auth"
	"go-auth/src/controllers/home"
	"go-auth/src/ctx"
)

var routes []ctx.Route = append(
	home.Routes,
	auth.Routes...
)
