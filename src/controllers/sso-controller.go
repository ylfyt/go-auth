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
	"time"
)

type SsoClient struct {
	Id          string
	CallbackUrl string
}

func (me *Controller) ssoLogin(w http.ResponseWriter, r *http.Request) {
	ssoClients := map[string]*SsoClient{
		"123": {
			Id:          "123",
			CallbackUrl: "https://www.google.com/hehe",
		},
	}

	data := utils.GetBodyContext[dtos.SsoLoginPayload](r)

	client := ssoClients[data.Client]
	if client == nil {
		sendBadRequestResponse(w, "Client is no valid")
		return
	}

	var user models.User
	err := me.db.Get(&user, `
	SELECT * FROM users WHERE username = ?
	`, data.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			sendBadRequestResponse(w, "User is not found")
			return
		}
		fmt.Println("ERR", err)
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

	refresh, access, _, err := services.CreateJwtToken(me.config, user)
	if err != nil {
		fmt.Println("Data:", err)
		sendDefaultInternalErrorResponse(w)
		return
	}

	token := dtos.TokenPayload{
		RefreshToken: refresh,
		AccessToken:  access,
		ExpiredIn:    int64(me.config.JwtAccessTokenExpiryTime),
	}
	exchangeToken := time.Now().Unix()
	err = me.ssoTokenService.Store(exchangeToken, token)
	if err != nil {
		fmt.Println("ERR", err)
		sendDefaultInternalErrorResponse(w)
		return
	}

	sendSuccessResponse(w, dtos.SsoLoginResponse{
		Callback: client.CallbackUrl,
		Exchange: exchangeToken,
	})
}
