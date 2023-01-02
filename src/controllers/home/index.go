package home

import (
	"fmt"
	"go-auth/src/ctx"
	"go-auth/src/db"
	"go-auth/src/dtos"
	"go-auth/src/l"
	"go-auth/src/utils"
	"net/http"
)

func home(r *http.Request, db db.DbConnection) dtos.Response {
	reqId := ctx.GetReqIdCtx(r)
	l.I(reqId)
	fmt.Printf("Data: %+v\n", db)
	return utils.GetSuccessResponse(reqId)
}
