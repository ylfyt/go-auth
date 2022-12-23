package auth

import (
	"fmt"
	"go-auth/src/db"
	"go-auth/src/dtos"
	"go-auth/src/services"
	"go-auth/src/utils"
	"net/http"
)

func logout(data dtos.RefreshPayload) dtos.Response {
	valid, jid := services.VerifyRefreshToken(data.RefreshToken)
	if !valid {
		return utils.GetErrorResponse(http.StatusBadRequest, "Token is not valid")
	}

	conn, err := db.BorrowDbConnection()
	if err != nil {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}
	defer db.ReturnDbConnection(conn)

	deleted, err := db.Write(conn, `
		DELETE FROM jwt_tokens WHERE id = $1
	`, jid)
	if err != nil {
		fmt.Println("Err", err)
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}

	if deleted == 0 {
		return utils.GetErrorResponse(http.StatusBadRequest, "Token is not found")
	}

	// Implement Blacklist Token Mechanism

	return utils.GetSuccessResponse(true)
}
