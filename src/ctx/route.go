package ctx

import (
	"go-auth/src/dtos"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(r *http.Request) dtos.Response
	Middlewares []mux.MiddlewareFunc
}
