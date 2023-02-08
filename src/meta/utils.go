package meta

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func (app *App) validateHandler(handler interface{}) error {
	ref := reflect.TypeOf(handler)

	// Checking for Handler Output
	if ref.Kind() != reflect.Func {
		return errors.New("handler should be a function")
	}
	if ref.NumOut() != 1 {
		return errors.New("return value should be responsedto")
	}

	if ref.Out(0).String() != "meta.ResponseDto" {
		return errors.New("return value should be responsedto")
	}

	// Checking for Handler Input
	numOfStruct := 0
	for i := 0; i < ref.NumIn(); i++ {
		// TODO: Dependecies
		// // For Dependecies Input
		// if depMaps[ref.In(i).String()] != nil {
		// 	continue
		// }

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

func getCallParams(c *fiber.Ctx, refFunc interface{}) ([]reflect.Value, int) {
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
	// var depIdxs []DependencyInfo
	var callParams []reflect.Value
	for i, v := range argTypes {
		// Dependecies Setup
		// if depMaps[v.String()] != nil {
		// 	depIdxs = append(depIdxs, DependencyInfo{
		// 		Key: v.String(),
		// 		Idx: i,
		// 	})
		// 	var temp interface{}
		// 	callParams = append(callParams, reflect.ValueOf(temp))
		// 	continue
		// }

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

	return callParams, structIdx
}

func validate(data *reflect.Value) []Error {
	val := validator.New()
	err := val.Struct(data.Interface())
	if err == nil {
		return nil
	}
	if _, ok := err.(*validator.InvalidValidationError); ok {
		fmt.Println("Err", err)
		return nil
	}

	var validationErrors []Error

	for _, err := range err.(validator.ValidationErrors) {

		validationErrors = append(validationErrors, Error{
			Field: err.Field(),
			Tag:   err.ActualTag(),
			Param: err.Param(),
		})
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return validationErrors
}

func sendResponse(c *fiber.Ctx, response ResponseDto) {
	err := c.JSON(response)
	if err != nil {
		fmt.Println("Failed to send response", err)
	}
}

func (app *App) setup() {
	for _, v := range app.endPoints {
		fmt.Printf("API SETUP: %s | %s\n", v.Path, v.Method)
		var handlers []func(c *fiber.Ctx) error = []func(c *fiber.Ctx) error{
			func(c *fiber.Ctx) error {
				params, shouldBeValidateIdx := getCallParams(c, v.HandlerFunc)

				// Calling route handler
				if shouldBeValidateIdx == -1 {
					response := reflect.ValueOf(v.HandlerFunc).Call(params)[0].Interface().(ResponseDto)
					sendResponse(c, response)
					return nil
				}

				// Applying Validation For Request Payload
				errs := validate(&params[shouldBeValidateIdx])
				if errs != nil {
					sendResponse(c, ResponseDto{
						Status:  fiber.ErrBadRequest.Code,
						Message: "VALIDATION_ERROR",
						Errors:  errs,
						Success: false,
						Data:    nil,
					})
					return nil
				}

				// Calling route handler
				response := reflect.ValueOf(v.HandlerFunc).Call(params)[0].Interface().(ResponseDto)
				sendResponse(c, response)
				return nil
			},
		}

		app.router.Add(v.Method, v.Path, handlers...)
	}
}
