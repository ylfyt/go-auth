package middlewares

import (
	"fmt"
	"go-auth/src/dtos"
	"go-auth/src/utils"
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

type BodyParserOption struct {
	Validate bool
}

func BodyParser[T any](opts ...BodyParserOption) func(http.Handler) http.Handler {
	opt := BodyParserOption{
		Validate: true,
	}
	if len(opts) > 0 {
		opt = opts[0]
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var data T
			err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(body, &data)
			if err != nil {
				fmt.Println("ERR", err)
				w.Header().Add("content-type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				response := dtos.Response{
					Status:  http.StatusBadRequest,
					Message: "Payload is not valid",
					Success: false,
				}
				err := jsoniter.ConfigCompatibleWithStandardLibrary.NewEncoder(w).Encode(response)
				if err != nil {
					fmt.Println("ERR", err)
				}
				return
			}

			if opt.Validate {
				fieldErrors := utils.Validation(&data)
				if len(fieldErrors) != 0 {
					w.Header().Add("content-type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
					response := dtos.Response{
						Status:  http.StatusBadRequest,
						Message: "Payload is not valid",
						Success: false,
						Errors:  fieldErrors,
					}
					err := jsoniter.ConfigCompatibleWithStandardLibrary.NewEncoder(w).Encode(response)
					if err != nil {
						fmt.Println("ERR", err)
					}
					return
				}
				next.ServeHTTP(w, utils.SetContext(r, utils.CTX_BODY_KEY, data))
				return
			}

			next.ServeHTTP(w, utils.SetContext(r, utils.CTX_BODY_KEY, data))
		})
	}
}
