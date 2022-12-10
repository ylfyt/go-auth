package home

import (
	"go-auth/src/ctx"
	"go-auth/src/middlewares"

	"github.com/gorilla/mux"
)

var Routes = []ctx.Route{
	{
		Name:        "Home",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: home,
		Middlewares: []mux.MiddlewareFunc{
			middlewares.Authorization,
		},
	},
	{
		Name:        "Ping",
		Method:      "GET",
		Pattern:     "/ping",
		HandlerFunc: ping,
	},
}
