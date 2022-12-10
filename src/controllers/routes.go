package controllers

import (
	"github.com/gorilla/mux"
	"go-auth/src/controllers/auth"
	"go-auth/src/controllers/home"
	"go-auth/src/middlewares"
)

var routes = []Route{
	{
		Name:        "Home",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: home.Home,
		Middlewares: []mux.MiddlewareFunc{
			middlewares.Authorization,
		},
	},
	{
		Name:        "Ping",
		Method:      "GET",
		Pattern:     "/ping",
		HandlerFunc: home.Ping,
	},
	{
		Name:        "AuthRegister",
		Method:      "POST",
		Pattern:     "/auth/register",
		HandlerFunc: auth.Register,
	},
	{
		Name:        "AuthLogin",
		Method:      "POST",
		Pattern:     "/auth/login",
		HandlerFunc: auth.Login,
	},
	{
		Name:        "AuthRefresh",
		Method:      "POST",
		Pattern:     "/auth/refresh",
		HandlerFunc: auth.RefreshToken,
	},
	{
		Name:        "AuthLogout",
		Method:      "POST",
		Pattern:     "/auth/logout",
		HandlerFunc: auth.Logout,
	},
}
