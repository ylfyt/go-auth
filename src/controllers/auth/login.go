package auth

import (
	"go-auth/src/config"
	"go-auth/src/dtos"
	"go-auth/src/models"
	"go-auth/src/services"
	"go-auth/src/utils"
	"net/http"
	"strings"

	"github.com/ylfyt/go_db/go_db"
	"github.com/ylfyt/meta/meta"
)

func login(data dtos.Register, db *go_db.DB) meta.ResponseDto {
	var user *models.User
	err := db.GetFirst(&user, `
	SELECT * FROM users WHERE username = $1
	`, data.Username)
	if err != nil {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}

	if user == nil {
		return utils.GetErrorResponse(http.StatusBadRequest, "Username or password is wrong")
	}

	passwordData := strings.Split(user.Password, ":")
	if len(passwordData) != 2 {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}
	isValid := utils.VerifyPassword(passwordData[0], data.Password, user.Username, []byte(passwordData[1]))
	if !isValid {
		return utils.GetErrorResponse(http.StatusBadRequest, "Username or password is wrong")
	}

	refresh, access, jid, err := services.CreateJwtToken(*user)
	if err != nil {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}
	_, err = db.Write(`
		INSERT INTO jwt_tokens VALUES($1, $2, NOW())
	`, jid, user.Id)
	if err != nil {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}

	return utils.GetSuccessResponse(dtos.LoginResponse{
		User: *user,
		Token: dtos.TokenPayload{
			RefreshToken: refresh,
			AccessToken:  access,
			ExpiredIn:    int64(config.JWT_ACCESS_TOKEN_EXPIRY_TIME),
		},
	})
}
