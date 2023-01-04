package home

import (
	"fmt"
	"go-auth/src/ctx"
	"go-auth/src/db"
	"go-auth/src/dtos"
	"go-auth/src/l"
	"go-auth/src/models"
	"go-auth/src/services"
	"go-auth/src/utils"
	"net/http"
)

func home(r *http.Request, dbCtx services.DbContext) dtos.Response {
	reqId := ctx.GetReqIdCtx(r)
	l.I(reqId)

	product, err := db.GetOne[models.Product](dbCtx.Db, `SELECT * FROM products LIMIT 1`)

	if err != nil {
		fmt.Println("Err", err)
		return utils.GetInternalErrorResponse("Something wrong")
	}

	return utils.GetSuccessResponse(product)
}
