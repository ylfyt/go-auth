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
		// For Dependecies Input
		if app.dependencies[ref.In(i).String()] != nil {
			continue
		}

		// For Http Request Pointer
		if ref.In(i).Kind() == reflect.Pointer {
			if ref.In(i).String() != "*fiber.Ctx" {
				return errors.New("pointer arg only allowed with type *fiber.ctx, or from dependency services")
			}
			continue
		}

		// For Payload on Request Body
		if ref.In(i).Kind() == reflect.Struct {
			numOfStruct++
			continue
		}

		if ref.In(i).Kind() != reflect.String {
			return errors.New("arg for request param only allowed with type 'string'")
		}
	}

	// Limiting 1 Variable For Handling Request Payload
	if numOfStruct > 1 {
		return errors.New("number of struct arg should be 1")
	}

	return nil
}

func (app *App) getCallParams(c *fiber.Ctx, refFunc interface{}) ([]reflect.Value, int) {
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
	var callParams []reflect.Value
	for i, v := range argTypes {
		// Dependecies Setup
		if app.dependencies[v.String()] != nil {
			callParams = append(callParams, reflect.ValueOf(app.dependencies[v.String()]))
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
		fn := v.HandlerFunc
		handlers := append(v.Middlewares, func(c *fiber.Ctx) error {
			params, shouldBeValidateIdx := app.getCallParams(c, fn)

			// Calling route handler
			if shouldBeValidateIdx == -1 {
				response := reflect.ValueOf(fn).Call(params)[0].Interface().(ResponseDto)
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
			response := reflect.ValueOf(fn).Call(params)[0].Interface().(ResponseDto)
			sendResponse(c, response)
			return nil
		})
		
		fmt.Printf("API SETUP: %s | %s\n", v.Path, v.Method)
		app.router.Add(v.Method, v.Path, handlers...)
	}
}
