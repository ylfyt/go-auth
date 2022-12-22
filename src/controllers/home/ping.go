package home

import (
	"go-auth/src/dtos"
	"go-auth/src/utils"
)

func ping() dtos.Response {
	return utils.GetSuccessResponse("ok")
}
