package home

import (
	"go-auth/src/ctx"
	"go-auth/src/dtos"
	"go-auth/src/utils"
	"net/http"
)

func home(r *http.Request) dtos.Response {
	reqId := ctx.GetReqIdCtx(r)
	return utils.GetSuccessResponse(reqId)
}
