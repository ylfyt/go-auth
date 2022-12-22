package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-auth/src/ctx"
	"go-auth/src/middlewares"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
)

const (
	API_BASE_URL = "/api"
)

func validateHandler(handler interface{}) error {
	ref := reflect.TypeOf(handler)
	if ref.Kind() != reflect.Func {
		return errors.New("handler should be a function")
	}
	if ref.NumOut() != 1 {
		return errors.New("number of return value should be 1")
	}
	if ref.Out(0).String() != "dtos.Response" {
		return errors.New("return value should be dto.response")
	}

	numOfStruct := 0
	for i := 0; i < ref.NumIn(); i++ {
		if ref.In(i).Kind() == reflect.Struct {
			numOfStruct++
		}
		if ref.In(i).Kind() == reflect.Pointer {
			return errors.New("arg cannot a pointer")
		}
	}
	if numOfStruct > 1 {
		return errors.New("number of struct arg should be 1")
	}

	return nil
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.Use(middlewares.Cors)
	router.Use(middlewares.AccessLogger)

	for _, routes := range appRoutes {
		for _, route := range routes {
			fnHandler := route.HandlerFunc
			err := validateHandler(fnHandler)
			if err != nil {
				panic(err)
			}
			fmt.Printf("API SETUP: %s | m:%d | %s\n", route.Pattern, len(route.Middlewares), route.Method)
			sub := router.
				Methods(route.Method, "OPTIONS").Subrouter()

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
					reqId := r.Context().Value(ctx.ReqIdCtxKey)
					if response.Status == http.StatusOK {
						fmt.Printf("[%s] REQUEST SUCCESS\n", reqId)
						return
					}
					fmt.Printf("[%s] REQUEST FAILED with RESPONSE:%+v\n", reqId, response)
				})
			for _, mid := range route.Middlewares {
				sub.Use(mid)
			}
		}
	}

	return router
}
