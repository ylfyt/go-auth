package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"go-auth/src/dtos"
	"go-auth/src/services"
	"go-auth/src/shared"
	"go-auth/src/utils"

	jsoniter "github.com/json-iterator/go"
)

func validate(authHeader string, config *shared.EnvConf) (bool, dtos.AccessClaims) {
	if authHeader == "" {
		return false, dtos.AccessClaims{}
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	valid, claims := services.VerifyAccessToken(config, token)

	return valid, claims
}

func Authorization(config *shared.EnvConf) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqId := utils.GetContext[int64](r, utils.CTX_REQ_ID_KEY)

			valid, claim := validate(r.Header.Get("Authorization"), config)
			if valid {
				fmt.Printf("[%d] AUTHORIZED: %+v\n", *reqId, claim)
				next.ServeHTTP(w, utils.SetContext(r, utils.CTX_AUTH_CLAIM_KEY, claim))
				return
			}
			fmt.Printf("[%d] UNAUTHORIZED\n", reqId)

			w.Header().Add("content-type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			response := dtos.Response{
				Code:  http.StatusUnauthorized,
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
