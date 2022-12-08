package controllers

import (
	"encoding/json"
	"go-auth/src/db"
	"go-auth/src/dtos"
	"go-auth/src/models"
	"io"
	"net/http"

	"github.com/google/uuid"
)

func Register(r *http.Request) dtos.Response {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return getErrorResponse(http.StatusBadRequest, "Failed to get request body")
	}

	var data dtos.Register
	err = json.Unmarshal(body, &data)
	if err != nil {
		return getErrorResponse(http.StatusBadRequest, "Failed to get payload")
	}

	if len(data.Username) < 4 || len(data.Username) > 20 {
		return getErrorResponse(http.StatusBadRequest, "Username should be > 4 and < 20")
	}

	if len(data.Password) < 4 {
		return getErrorResponse(http.StatusBadRequest, "Password should be > 4")

	}

	conn, err := db.BorrowDbConnection()
	if err != nil {
		return getErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}
	defer db.ReturnDbConnection(conn)

	user := models.User{
		Id:       uuid.New(),
		Username: data.Username,
		Password: data.Password,
	}

	inserted := conn.WriteQuery(`
		INSERT INTO users VALUES($1, $2, $3)
	`, user.Id, user.Username, user.Password)

	if inserted == 0 {
		return getErrorResponse(http.StatusInternalServerError, "Failed to insert new user")
	}

	return getSuccessResponse(user)
}
