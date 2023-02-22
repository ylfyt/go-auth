package auth

import (
	"go-auth/src/config"
	"go-auth/src/db"
	"go-auth/src/dtos"
	"go-auth/src/meta"
	"go-auth/src/models"
	"go-auth/src/services"
	"go-auth/src/utils"
	"net/http"
)

func refreshToken(data dtos.RefreshPayload, dbCtx services.DbContext) meta.ResponseDto {
	valid, jid := services.VerifyRefreshToken(data.RefreshToken)
	if !valid {
		return utils.GetErrorResponse(http.StatusBadRequest, "Token is not valid")
	}

	token, err := db.GetOne[models.JwtToken](dbCtx.Db, `
		SELECT * FROM jwt_tokens WHERE id = $1
	`, jid)
	if err != nil {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}
	if token == nil {
		return utils.GetErrorResponse(http.StatusBadRequest, "Token is not found")
	}

	user, err := db.GetOne[models.User](dbCtx.Db, `
	SELECT * FROM users WHERE id = $1
	`, token.UserId)
	if err != nil {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}

	refresh, access, newJid, err := services.CreateJwtToken(*user)
	if err != nil {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}
	inserted, _ := db.Write(dbCtx.Db, `
		INSERT INTO jwt_tokens VALUES($1, $2, NOW())
	`, newJid, user.Id)
	if inserted == 0 {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}

	deleted, _ := db.Write(dbCtx.Db, `
	DELETE FROM jwt_tokens WHERE id = $1
	`, jid)
	if deleted == 0 {
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
