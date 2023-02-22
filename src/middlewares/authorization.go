package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"go-auth/src/meta"
	"go-auth/src/models"
	"go-auth/src/services"

	"github.com/gofiber/fiber/v2"
	"go-auth/src/utils"
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
	reqId := utils.GetContext[string](c, "reqId")
	authHeader := c.GetReqHeaders()["Authorization"]

	valid, claims := validate(authHeader)
	if valid {
		utils.SetContext(c, "claims", claims)
		fmt.Printf("[%s] REQUEST AUTHORIZED with User Claims: %+v\n", *reqId, claims)
		return c.Next()
	}

	response := meta.ResponseDto{
		Status:  http.StatusUnauthorized,
		Message: "Unauthorized",
		Success: false,
		Data:    nil,
	}

	fmt.Printf("[%s] REQUEST FAILED with RESPONSE:%+v\n", *reqId, response)
	return c.Status(response.Status).JSON(response)
}
