package product

import (
	"go-auth/src/ctx"
	"go-auth/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

var Routes = []ctx.Route{
	{
		Name:        "GetProduct",
		Method:      "GET",
		Pattern:     "/product",
		HandlerFunc: getProduct,
	},
	{
		Name:        "CreateProduct",
		Method:      "POST",
		Pattern:     "/product",
		HandlerFunc: createProduct,
		Middlewares: []func(c *fiber.Ctx) error{
			middlewares.Authorization,
		},
	},
}
