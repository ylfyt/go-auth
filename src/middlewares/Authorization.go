package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"go-auth/src/ctx"
)

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Auth Check")

		r = r.WithContext(context.WithValue(r.Context(), ctx.TokenPayloadCtxKey, "1"))

		next.ServeHTTP(w, r)
	})
}
