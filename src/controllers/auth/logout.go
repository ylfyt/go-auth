package auth

import (
	"fmt"
	"go-auth/src/dtos"
	"go-auth/src/services"
	"go-auth/src/utils"
	"net/http"

	"github.com/ylfyt/go_db/go_db"
	"github.com/ylfyt/meta/meta"
)

func logout(data dtos.RefreshPayload, db *go_db.DB) meta.ResponseDto {
	valid, jid := services.VerifyRefreshToken(data.RefreshToken)
	if !valid {
		return utils.GetErrorResponse(http.StatusBadRequest, "Token is not valid")
	}

	deleted, err := db.Write(`
		DELETE FROM jwt_tokens WHERE id = $1
	`, jid)
	if err != nil {
		fmt.Println("Err", err)
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}

	if deleted == 0 {
		return utils.GetErrorResponse(http.StatusBadRequest, "Token is not found")
	}

	// TODO: Implement Blacklist Token Mechanism

	return utils.GetSuccessResponse(true)
}
