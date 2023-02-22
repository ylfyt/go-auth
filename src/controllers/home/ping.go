package home

import (
	"go-auth/src/meta"
	"go-auth/src/utils"
)

func ping() meta.ResponseDto {
	return utils.GetSuccessResponse("ok")
}
