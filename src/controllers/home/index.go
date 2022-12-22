package home

import (
	"go-auth/src/dtos"
	"go-auth/src/utils"
)

func home() dtos.Response {
	// claims := ctx.GetUserClaimsCtx(r)
	// fmt.Printf("Data: %+v\n", claims)
	return utils.GetSuccessResponse("TODO: CTX")
}
