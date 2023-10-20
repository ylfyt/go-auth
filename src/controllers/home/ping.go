package home

import (
	"go-auth/src/utils"

	"github.com/ylfyt/meta/meta"
)

func ping() meta.ResponseDto {
	return utils.GetSuccessResponse("pong")
}
