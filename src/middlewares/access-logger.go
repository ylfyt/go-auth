package middlewares

import (
	"go-auth/src/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AccessLogger(c *fiber.Ctx) error {
	reqId := time.Now().Format("REQ_2006-01-02_15:04:05.000")
	utils.SetContext(c, "reqId", reqId)

	return c.Next()
}
