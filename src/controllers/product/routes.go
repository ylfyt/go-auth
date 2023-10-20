package product

import (
	"go-auth/src/middlewares"

	"github.com/ylfyt/meta/meta"
)

var Routes = []meta.EndPoint{
	{
		Method:      "GET",
		Path:        "/product",
		HandlerFunc: getProduct,
	},
	{
		Method:      "POST",
		Path:        "/product",
		HandlerFunc: createProduct,
		Middlewares: []any{
			middlewares.Authorization,
		},
	},
}
