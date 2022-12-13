package product

import "go-auth/src/ctx"

var Routes = []ctx.Route{
	{
		Name:        "GetProduct",
		Method:      "GET",
		Pattern:     "/product",
		HandlerFunc: getProduct,
	},
}
