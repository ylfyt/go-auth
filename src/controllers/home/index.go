package home

import (
	"go-auth/src/ctx"
	"go-auth/src/dtos"
	"go-auth/src/utils"
	"net/http"
	"go-auth/src/l"
)

func home(r *http.Request) dtos.Response {
	reqId := ctx.GetReqIdCtx(r)
	l.I(reqId)
	return utils.GetSuccessResponse(reqId)
}
