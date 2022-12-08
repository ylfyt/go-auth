package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go-auth/src/ctx"
	"go-auth/src/middlewares"
)

func Home(r *http.Request) ResponseDTO {
	return getSuccessResponse("ok")
}

func Ping(r *http.Request) ResponseDTO {
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
}
