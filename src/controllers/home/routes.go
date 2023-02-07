package home

import (
	"go-auth/src/ctx"
	"go-auth/src/middlewares"

	"github.com/gofiber/fiber/v2"
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
		Middlewares: []func(c *fiber.Ctx) error{
			middlewares.Authorization,
		},
	},
}
