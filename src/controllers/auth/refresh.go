package auth

import (
	"encoding/json"
	"fmt"
	"go-auth/src/db"
	"go-auth/src/dtos"
	"go-auth/src/models"
	"go-auth/src/services"
	"go-auth/src/utils"
	"io"
	"net/http"
)

func RefreshToken(r *http.Request) dtos.Response {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return utils.GetErrorResponse(http.StatusBadRequest, "Failed to get request body")
	}

	var data dtos.RefreshPayload
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error", err)
		return utils.GetErrorResponse(http.StatusBadRequest, "Failed to get payload")
	}

	valid, jid := services.VerifyRefreshToken(data.RefreshToken)
	if !valid {
		return utils.GetErrorResponse(http.StatusBadRequest, "Token is not valid")
	}

	conn, err := db.BorrowDbConnection()
	if err != nil {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}
	defer db.ReturnDbConnection(conn)

	token, err := db.GetFirst[models.JwtToken](*conn, `
		SELECT * FROM jwt_tokens WHERE id = $1
	`, jid)
	if err != nil {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}
	if token == nil {
		return utils.GetErrorResponse(http.StatusBadRequest, "Token is not found")
	}

	user, err := db.GetFirst[models.User](*conn, `
	SELECT * FROM users WHERE id = $1
	`, token.UserId)
	if err != nil {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}

	refresh, access, newJid, err := services.CreateJwtToken(*user)
	if err != nil {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}
	inserted, _ := db.Write(*conn, `
		INSERT INTO jwt_tokens VALUES($1, $2, NOW())
	`, newJid, user.Id)
	if inserted == 0 {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}

	deleted, _ := db.Write(*conn, `
	DELETE FROM jwt_tokens WHERE id = $1
	`, jid)
	if deleted == 0 {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}

	return utils.GetSuccessResponse(dtos.TokenPayload{
		RefreshToken: refresh,
		AccessToken:  access,
	})
}