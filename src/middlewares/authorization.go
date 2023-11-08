package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"go-auth/src/dtos"
	"go-auth/src/models"
	"go-auth/src/services"
	"go-auth/src/structs"

	jsoniter "github.com/json-iterator/go"
)

func validate(authHeader string, config *structs.EnvConf) (bool, models.AccessClaims) {
	if authHeader == "" {
		return false, models.AccessClaims{}
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	valid, claims := services.VerifyAccessToken(config, token)

	return valid, claims
}

func Authorization(config *structs.EnvConf) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			valid, claims := validate(r.Header.Get("Authorization"), config)
			if valid {
				fmt.Println("OK", claims)
				next.ServeHTTP(w, r)
				return
			}
			fmt.Println("UNAUTHORIZED")

			w.Header().Add("content-type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			response := dtos.Response{
				Status:  http.StatusUnauthorized,
				Message: "UNAUTHORIZED",
				Success: false,
			}
			err := jsoniter.ConfigCompatibleWithStandardLibrary.NewEncoder(w).Encode(response)
			if err != nil {
				fmt.Println("ERR", err)
			}
		})
	}
}
