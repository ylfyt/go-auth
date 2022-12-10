package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	API_BASE_URL = "/api"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

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
