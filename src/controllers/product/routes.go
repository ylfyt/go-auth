package product

import (
	"go-auth/src/meta"
	"go-auth/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

var Routes = []meta.EndPoint{
	{
		Method:      "GET",
		Path:        "/product",
		HandlerFunc: getProduct,
	},
	{
		Method:      "POST",
		Path:        "/product",
		HandlerFunc: createProduct,
		Middlewares: []func(c *fiber.Ctx) error{
			middlewares.Authorization,
		},
	},
}
