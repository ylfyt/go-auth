package auth

import (
	"fmt"
	"go-auth/src/dtos"
	"go-auth/src/models"
	"go-auth/src/services"
	"go-auth/src/utils"
	"net/http"

	"github.com/ylfyt/meta/meta"
)

func (me *AuthController) refreshToken(data dtos.RefreshPayload) meta.ResponseDto {
	valid, jid := services.VerifyRefreshToken(me.config, data.RefreshToken)
	if !valid {
		return utils.GetErrorResponse(http.StatusBadRequest, "Token is not valid")
	}

	var token *models.JwtToken
	err := me.db.GetFirst(&token, `
		SELECT * FROM jwt_tokens WHERE id = $1
	`, jid)
	if err != nil {
		fmt.Println("ERROR", err)
		return utils.GetInternalErrorResponse("Something wrong!")
	}
	if token == nil {
		return utils.GetErrorResponse(http.StatusBadRequest, "Token is not found")
	}

	var user *models.User
	err = me.db.GetFirst(&user, `
	SELECT * FROM users WHERE id = $1
	`, token.UserId)
	if err != nil {
		return utils.GetInternalErrorResponse("Something wrong!")
	}
	if user == nil {
		return utils.GetBadRequestResponse("User is not found")
	}

	refresh, access, newJid, err := services.CreateJwtToken(me.config, *user)
	if err != nil {
		return utils.GetInternalErrorResponse("Something wrong!")
	}
	_, err = me.db.Write(`
		INSERT INTO jwt_tokens VALUES($1, $2, NOW())
	`, newJid, user.Id)
	if err != nil {
		return utils.GetInternalErrorResponse("Something wrong!")
	}

	_, err = me.db.Write(`
	DELETE FROM jwt_tokens WHERE id = $1
	`, jid)
	if err != nil {
		fmt.Println("ERROR", err)
	}

	return utils.GetSuccessResponse(dtos.LoginResponse{
		User: *user,
		Token: dtos.TokenPayload{
			RefreshToken: refresh,
			AccessToken:  access,
			ExpiredIn:    int64(me.config.JwtAccessTokenExpiryTime),
		},
	})
}
