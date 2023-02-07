package middlewares

import (
	"go-auth/src/ctx"
	"go-auth/src/l"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AccessLogger(c *fiber.Ctx) error {
	reqId := time.Now().Format("REQ_2006-01-02_15:04:05.000")
	c.Locals(ctx.ReqIdCtxKey, reqId)

	l.I("[%s] NEW REQUEST:", reqId)

	return c.Next()
}
