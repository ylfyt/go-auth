package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"go-auth/src/models"
	"go-auth/src/services"
	"go-auth/src/structs"

	"go-auth/src/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/ylfyt/meta/meta"
)

func validate(authHeader string, config *structs.EnvConf) (bool, models.AccessClaims) {
	if authHeader == "" {
		return false, models.AccessClaims{}
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	valid, claims := services.VerifyAccessToken(config, token)

	return valid, claims
}

func Authorization(c *fiber.Ctx, config *structs.EnvConf) error {
	reqId := utils.GetContext[string](c, "reqId")
	authHeader := ""
	if len(c.GetReqHeaders()["Authorization"]) > 0 {
		authHeader = c.GetReqHeaders()["Authorization"][0]
	}

	valid, claims := validate(authHeader, config)
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
