package controllers

import (
	"fmt"
	"net/http"

	"go-auth/src/ctx"
	"go-auth/src/dtos"
	"go-auth/src/middlewares"

	"github.com/gorilla/mux"
)

func Home(r *http.Request) dtos.Response {
	return getSuccessResponse("ok")
}

func Ping(r *http.Request) dtos.Response {
	userId := r.Context().Value(ctx.TokenPayloadCtxKey)
	fmt.Println("Data:", userId)
	return getErrorResponse(http.StatusBadRequest, userId.(string))
}

var routes = []Route{
	{
		Name:        "Home",
		Method:      "GET",
		Pattern:     "/",
		HandlerFunc: Home,
	},
	{
		Name:        "Ping",
		Method:      "GET",
		Pattern:     "/ping",
		HandlerFunc: Ping,
		Middlewares: []mux.MiddlewareFunc{
			middlewares.Authorization,
		},
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
	{
		Name:        "Test",
		Method:      "POST",
		Pattern:     "/auth/test",
		HandlerFunc: Test,
	},
}
