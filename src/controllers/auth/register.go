package auth

import (
	"fmt"
	"go-auth/src/dtos"
	"go-auth/src/utils"

	"github.com/google/uuid"
	"github.com/ylfyt/meta/meta"
)

func (me *AuthController) register(data dtos.Register) meta.ResponseDto {
	var count int
	err := me.db.ColFirst(&count, `
		SELECT count(*) FROM users WHERE username = $1
	`, data.Username)
	if err != nil {
		fmt.Println("Error", err)
		return utils.GetInternalErrorResponse("Something wrong!")
	}

	if count != 0 {
		return utils.GetBadRequestResponse("Username already exist")
	}

	rawSalt := utils.GenerateRawSalt()
	realSalt := utils.GetRealSalt(rawSalt, data.Username)
	hashedPassword := utils.HashPassword(data.Password, realSalt)

	newId := uuid.New()
	_, err = me.db.Write(`INSERT INTO users VALUES($1, $2, $3, NOW())`, newId, data.Username, hashedPassword+":"+string(rawSalt))
	if err != nil {
		fmt.Println("Error:", err)
		return utils.GetInternalErrorResponse("Failed to insert new user")
	}

	return utils.GetSuccessResponse(newId)
}
