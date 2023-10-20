package home

import (
	"go-auth/src/utils"

	"github.com/ylfyt/meta/meta"
)

func (me *HomeController) ping() meta.ResponseDto {
	return utils.GetSuccessResponse("pong")
}
