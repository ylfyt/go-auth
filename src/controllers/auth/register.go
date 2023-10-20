package auth

import (
	"database/sql"
	"fmt"
	"go-auth/src/db"
	"go-auth/src/dtos"
	"go-auth/src/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/ylfyt/meta/meta"
)

func register(data dtos.Register, dbCtx *sql.DB) meta.ResponseDto {
	exists, err := db.GetRowCount(dbCtx, `
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

	inserted, err := db.Write(dbCtx, `
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
