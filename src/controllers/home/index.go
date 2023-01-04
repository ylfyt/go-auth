package home

import (
	"fmt"
	"go-auth/src/ctx"
	"go-auth/src/db"
	"go-auth/src/dtos"
	"go-auth/src/l"
	// "go-auth/src/models"
	"go-auth/src/services"
	"go-auth/src/utils"
	"net/http"
)

func home(r *http.Request, dbCtx services.DbContext) dtos.Response {
	reqId := ctx.GetReqIdCtx(r)
	l.I(reqId)

	name, err := db.GetFieldFirst[int](dbCtx.Db, `SELECT name FROM products LIMIT 1`)

	if err != nil {
		fmt.Println("Err", err)
		return utils.GetInternalErrorResponse("Something wrong")
	}

	fmt.Println("Data:", name)

	return utils.GetSuccessResponse(name)
}
