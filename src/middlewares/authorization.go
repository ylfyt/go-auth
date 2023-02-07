package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"go-auth/src/ctx"
	"go-auth/src/dtos"
	"go-auth/src/models"
	"go-auth/src/services"

	"github.com/gofiber/fiber/v2"
)

func validate(authHeader string) (bool, models.AccessClaims) {
	if authHeader == "" {
		return false, models.AccessClaims{}
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	valid, claims := services.VerifyAccessToken(token)

	return valid, claims
}

func Authorization(c *fiber.Ctx) error {
	reqId := c.Context().Value(ctx.ReqIdCtxKey)
	authHeader := c.GetReqHeaders()["authorization"]

	valid, claims := validate(authHeader)
	if valid {
		c.Locals(ctx.UserClaimsCtxKey, claims)
		fmt.Printf("[%s] REQUEST AUTHORIZED with User Claims: %+v\n", reqId, claims)
		return c.Next()
	}

	response := dtos.Response{
		Status:  http.StatusUnauthorized,
		Message: "Unauthorized",
		Success: false,
		Data:    nil,
	}

	fmt.Printf("[%s] REQUEST FAILED with RESPONSE:%+v\n", reqId, response)
	return c.Status(response.Status).JSON(response)
}
