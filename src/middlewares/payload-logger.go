package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-auth/src/ctx"
	"go-auth/src/meta"
	"io"
	"net/http"
)

func PayloadLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := r.Context().Value(ctx.ReqIdCtxKey)
		body, err := io.ReadAll(r.Body)
		if err == nil {
			r.Body.Close()
			r.Body = io.NopCloser(bytes.NewBuffer(body))
			fmt.Printf("[%s] REQUEST PAYLOAD : %s\n", reqId, string(body))

			next.ServeHTTP(w, r)
			return
		}

		response := meta.ResponseDto{
			Status:  http.StatusInternalServerError,
			Message: "Something wrong",
			Success: false,
			Data:    nil,
		}
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(response.Status)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			fmt.Println("Failed to send response", err)
		}
		fmt.Printf("[%s] REQUEST FAILED with RESPONSE:%+v\n", reqId, response)
	})
}
