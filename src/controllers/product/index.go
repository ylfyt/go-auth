package product

import (
	"fmt"
	"go-auth/src/middlewares"
	"go-auth/src/models"
	"go-auth/src/utils"

	"github.com/ylfyt/go_db/go_db"
	"github.com/ylfyt/meta/meta"
)

type ProductController struct {
	db *go_db.DB
}

func (me *ProductController) Setup(db *go_db.DB) []meta.EndPoint {
	me.db = db

	return []meta.EndPoint{
		{
			Method:      "GET",
			Path:        "/product",
			HandlerFunc: me.getProduct,
		},
		{
			Method:      "POST",
			Path:        "/product",
			HandlerFunc: me.createProduct,
			Middlewares: []any{
				middlewares.Authorization,
			},
		},
	}

}

func (me *ProductController) getProduct() meta.ResponseDto {
	var products []models.Product
	err := me.db.Get(&products, `SELECT * FROM products`)
	if err != nil {
		fmt.Println("Err", err)
		return utils.GetInternalErrorResponse("Something wrong")
	}

	return utils.GetSuccessResponse(products)
}
