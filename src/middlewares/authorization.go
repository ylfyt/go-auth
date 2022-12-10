package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go-auth/src/ctx"
	"go-auth/src/dtos"
	"go-auth/src/models"
	"go-auth/src/services"
)

func validate(authHeader string) (bool, models.AccessClaims) {
	if authHeader == "" {
		return false, models.AccessClaims{}
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	valid, claims := services.VerifyAccessToken(token)

	return valid, claims
}

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := r.Context().Value(ctx.ReqIdCtxKey)

		authHeader := r.Header.Get("authorization")

		valid, claims := validate(authHeader)
		if !valid {
			response := dtos.Response{
				Status:  http.StatusUnauthorized,
				Message: "Unauthorized",
				Success: false,
				Data:    nil,
			}
			w.Header().Add("content-type", "application/json")
			w.WriteHeader(response.Status)
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				fmt.Println("Failed to send response", err)
			}
			fmt.Printf("[%s] REQUEST FAILED with RESPONSE:%+v\n", reqId, response)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), ctx.UserClaimsCtxKey, claims))

		fmt.Printf("[%s] REQUEST AUTHORIZED with User Claims: %+v\n", reqId, claims)

		next.ServeHTTP(w, r)
	})
}
