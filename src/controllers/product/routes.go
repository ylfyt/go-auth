package product

import (
	"go-auth/src/ctx"
	"go-auth/src/middlewares"
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
		Middlewares: []interface{}{
			middlewares.Authorization,
		},
	},
}
