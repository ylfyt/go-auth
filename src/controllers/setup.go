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
	"net/http"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

const (
	API_BASE_URL = "/api"
)

type DependencyInfo struct {
	Key string
	Idx int
}

var depMaps map[string]interfaces.DependencyInjection

func validateHandler(handler interface{}) error {
	ref := reflect.TypeOf(handler)

	// Checking for Handler Output
	if ref.Kind() != reflect.Func {
		return errors.New("handler should be a function")
	}
	if ref.NumOut() != 1 {
		return errors.New("number of return value should be 1")
	}
	if ref.Out(0).String() != "dtos.Response" {
		return errors.New("return value should be dto.response")
	}

	// Checking for Handler Input
	numOfStruct := 0
	for i := 0; i < ref.NumIn(); i++ {
		// For Dependecies Input
		if depMaps[ref.In(i).String()] != nil {
			continue
		}

		// For Payload on Request Body
		if ref.In(i).Kind() == reflect.Struct {
			numOfStruct++
		}

		// For Http Request Pointer
		if ref.In(i).Kind() == reflect.Pointer && ref.In(i).String() != "*fiber.Ctx" {
			return errors.New("pointer arg only allowed with type *http.request")
		}
	}

	// Limiting 1 Variable For Handling Request Payload
	if numOfStruct > 1 {
		return errors.New("number of struct arg should be 1")
	}

	return nil
}

func getCallParams(c *fiber.Ctx, refFunc interface{}) ([]reflect.Value, int, []DependencyInfo) {
	refType := reflect.TypeOf(refFunc)
	var argTypes []reflect.Type
	for i := 0; i < refType.NumIn(); i++ {
		argType := refType.In(i)
		argTypes = append(argTypes, argType)
	}

	// Request Params Setup
	paramIdx := 0
	tempParams := c.AllParams()
	var urlParams = make([]string, 0, len(tempParams))
	for key := range tempParams {
		urlParams = append(urlParams, tempParams[key])
	}

	structIdx := -1
	var depIdxs []DependencyInfo
	var callParams []reflect.Value
	for i, v := range argTypes {
		// Dependecies Setup
		if depMaps[v.String()] != nil {
			depIdxs = append(depIdxs, DependencyInfo{
				Key: v.String(),
				Idx: i,
			})
			var temp interface{}
			callParams = append(callParams, reflect.ValueOf(temp))
			continue
		}

		// Applying Http Request Pointer
		if v.Kind() == reflect.Pointer {
			callParams = append(callParams, reflect.ValueOf(c))
			continue
		}

		// Applying Request Body
		if v.Kind() == reflect.Struct {
			jsonString := c.Body()
			temp := reflect.New(v).Interface()
			_ = json.Unmarshal(jsonString, &temp)
			callParams = append(callParams, reflect.ValueOf(temp).Elem())
			structIdx = i
			continue
		}

		// TODO: Applying Request URL Queries

		// Applying Request Params
		if paramIdx < len(urlParams) {
			callParams = append(callParams, reflect.ValueOf(urlParams[paramIdx]))
			paramIdx++
		} else {
			callParams = append(callParams, reflect.ValueOf(""))
		}
	}

	return callParams, structIdx, depIdxs
}

func sendResponse(c *fiber.Ctx, response dtos.Response) {
	err := c.JSON(response)
	if err != nil {
		fmt.Println("Failed to send response", err)
	}
	reqId := c.Context().Value(ctx.ReqIdCtxKey)
	if response.Status == http.StatusOK {
		l.I("[%s] REQUEST SUCCESS", reqId)
		return
	}
	l.I("[%s] REQUEST FAILED with RESPONSE:%+v\n", reqId, response)
}

func NewRouter() *fiber.App {
	app := fiber.New(fiber.Config{
		StrictRouting: true,
	})
	api := app.Group("/api")

	depMaps = make(map[string]interfaces.DependencyInjection)
	var dependencies []interfaces.DependencyInjection

	// Middleware Setup
	// app.Use(middlewares.Cors)
	api.Use(middlewares.AccessLogger)

	// Dependency Injection Setup
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
			// sub := router.
			// Methods(route.Method, "OPTIONS").Subrouter()

			// Applying Middlewares For each Route
			handlers := append(route.Middlewares, func(c *fiber.Ctx) error {
				params, shouldBeValidateIdx, depIdxs := getCallParams(c, fnHandler)

				// Applying Dependecies
				if len(depIdxs) > 0 {
					for _, v := range depIdxs {
						depVal := depMaps[v.Key].Get()
						params[v.Idx] = reflect.ValueOf(depVal)
						defer depMaps[v.Key].Return(depVal)
					}
				}

				// Calling route handler
				if shouldBeValidateIdx == -1 {
					response := reflect.ValueOf(fnHandler).Call(params)[0].Interface().(dtos.Response)
					sendResponse(c, response)
					return nil
				}

				// Applying Validation For Request Payload
				errs := validate(&params[shouldBeValidateIdx])
				if errs != nil {
					sendResponse(c, utils.GetBadRequestResponse("VALIDATION_ERROR", errs...))
					return nil
				}

				// Calling route handler
				response := reflect.ValueOf(fnHandler).Call(params)[0].Interface().(dtos.Response)
				sendResponse(c, response)
				return nil
			})

			api.Add(route.Method, route.Pattern, handlers...)
		}
	}

	return app
}
