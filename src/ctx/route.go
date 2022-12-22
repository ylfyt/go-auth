package ctx

import (
	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc interface{}
	Middlewares []mux.MiddlewareFunc
}
