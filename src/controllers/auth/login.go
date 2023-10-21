package auth

import (
	"fmt"
	"go-auth/src/dtos"
	"go-auth/src/models"
	"go-auth/src/services"
	"go-auth/src/utils"
	"strings"

	"github.com/ylfyt/meta/meta"
)

func (me *AuthController) login(data dtos.Register) meta.ResponseDto {
	var user *models.User
	err := me.db.GetFirst(&user, `
	SELECT * FROM users WHERE username = $1
	`, data.Username)
	if err != nil {
		fmt.Println("ERR", err)
		return utils.GetInternalErrorResponse("Something wrong!")
	}

	if user == nil {
		return utils.GetBadRequestResponse("Username or password is wrong")
	}

	passwordData := strings.Split(user.Password, ":")
	if len(passwordData) != 2 {
		fmt.Println("???")
		return utils.GetInternalErrorResponse("Something wrong!")
	}
	isValid := utils.VerifyPassword(passwordData[0], data.Password, user.Username, []byte(passwordData[1]))
	if !isValid {
		return utils.GetBadRequestResponse("Username or password is wrong")
	}

	refresh, access, jid, err := services.CreateJwtToken(me.config, *user)
	if err != nil {
		fmt.Println("Data:", err)
		return utils.GetInternalErrorResponse("Something wrong!")
	}
	_, err = me.db.Write(`
		INSERT INTO jwt_tokens VALUES($1, $2, NOW())
	`, jid, user.Id)
	if err != nil {
		fmt.Println("Data:", err)
		return utils.GetInternalErrorResponse("Something wrong!")
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
