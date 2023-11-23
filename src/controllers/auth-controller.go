package controllers

import (
	"fmt"
	"go-auth/src/dtos"
	"go-auth/src/models"
	"go-auth/src/services"
	"go-auth/src/utils"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (me *Controller) login(w http.ResponseWriter, r *http.Request) {
	data, err := utils.ParseBody[dtos.Register](r)
	if err != nil {
		fmt.Println("ERR", err)
		sendBadRequestResponse(w, "Payload is not valid")
		return
	}

	var user *models.User
	err = me.db.GetFirst(&user, `
	SELECT * FROM users WHERE username = $1
	`, data.Username)
	if err != nil {
		fmt.Println("ERR", err)
		sendDefaultInternalErrorResponse(w)
		return
	}

	if user == nil {
		sendBadRequestResponse(w, "Username or password is wrong")
		return
	}

	passwordData := strings.Split(user.Password, ":")
	if len(passwordData) != 2 {
		fmt.Println("???")
		sendDefaultInternalErrorResponse(w)
		return
	}
	isValid := utils.VerifyPassword(passwordData[0], data.Password, user.Username, []byte(passwordData[1]))
	if !isValid {
		sendBadRequestResponse(w, "Username or password is wrong")
		return
	}

	refresh, access, jid, err := services.CreateJwtToken(me.config, *user)
	if err != nil {
		fmt.Println("Data:", err)
		sendDefaultInternalErrorResponse(w)
		return
	}
	_, err = me.db.Write(`
		INSERT INTO jwt_tokens(id, user_id, created_at) VALUES($1, $2, NOW())
	`, jid, user.Id)
	if err != nil {
		fmt.Println("Data:", err)
		sendDefaultInternalErrorResponse(w)
		return
	}

	sendSuccessResponse(w, dtos.LoginResponse{
		User: *user,
		Token: dtos.TokenPayload{
			RefreshToken: refresh,
			AccessToken:  access,
			ExpiredIn:    int64(me.config.JwtAccessTokenExpiryTime),
		},
	})
}

func (me *Controller) logoutAll(w http.ResponseWriter, r *http.Request) {
	data, err := utils.ParseBody[dtos.Register](r)
	if err != nil {
		fmt.Println("ERR", err)
		sendBadRequestResponse(w, "Payload is not valid")
	}

	var user *models.User
	err = me.db.GetFirst(&user, `
	SELECT * FROM users WHERE username = $1
	`, data.Username)
	if err != nil {
		sendDefaultInternalErrorResponse(w)
		return
	}

	if user == nil {
		sendBadRequestResponse(w, "Username or password is wrong")
		return
	}

	passwordData := strings.Split(user.Password, ":")
	if len(passwordData) != 2 {
		sendDefaultInternalErrorResponse(w)
		return
	}
	isValid := utils.VerifyPassword(passwordData[0], data.Password, user.Username, []byte(passwordData[1]))
	if !isValid {
		sendBadRequestResponse(w, "Username or password is wrong")
		return
	}

	deleted, err := me.db.Write(`
		DELETE FROM jwt_tokens WHERE user_id = $1
	`, user.Id)
	if err != nil {
		sendDefaultInternalErrorResponse(w)
		return
	}

	if deleted == 0 {
		sendBadRequestResponse(w, "There is no logged in")
		return
	}

	sendSuccessResponse(w)
}

func (me *Controller) logout(w http.ResponseWriter, r *http.Request) {
	data, err := utils.ParseBody[dtos.RefreshPayload](r)
	if err != nil {
		fmt.Println("ERR", err)
		sendBadRequestResponse(w, "Payload is not valid")
		return
	}
	valid, jid := services.VerifyRefreshToken(me.config, data.RefreshToken)
	if !valid {
		sendBadRequestResponse(w, "Token is not valid")
		return
	}

	deleted, err := me.db.Write(`
		DELETE FROM jwt_tokens WHERE id = $1
	`, jid)
	if err != nil {
		fmt.Println("Err", err)
		sendDefaultInternalErrorResponse(w)
		return
	}

	if deleted == 0 {
		sendBadRequestResponse(w, "Token is not found")
		return
	}

	// TODO: Implement Blacklist Token Mechanism

	sendSuccessResponse(w, true)
}

func (me *Controller) refreshToken(w http.ResponseWriter, r *http.Request) {
	data, err := utils.ParseBody[dtos.RefreshPayload](r)
	if err != nil {
		fmt.Println("ERR", err)
		sendBadRequestResponse(w, "Payload is not valid")
		return
	}
	valid, jid := services.VerifyRefreshToken(me.config, data.RefreshToken)
	if !valid {
		sendBadRequestResponse(w, "Token is not valid")
		return
	}

	var token *models.JwtToken
	err = me.db.GetFirst(&token, `
		SELECT * FROM jwt_tokens WHERE id = $1
	`, jid)
	if err != nil {
		fmt.Println("ERROR", err)
		sendDefaultInternalErrorResponse(w)
		return
	}
	if token == nil {
		sendBadRequestResponse(w, "Token is not found")
		return
	}

	var user *models.User
	err = me.db.GetFirst(&user, `
	SELECT * FROM users WHERE id = $1
	`, token.UserId)
	if err != nil {
		sendDefaultInternalErrorResponse(w)
		return
	}
	if user == nil {
		sendBadRequestResponse(w, "User is not found")
		return
	}

	refresh, access, newJid, err := services.CreateJwtToken(me.config, *user)
	if err != nil {
		sendDefaultInternalErrorResponse(w)
		return
	}
	_, err = me.db.Write(`
		INSERT INTO jwt_tokens VALUES($1, $2, NOW())
	`, newJid, user.Id)
	if err != nil {
		sendDefaultInternalErrorResponse(w)
		return
	}

	_, err = me.db.Write(`
	DELETE FROM jwt_tokens WHERE id = $1
	`, jid)
	if err != nil {
		fmt.Println("ERROR", err)
	}

	sendSuccessResponse(w, dtos.LoginResponse{
		User: *user,
		Token: dtos.TokenPayload{
			RefreshToken: refresh,
			AccessToken:  access,
			ExpiredIn:    int64(me.config.JwtAccessTokenExpiryTime),
		},
	})
}

func (me *Controller) register(w http.ResponseWriter, r *http.Request) {
	data, err := utils.ParseBody[dtos.Register](r)
	if err != nil {
		fmt.Println("ERR", err)
		sendBadRequestResponse(w, "Payload is not valid")
		return
	}
	var count int
	err = me.db.ColFirst(&count, `
		SELECT count(*) FROM users WHERE username = $1
	`, data.Username)
	if err != nil {
		fmt.Println("Error", err)
		sendDefaultInternalErrorResponse(w)
		return
	}

	if count != 0 {
		sendBadRequestResponse(w, "Username already exist")
		return
	}

	rawSalt := utils.GenerateRawSalt()
	realSalt := utils.GetRealSalt(rawSalt, data.Username)
	hashedPassword := utils.HashPassword(data.Password, realSalt)

	newId := uuid.New()
	_, err = me.db.Write(`INSERT INTO users(id, username, password, created_at) VALUES($1, $2, $3, NOW())`, newId, data.Username, hashedPassword+":"+string(rawSalt))
	if err != nil {
		fmt.Println("Error:", err)
		sendDefaultInternalErrorResponse(w)
		return
	}

	sendSuccessResponse(w, newId)
}
