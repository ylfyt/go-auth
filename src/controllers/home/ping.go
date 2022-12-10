package home

import (
	"go-auth/src/dtos"
	"go-auth/src/utils"
	"net/http"
)

func ping(r *http.Request) dtos.Response {
	return utils.GetSuccessResponse("ok")
}
