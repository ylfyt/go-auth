package product

import (
	"go-auth/src/ctx"
	"go-auth/src/middlewares"

	"github.com/gorilla/mux"
)

var Routes = []ctx.Route{
	{
		Name:        "GetProduct",
		Method:      "GET",
		Pattern:     "/product",
		HandlerFunc: getProduct,
	},
	{
		Name:        "CreateProduct",
		Method:      "POST",
		Pattern:     "/product",
		HandlerFunc: createProduct,
		Middlewares: []mux.MiddlewareFunc{
			middlewares.Authorization,
		},
	},
}
