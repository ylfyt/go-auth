package auth

import (
	"go-auth/src/db"
	"go-auth/src/meta"
	"go-auth/src/models"
	"go-auth/src/services"
	"go-auth/src/utils"
	"net/http"
)

func getUsers(dbCtx services.DbContext) meta.ResponseDto {
	users, err := db.Get[models.User](dbCtx.Db, `
	SELECT * FROM users
	`)
	if err != nil {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something Wrong!")
	}

	return utils.GetSuccessResponse(users)
}

func getUserById(id string, dbCtx services.DbContext) meta.ResponseDto {
	user, err := db.GetOne[models.User](dbCtx.Db, `
		SELECT * FROM users WHERE id = $1
	`, id)

	if err != nil {
		return utils.GetBadRequestResponse("Payload is not valid")
	}

	if user == nil {
		return utils.GetErrorResponse(http.StatusNotFound, "User not found")
	}

	return utils.GetSuccessResponse(user)
}
