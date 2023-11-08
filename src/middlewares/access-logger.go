package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

func AccessLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := time.Now().Unix()
		fmt.Println("Req", reqId)

		next.ServeHTTP(w, r)
	})
}
