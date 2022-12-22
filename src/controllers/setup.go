package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-auth/src/ctx"
	"go-auth/src/dtos"
	"go-auth/src/middlewares"
	"io"
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
		if ref.In(i).Kind() == reflect.Pointer && ref.In(i).String() != "*http.Request" {
			return errors.New("pointer arg only allowed with type *http.request")
		}
	}
	if numOfStruct > 1 {
		return errors.New("number of struct arg should be 1")
	}

	return nil
}

func getCallParams(r *http.Request, refFunc interface{}) ([]reflect.Value, int) {
	refType := reflect.TypeOf(refFunc)
	var argTypes []reflect.Type
	for i := 0; i < refType.NumIn(); i++ {
		argType := refType.In(i)
		argTypes = append(argTypes, argType)
	}

	idx := 0
	tempParams := mux.Vars(r)
	var urlParams = make([]string, 0, len(tempParams))
	for key := range tempParams {
		urlParams = append(urlParams, tempParams[key])
	}

	structIdx := -1
	var callParams []reflect.Value
	for i, v := range argTypes {
		if v.Kind() == reflect.Pointer {
			callParams = append(callParams, reflect.ValueOf(r))
			continue
		}
		if v.Kind() == reflect.Struct {
			jsonString, _ := io.ReadAll(r.Body)
			temp := reflect.New(v).Interface()
			_ = json.Unmarshal([]byte(jsonString), &temp)
			callParams = append(callParams, reflect.ValueOf(temp).Elem())
			structIdx = i
			continue
		}

		if idx < len(urlParams) {
			callParams = append(callParams, reflect.ValueOf(urlParams[idx]))
			idx++
		} else {
			callParams = append(callParams, reflect.ValueOf(""))
		}
	}

	return callParams, structIdx
}

func sendResponse(w http.ResponseWriter, r *http.Request, response dtos.Response) {
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
				fmt.Print(route.Name + " handler | ")
				panic(err)
			}

			fmt.Printf("API SETUP: %s | m:%d | %s\n", route.Pattern, len(route.Middlewares), route.Method)
			sub := router.
				Methods(route.Method, "OPTIONS").Subrouter()

			sub.PathPrefix(API_BASE_URL).
				Path(route.Pattern).
				Name(route.Name).
				HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					params, shouldBeValidateIdx := getCallParams(r, fnHandler)
					if shouldBeValidateIdx == -1 {
						response := reflect.ValueOf(fnHandler).Call(params)[0].Interface().(dtos.Response)
						sendResponse(w, r, response)
						return
					}

					// TODO: Validate payload
					response := reflect.ValueOf(fnHandler).Call(params)[0].Interface().(dtos.Response)
					sendResponse(w, r, response)
				})
			for _, mid := range route.Middlewares {
				sub.Use(mid)
			}
		}
	}

	return router
}
