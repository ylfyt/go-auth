package auth

import (
	"go-auth/src/db"
	"go-auth/src/dtos"
	"go-auth/src/models"
	"go-auth/src/utils"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func logoutAll(data dtos.Register) dtos.Response {
	conn, err := db.BorrowDbConnection()
	if err != nil {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}
	defer db.ReturnDbConnection(conn)

	user, err := db.GetFirst[models.User](*conn, `
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

	deleted, err := db.Write(*conn, `
		DELETE FROM jwt_tokens WHERE user_id = $1
	`, user.Id)
	if err != nil {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}

	if deleted == 0 {
		return utils.GetErrorResponse(http.StatusBadRequest, "There is no logged in")
	}

	return utils.GetSuccessResponse(true)
}
