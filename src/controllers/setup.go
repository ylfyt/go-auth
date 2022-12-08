package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	// "go-auth/src/middlewares"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(r *http.Request) ResponseDTO
	Middlewares []mux.MiddlewareFunc
}

type ResponseDTO struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func getSuccessResponse(data interface{}) ResponseDTO {
	return ResponseDTO{
		Status:  http.StatusOK,
		Message: "",
		Success: true,
		Data:    data,
	}
}

func getErrorResponse(statusCode int, message string) ResponseDTO {
	return ResponseDTO{
		Status:  statusCode,
		Message: message,
		Success: false,
		Data:    nil,
	}
}

const (
	API_BASE_URL = "/api"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	// router.Use(middlewares.AccessLogger)

	for _, route := range routes {
		fnHandler := route.HandlerFunc
		sub := router.
			Methods(route.Method).Subrouter()

		sub.PathPrefix(API_BASE_URL).
			Path(route.Pattern).
			Name(route.Name).
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				response := fnHandler(r)
				w.Header().Add("content-type", "application/json")
				w.WriteHeader(response.Status)
				err := json.NewEncoder(w).Encode(response)
				if err != nil {
					fmt.Println("Failed to send response", err)
				}
			})
		for _, mid := range route.Middlewares {
			sub.Use(mid)
		}
	}

	return router
}
