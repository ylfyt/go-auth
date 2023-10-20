package product

import (
	"fmt"
	"go-auth/src/models"
	"go-auth/src/utils"

	"github.com/ylfyt/go_db/go_db"
	"github.com/ylfyt/meta/meta"
)

func getProduct(db *go_db.DB) meta.ResponseDto {
	var products []models.Product
	err := db.Get(&products, `SELECT * FROM products`)
	if err != nil {
		fmt.Println("Err", err)
		return utils.GetInternalErrorResponse("Something wrong")
	}

	return utils.GetSuccessResponse(products)
}
