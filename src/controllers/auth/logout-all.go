package auth

import (
	"encoding/json"
	"go-auth/src/db"
	"go-auth/src/dtos"
	"go-auth/src/models"
	"go-auth/src/utils"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func logoutAll(r *http.Request) dtos.Response {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return utils.GetErrorResponse(http.StatusBadRequest, "Failed to get request body")
	}

	var data dtos.Register
	err = json.Unmarshal(body, &data)
	if err != nil {
		return utils.GetErrorResponse(http.StatusBadRequest, "Failed to get payload")
	}

	if len(data.Username) < 4 || len(data.Username) > 20 {
		return utils.GetErrorResponse(http.StatusBadRequest, "Username should be > 4 and < 20")
	}

	if len(data.Password) < 4 {
		return utils.GetErrorResponse(http.StatusBadRequest, "Password should be > 4")
	}

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
