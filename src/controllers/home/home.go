package home

import (
	"fmt"
	"go-auth/src/ctx"
	"go-auth/src/dtos"
	"go-auth/src/models"
	"go-auth/src/utils"
	"net/http"
)

func Home(r *http.Request) dtos.Response {
	claims := r.Context().Value(ctx.UserClaimsCtxKey).(models.AccessClaims)
	fmt.Printf("Data: %+v\n", claims)
	return utils.GetSuccessResponse(claims)
}
