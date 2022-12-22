package controllers

import (
	"fmt"
	"go-auth/src/dtos"
	"reflect"

	"github.com/go-playground/validator/v10"
)

func validate(data *reflect.Value) []dtos.Error {
	val := validator.New()
	err := val.Struct(data.Interface())
	if err == nil {
		return nil
	}
	if _, ok := err.(*validator.InvalidValidationError); ok {
		fmt.Println("Err", err)
		return nil
	}

	var validationErrors []dtos.Error

	for _, err := range err.(validator.ValidationErrors) {

		validationErrors = append(validationErrors, dtos.Error{
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
