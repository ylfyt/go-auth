package auth

import (
	"fmt"
	"go-auth/src/db"
	"go-auth/src/dtos"
	"go-auth/src/models"
	"go-auth/src/services"
	"go-auth/src/utils"
	"net/http"
	"strings"
)

func logoutAll(data dtos.Register, dbCtx services.DbContext) dtos.Response {
	user, err := db.GetOne[models.User](dbCtx.Db, `
	SELECT * FROM users WHERE username = $1
	`, data.Username)
	if err != nil {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}

	if user == nil {
		return utils.GetErrorResponse(http.StatusBadRequest, "Username or password is wrong")
	}

	passwordData := strings.Split(user.Password, ":")
	if len(passwordData) != 2 {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}
	isValid := utils.VerifyPassword(passwordData[0], data.Password, user.Username, []byte(passwordData[1]))
	if !isValid {
		return utils.GetErrorResponse(http.StatusBadRequest, "Username or password is wrong")
	}

	deleted, err := db.Write(dbCtx.Db, `
		DELETE FROM jwt_tokens WHERE user_id = $1
	`, user.Id)
	if err != nil {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}

	if deleted == 0 {
		return utils.GetErrorResponse(http.StatusBadRequest, "There is no logged in")
	}

	return utils.GetSuccessResponse(fmt.Sprintf("Logout from %d devices", deleted))
}
