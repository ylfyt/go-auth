package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go-auth/src/ctx"
)

func AccessLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := time.Now().Format("REQ_2006-01-02_15:04:05.000")
		r = r.WithContext(context.WithValue(r.Context(), ctx.ReqIdCtxKey, reqId))

		fmt.Printf("NEW_REQUEST %s : %+v\n", reqId, r)

		next.ServeHTTP(w, r)
	})
}
