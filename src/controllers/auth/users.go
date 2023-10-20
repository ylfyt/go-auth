package auth

import (
	"go-auth/src/models"
	"go-auth/src/utils"
	"net/http"

	"github.com/ylfyt/go_db/go_db"
	"github.com/ylfyt/meta/meta"
)

func getUsers(db *go_db.DB) meta.ResponseDto {
	var users []models.User
	err := db.Get(&users, `
	SELECT * FROM users
	`)
	if err != nil {
		return utils.GetInternalErrorResponse("Something wrong!")
	}

	return utils.GetSuccessResponse(users)
}

func getUserById(id string, db *go_db.DB) meta.ResponseDto {
	var user *models.User
	err := db.GetFirst(&user, `
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
