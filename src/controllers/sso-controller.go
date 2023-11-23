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

	data, err := utils.ParseBody[dtos.SsoLoginPayload](r)
	if err != nil {
		sendBadRequestResponse(w, "Payload is not valid")
		return
	}

	client := ssoClients[data.Client]
	if client == nil {
		sendBadRequestResponse(w, "Client is no valid")
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

	refresh, access, _, err := services.CreateJwtToken(me.config, *user)
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
	exchangeToken := uuid.New().String()
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
