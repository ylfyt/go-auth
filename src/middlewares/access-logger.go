package middlewares

import (
	"context"
	"go-auth/src/ctx"
	"go-auth/src/l"
	"net/http"
	"time"
)

func AccessLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := time.Now().Format("REQ_2006-01-02_15:04:05.000")
		r = r.WithContext(context.WithValue(r.Context(), ctx.ReqIdCtxKey, reqId))

		l.I("[%s] NEW REQUEST: %s %s", reqId, r.Method, r.URL)

		next.ServeHTTP(w, r)
	})
}
