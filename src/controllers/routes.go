package controllers

import (
	"fmt"
	"net/http"

	"go-auth/src/controllers/auth"
	"go-auth/src/ctx"
	"go-auth/src/dtos"
	"go-auth/src/middlewares"
	"go-auth/src/models"
	"go-auth/src/utils"

	"github.com/gorilla/mux"
)

func Home(r *http.Request) dtos.Response {
	claims := r.Context().Value(ctx.UserClaimsCtxKey).(models.AccessClaims)
	fmt.Printf("Data: %+v\n", claims)
	return utils.GetSuccessResponse(claims)
}

func Ping(r *http.Request) dtos.Response {
	return utils.GetSuccessResponse("ok")
}

var routes = []Route{
	{
		Name:        "Home",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: Home,
		Middlewares: []mux.MiddlewareFunc{
			middlewares.Authorization,
		},
	},
	{
		Name:        "Ping",
		Method:      "GET",
		Pattern:     "/ping",
		HandlerFunc: Ping,
	},
	{
		Name:        "AuthRegister",
		Method:      "POST",
		Pattern:     "/auth/register",
		HandlerFunc: auth.Register,
	},
	{
		Name:        "AuthLogin",
		Method:      "POST",
		Pattern:     "/auth/login",
		HandlerFunc: auth.Login,
	},
	{
		Name:        "AuthRefresh",
		Method:      "POST",
		Pattern:     "/auth/refresh",
		HandlerFunc: auth.RefreshToken,
	},
	{
		Name:        "AuthLogout",
		Method:      "POST",
		Pattern:     "/auth/logout",
		HandlerFunc: auth.Logout,
	},
}
