package home

import (
	"go-auth/src/ctx"
	"go-auth/src/middlewares"
)

var Routes = []ctx.Route{
	{
		Name:        "Home",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: home,
	},
	{
		Name:        "Ping",
		Method:      "GET",
		Pattern:     "/ping",
		HandlerFunc: ping,
		Middlewares: []interface{}{
			middlewares.Authorization,
		},
	},
}
