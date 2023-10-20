package auth

import (
	"go-auth/src/models"
	"go-auth/src/utils"
	"net/http"

	"github.com/ylfyt/meta/meta"
)

func (me *AuthController) getUsers() meta.ResponseDto {
	var users []models.User
	err := me.db.Get(&users, `
	SELECT * FROM users
	`)
	if err != nil {
		return utils.GetInternalErrorResponse("Something wrong!")
	}

	return utils.GetSuccessResponse(users)
}

func (me *AuthController) getUserById(id string) meta.ResponseDto {
	var user *models.User
	err := me.db.GetFirst(&user, `
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
