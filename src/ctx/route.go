package ctx

import "github.com/gofiber/fiber/v2"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc interface{}
	Middlewares []func(c *fiber.Ctx) error
}
