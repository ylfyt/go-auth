package auth

import (
	"fmt"
	"go-auth/src/db"
	"go-auth/src/dtos"
	"go-auth/src/services"
	"go-auth/src/utils"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func register(data dtos.Register, dbCtx services.DbContext) dtos.Response {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong!")
	}

	exists, err := db.GetRowCount(dbCtx.Db, `
		SELECT count(*) FROM users WHERE username = $1
	`, data.Username)
	if err != nil {
		fmt.Println("Error", err)
		return utils.GetErrorResponse(http.StatusInternalServerError, "Something wrong")
	}

	if exists != 0 {
		return utils.GetErrorResponse(http.StatusBadRequest, "Username already exist")
	}

	newId := uuid.New()

	inserted, err := db.Write(dbCtx.Db, `
		INSERT INTO users VALUES($1, $2, $3, NOW())
	`, newId, data.Username, string(hashedPassword))

	if err != nil {
		fmt.Println("Error:", err)
		return utils.GetErrorResponse(http.StatusInternalServerError, "Failed to insert new user")
	}

	if inserted == 0 {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Failed to insert new user")
	}

	return utils.GetSuccessResponse(newId)
}
