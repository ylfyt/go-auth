package middlewares

import (
	"fmt"
	"go-auth/src/utils"
	"net/http"
	"time"
)

func AccessLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := time.Now().Unix()
		fmt.Println("NEW REQ:", reqId)
		next.ServeHTTP(w, utils.SetContext(r, utils.CTX_REQ_ID_KEY, reqId))
	})
}
