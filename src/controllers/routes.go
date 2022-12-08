package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go-auth/src/ctx"
	"go-auth/src/middlewares"
)

func Home(w http.ResponseWriter, r *http.Request) {
	sendSuccessResponse(w, r, "ok")
}

func Ping(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(ctx.TokenPayloadCtxKey)
	fmt.Println("Data:", userId)
	sendErrorResponse(w, r, http.StatusBadRequest, "Bad Request")
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
