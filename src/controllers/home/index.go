package home

import (
	"fmt"
	"go-auth/src/ctx"
	"go-auth/src/dtos"
	"go-auth/src/utils"
	"net/http"
)

func home(r *http.Request) dtos.Response {
	claims := ctx.GetUserClaimsCtx(r)
	fmt.Printf("Data: %+v\n", claims)
	return utils.GetSuccessResponse(claims)
}
