package utils

import "github.com/gofiber/fiber/v2"

func GetContext[T any](c *fiber.Ctx, key string) *T {
	val := c.Locals(key)
	if val == nil {
		return nil
	}

	if z, ok := val.(T); ok {
		return &z
	}

	return nil
}

func SetContext(c *fiber.Ctx, key string, val interface{}) {
	c.Locals(key, val)
}
