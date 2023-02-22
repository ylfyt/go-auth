package auth

import (
	"fmt"
	"go-auth/src/db"
	"go-auth/src/dtos"
	"go-auth/src/meta"
	"go-auth/src/services"
	"go-auth/src/utils"
	"net/http"

	"github.com/google/uuid"
)

func register(data dtos.Register, dbCtx services.DbContext) meta.ResponseDto {
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

	rawSalt := utils.GenerateRawSalt()
	realSalt := utils.GetRealSalt(rawSalt, data.Username)
	hashedPassword := utils.HashPassword(data.Password, realSalt)

	newId := uuid.New()

	inserted, err := db.Write(dbCtx.Db, `
		INSERT INTO users VALUES($1, $2, $3, NOW())
	`, newId, data.Username, hashedPassword+":"+string(rawSalt))

	if err != nil {
		fmt.Println("Error:", err)
		return utils.GetErrorResponse(http.StatusInternalServerError, "Failed to insert new user")
	}

	if inserted == 0 {
		return utils.GetErrorResponse(http.StatusInternalServerError, "Failed to insert new user")
	}

	return utils.GetSuccessResponse(newId)
}
