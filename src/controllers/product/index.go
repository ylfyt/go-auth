package product

import (
	"database/sql"
	"fmt"
	"go-auth/src/db"
	"go-auth/src/meta"
	"go-auth/src/models"
	"go-auth/src/utils"
	"net/http"
)

func getProduct(dbCtx *sql.DB) meta.ResponseDto {
	products, err := db.Get[models.Product](dbCtx, `SELECT * FROM products`)
	if err != nil {
		fmt.Println("Err", err)
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong")
	}

	return utils.GetSuccessResponse(products)
}
