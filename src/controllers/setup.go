package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go-auth/src/ctx"
	"go-auth/src/middlewares"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	Middlewares []mux.MiddlewareFunc
}

type ResponseDTO struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func sendSuccessResponse(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Add("content-type", "application/json")
	response := ResponseDTO{
		Status:  http.StatusOK,
		Message: "",
		Success: true,
		Data:    data,
	}
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println("Failed to send response", err)
	}
	reqId := r.Context().Value(ctx.ReqIdCtxKey)
	fmt.Printf("REQUEST_SUCCESS %s : %+v\n", reqId, w)
}

func sendErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	w.Header().Add("content-type", "application/json")
	response := ResponseDTO{
		Status:  statusCode,
		Message: message,
		Success: false,
		Data:    nil,
	}
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println("Failed to send response", err)
	}
	reqId := r.Context().Value(ctx.ReqIdCtxKey)
	fmt.Printf("REQUEST_ERROR %s : %s  |  %+v\n", reqId, message, w)
}

const (
	API_BASE_URL = "/api"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.Use(middlewares.AccessLogger)

	for _, route := range routes {
		var handler http.Handler = route.HandlerFunc
		sub := router.
			Methods(route.Method).Subrouter()

		sub.PathPrefix(API_BASE_URL).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
		for _, mid := range route.Middlewares {
			sub.Use(mid)
		}
	}

	return router
}
