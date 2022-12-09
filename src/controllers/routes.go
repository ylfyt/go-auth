package controllers

import (
	"fmt"
	"net/http"

	"go-auth/src/ctx"
	"go-auth/src/dtos"
	"go-auth/src/middlewares"
	"go-auth/src/models"

	"github.com/gorilla/mux"
)

func Home(r *http.Request) dtos.Response {
	claims := r.Context().Value(ctx.UserClaimsCtxKey).(models.AccessClaims)
	fmt.Printf("Data: %+v\n", claims)
	return getSuccessResponse(claims)
}

func Ping(r *http.Request) dtos.Response {
	return getSuccessResponse("ok")
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
		HandlerFunc: Register,
	},
	{
		Name:        "AuthLogin",
		Method:      "POST",
		Pattern:     "/auth/login",
		HandlerFunc: Login,
	},
}
