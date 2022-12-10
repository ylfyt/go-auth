package auth

import (
	"encoding/json"
	"go-auth/src/db"
	"go-auth/src/dtos"
	"go-auth/src/models"
	"go-auth/src/services"
	"go-auth/src/utils"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Login(r *http.Request) dtos.Response {
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

	refresh, access, jid, err := services.CreateJwtToken(*user)
	if err != nil {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}
	inserted, _ := db.Write(*conn, `
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
		},
	})
}