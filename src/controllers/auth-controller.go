package controllers

import (
	"database/sql"
	"fmt"
	"go-auth/src/dtos"
	"go-auth/src/models"
	"go-auth/src/services"
	"go-auth/src/utils"
	"net/http"
	"strings"
)

func (me *Controller) login(w http.ResponseWriter, r *http.Request) {
	data := utils.GetBodyContext[dtos.Register](r)

	var user models.User
	err := me.db.Get(&user, `SELECT * FROM users WHERE username = ?`, data.Username)
	if err != nil {
		fmt.Println("ERR", err)
		if err == sql.ErrNoRows {
			sendBadRequestResponse(w, "Username or password is wrong")
			return
		}
		sendDefaultInternalErrorResponse(w)
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

	refresh, access, jid, err := services.CreateJwtToken(me.config, user)
	if err != nil {
		fmt.Println("Data:", err)
		sendDefaultInternalErrorResponse(w)
		return
	}
	_, err = me.db.Exec(`
		INSERT INTO jwt_tokens(id, user_id, created_at) VALUES(?, ?, CURRENT_TIMESTAMP)
	`, jid, user.Id)
	if err != nil {
		fmt.Println("Data:", err)
		sendDefaultInternalErrorResponse(w)
		return
	}

	sendSuccessResponse(w, dtos.LoginResponse{
		User: user,
		Token: dtos.TokenPayload{
			RefreshToken: refresh,
			AccessToken:  access,
			ExpiredIn:    int64(me.config.JwtAccessTokenExpiryTime),
		},
	})
}

func (me *Controller) logoutAll(w http.ResponseWriter, r *http.Request) {
	data := utils.GetBodyContext[dtos.Register](r)

	var user models.User
	err := me.db.Get(&user, `
		SELECT * FROM users WHERE username = ?
	`, data.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			sendBadRequestResponse(w, "Username or password is wrong")
			return
		}
		sendDefaultInternalErrorResponse(w)
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

	result, err := me.db.Exec(`
		DELETE FROM jwt_tokens WHERE user_id = ?
	`, user.Id)
	if err != nil {
		sendDefaultInternalErrorResponse(w)
		return
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		sendBadRequestResponse(w, "There is no logged in")
		return
	}

	sendSuccessResponse(w)
}

func (me *Controller) logout(w http.ResponseWriter, r *http.Request) {
	data := utils.GetBodyContext[dtos.RefreshPayload](r)

	valid, jid := services.VerifyRefreshToken(me.config, data.RefreshToken)
	if !valid {
		sendBadRequestResponse(w, "Token is not valid")
		return
	}

	result, err := me.db.Exec(`
		DELETE FROM jwt_tokens WHERE id = ?
	`, jid)
	if err != nil {
		fmt.Println("Err", err)
		sendDefaultInternalErrorResponse(w)
		return
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		sendBadRequestResponse(w, "Token is not found")
		return
	}

	// TODO: Implement Blacklist Token Mechanism

	sendSuccessResponse(w, true)
}

func (me *Controller) refreshToken(w http.ResponseWriter, r *http.Request) {
	data := utils.GetBodyContext[dtos.RefreshPayload](r)

	valid, jid := services.VerifyRefreshToken(me.config, data.RefreshToken)
	if !valid {
		sendBadRequestResponse(w, "Token is not valid")
		return
	}

	var token models.JwtToken
	err := me.db.Get(&token, `
		SELECT * FROM jwt_tokens WHERE id = ?
	`, jid)
	if err != nil {
		fmt.Println("ERROR", err)
		if err == sql.ErrNoRows {
			sendBadRequestResponse(w, "Token is not found")
			return
		}
		sendDefaultInternalErrorResponse(w)
		return
	}

	var user models.User
	err = me.db.Get(&user, `
	SELECT * FROM users WHERE id = ?
	`, token.UserId)
	if err != nil {
		if err != sql.ErrNoRows {
			sendBadRequestResponse(w, "User is not found")
			return
		}
		fmt.Println("ERR", err)
		sendDefaultInternalErrorResponse(w)
		return
	}

	refresh, access, newJid, err := services.CreateJwtToken(me.config, user)
	if err != nil {
		sendDefaultInternalErrorResponse(w)
		return
	}
	_, err = me.db.Exec(`
		INSERT INTO jwt_tokens (id, created_at, user_id) VALUES(?, CURRENT_TIMESTAMP, ?)
	`, newJid, user.Id)
	if err != nil {
		fmt.Println("ERR", err)
		sendDefaultInternalErrorResponse(w)
		return
	}

	_, err = me.db.Exec(`
		DELETE FROM jwt_tokens WHERE id = ?
	`, jid)
	if err != nil {
		fmt.Println("ERROR", err)
	}

	sendSuccessResponse(w, dtos.LoginResponse{
		User: user,
		Token: dtos.TokenPayload{
			RefreshToken: refresh,
			AccessToken:  access,
			ExpiredIn:    int64(me.config.JwtAccessTokenExpiryTime),
		},
	})
}

func (me *Controller) register(w http.ResponseWriter, r *http.Request) {
	data := utils.GetBodyContext[dtos.Register](r)

	var count int
	err := me.db.Get(&count, `
		SELECT count(*) FROM users WHERE username = ?
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

	res, err := me.db.Exec(`INSERT INTO users(username, password, created_at) VALUES(?, ?, CURRENT_TIMESTAMP)`, data.Username, hashedPassword+":"+string(rawSalt))
	if err != nil {
		fmt.Println("Error:", err)
		sendDefaultInternalErrorResponse(w)
		return
	}
	newId, _ := res.LastInsertId()
	sendSuccessResponse(w, newId)
}
