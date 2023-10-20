package home

import (
	"database/sql"
	"fmt"
	"go-auth/src/db"
	"go-auth/src/models"
	"go-auth/src/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/ylfyt/meta/meta"
)

func home(c *fiber.Ctx, dbCtx *sql.DB) meta.ResponseDto {
	reqId := utils.GetContext[string](c, "reqId")
	fmt.Println("ReqId:", *reqId)
	product, err := db.GetOne[models.Product](dbCtx, `SELECT * FROM products LIMIT 1`)

	if err != nil {
		fmt.Println("Err", err)
		return utils.GetInternalErrorResponse("Something wrong")
	}

	return utils.GetSuccessResponse(product)
}
