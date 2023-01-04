package auth

import (
	"go-auth/src/config"
	"go-auth/src/db"
	"go-auth/src/dtos"
	"go-auth/src/models"
	"go-auth/src/services"
	"go-auth/src/utils"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func login(data dtos.Register, dbCtx services.DbContext) dtos.Response {
	user, err := db.GetOne[models.User](dbCtx.Db, `
	SELECT * FROM users WHERE username = $1
	`, data.Username)
	if err != nil {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}

	if user == nil {
		return utils.GetErrorResponse(http.StatusBadRequest, "Username or password is wrong")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		return utils.GetErrorResponse(http.StatusBadRequest, "Username or password is wrong")
	}

	refresh, access, jid, err := services.CreateJwtToken(*user)
	if err != nil {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}
	inserted, _ := db.Write(dbCtx.Db, `
		INSERT INTO jwt_tokens VALUES($1, $2, NOW())
	`, jid, user.Id)
	if inserted == 0 {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}

	return utils.GetSuccessResponse(dtos.LoginResponse{
		User: *user,
		Token: dtos.TokenPayload{
			RefreshToken: refresh,
			AccessToken:  access,
			ExpiredIn: int64(config.JWT_ACCESS_TOKEN_EXPIRY_TIME),
		},
	})
}