package home

import (
	"go-auth/src/meta"
	"go-auth/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

var Routes = []meta.EndPoint{
	{
		Method:      "GET",
		Path:        "/",
		HandlerFunc: home,
	},
	{
		Method:      "GET",
		Path:        "/ping",
		HandlerFunc: ping,
		Middlewares: []func(c *fiber.Ctx) error{
			middlewares.Authorization,
		},
	},
}
