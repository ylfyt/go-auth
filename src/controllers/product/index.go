package product

import (
	"fmt"
	"go-auth/src/db"
	"go-auth/src/meta"
	"go-auth/src/models"
	"go-auth/src/services"
	"go-auth/src/utils"
	"net/http"
)

func getProduct(dbCtx services.DbContext) meta.ResponseDto {
	products, err := db.Get[models.Product](dbCtx.Db, `SELECT * FROM products`)
	if err != nil {
		fmt.Println("Err", err)
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong")
	}

	return utils.GetSuccessResponse(products)
}
