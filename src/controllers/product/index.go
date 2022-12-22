package product

import (
	"fmt"
	"go-auth/src/db"
	"go-auth/src/dtos"
	"go-auth/src/models"
	"go-auth/src/utils"
	"net/http"
)

func getProduct() dtos.Response {
	conn, err := db.BorrowDbConnection()
	if err != nil {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong")
	}
	defer db.ReturnDbConnection(conn)

	products, err := db.Get[models.Product](*conn, `SELECT * FROM products`)
	if err != nil {
		fmt.Println("Err", err)
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong")
	}

	return utils.GetSuccessResponse(products)
}
