package home

import (
	"database/sql"
	"fmt"
	"go-auth/src/ctx"
	"go-auth/src/db"
	"go-auth/src/l"
	"go-auth/src/meta"
	"go-auth/src/models"
	"go-auth/src/utils"

	"github.com/gofiber/fiber/v2"
)

func home(c *fiber.Ctx, dbCtx *sql.DB) meta.ResponseDto {
	reqId := c.Locals(ctx.ReqIdCtxKey)
	l.I(reqId.(string))
	product, err := db.GetOne[models.Product](dbCtx, `SELECT * FROM products LIMIT 1`)

	if err != nil {
		fmt.Println("Err", err)
		return utils.GetInternalErrorResponse("Something wrong")
	}

	return utils.GetSuccessResponse(product)
}
