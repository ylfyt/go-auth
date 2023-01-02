package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-auth/src/ctx"
	"go-auth/src/dtos"
	"go-auth/src/interfaces"
	"go-auth/src/l"
	"go-auth/src/middlewares"
	"go-auth/src/services"
	"go-auth/src/utils"
	"io"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
)

const (
	API_BASE_URL = "/api"
)

type DepInfo struct {
	Key string
	Idx int
}

var depMaps map[string]interfaces.DependencyInjection

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
		if depMaps[ref.In(i).String()] != nil {
			continue
		}
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

func getCallParams(r *http.Request, refFunc interface{}) ([]reflect.Value, int, []DepInfo) {
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
	var depIdxs []DepInfo
	for i, v := range argTypes {
		if depMaps[v.String()] != nil {
			depIdxs = append(depIdxs, DepInfo{
				Key: v.String(),
				Idx: i,
			})
			var temp interface{}
			callParams = append(callParams, reflect.ValueOf(temp))
			continue
		}
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

	return callParams, structIdx, depIdxs
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
		l.I("[%s] REQUEST SUCCESS", reqId)
		return
	}
	l.I("[%s] REQUEST FAILED with RESPONSE:%+v\n", reqId, response)
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.Use(middlewares.Cors)
	router.Use(middlewares.AccessLogger)

	depMaps = make(map[string]interfaces.DependencyInjection)

	var dependencies []interfaces.DependencyInjection
	dependencies = append(dependencies, services.DbContext{})

	for _, dep := range dependencies {
		ref := reflect.TypeOf(dep)
		depMaps[ref.String()] = dep
	}

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
					params, shouldBeValidateIdx, depIdxs := getCallParams(r, fnHandler)
					if len(depIdxs) > 0 {
						for _, v := range depIdxs {
							depVal := depMaps[v.Key].Get()
							params[v.Idx] = reflect.ValueOf(depVal)
							defer depMaps[v.Key].Return(depVal)
						}
					}

					if shouldBeValidateIdx == -1 {
						response := reflect.ValueOf(fnHandler).Call(params)[0].Interface().(dtos.Response)
						sendResponse(w, r, response)
						return
					}

					errs := validate(&params[shouldBeValidateIdx])

					if errs != nil {
						sendResponse(w, r, utils.GetBadRequestResponse("VALIDATION_ERROR", errs...))
						return
					}

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
