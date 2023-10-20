package product

import (
	"fmt"
	"go-auth/src/dtos"
	"go-auth/src/models"
	"go-auth/src/utils"
	"time"

	"github.com/google/uuid"
	"github.com/ylfyt/go_db/go_db"
	"github.com/ylfyt/meta/meta"
)

func createProduct(data dtos.CreateProduct, db *go_db.DB) meta.ResponseDto {
	newId := uuid.New()
	now := time.Now()
	inserted, err := db.Write(`
		INSERT INTO products(id, name, description, price, created_at)
		VALUES($1, $2, $3, $4, $5)
	`, newId, data.Name, data.Description, data.Price, now)
	if err != nil {
		fmt.Println("Err", err)
		return utils.GetInternalErrorResponse("Something wrong")
	}
	if inserted == 0 {
		fmt.Println("Err", err)
		return utils.GetInternalErrorResponse("Something wrong")
	}

	return utils.GetSuccessResponse(models.Product{
		Id:          newId,
		CreatedAt:   now,
		Name:        data.Name,
		Description: data.Description,
		Price:       data.Price,
	})
}
