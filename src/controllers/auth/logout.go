package auth

import (
	"fmt"
	"go-auth/src/dtos"
	"go-auth/src/services"
	"go-auth/src/utils"

	"github.com/ylfyt/meta/meta"
)

func (me *AuthController) logout(data dtos.RefreshPayload) meta.ResponseDto {
	valid, jid := services.VerifyRefreshToken(data.RefreshToken)
	if !valid {
		return utils.GetBadRequestResponse("Token is not valid")
	}

	deleted, err := me.db.Write(`
		DELETE FROM jwt_tokens WHERE id = $1
	`, jid)
	if err != nil {
		fmt.Println("Err", err)
		return utils.GetInternalErrorResponse("Something wrong!")
	}

	if deleted == 0 {
		return utils.GetBadRequestResponse("Token is not found")
	}

	// TODO: Implement Blacklist Token Mechanism

	return utils.GetSuccessResponse(true)
}
