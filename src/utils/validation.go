package utils

import (
	"go-auth/src/dtos"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func Validation[T any](data *T) []dtos.FieldError {
	validate := validator.New()
	err := validate.Struct(data)
	if err == nil {
		return nil
	}
	if _, ok := err.(validator.ValidationErrors); !ok {
		return nil
	}

	var fieldErrors []dtos.FieldError
	for _, err := range err.(validator.ValidationErrors) {
		fieldName := err.StructField()
		field, _ := reflect.TypeOf(data).Elem().FieldByName(fieldName)
		if jsonTag, ok := field.Tag.Lookup("json"); ok {
			name := strings.SplitN(jsonTag, ",", 2)[0]
			if name != "-" {
				fieldName = name
			}
		} else {
			first := fieldName[0]
			if 'A' <= first && first <= 'Z' {
				first += 'a' - 'A'
			}
			strBytes := []byte(fieldName)
			strBytes[0] = first
			fieldName = string(strBytes)
		}
		fieldErrors = append(fieldErrors, dtos.FieldError{
			Field: fieldName,
			Tag:   err.Tag(),
			Param: err.Param(),
		})
	}
	return fieldErrors
}
