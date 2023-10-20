package utils

import (
	"net/http"

	"github.com/ylfyt/meta/meta"
)

func GetSuccessResponse(data interface{}) meta.ResponseDto {
	return meta.ResponseDto{
		Status:  http.StatusOK,
		Message: "",
		Success: true,
		Data:    data,
		Errors:  nil,
	}
}

func GetErrorResponse(statusCode int, message string, errors ...meta.FieldError) meta.ResponseDto {
	err := errors
	if err == nil {
		err = []meta.FieldError{}
	}
	return meta.ResponseDto{
		Status:  statusCode,
		Message: message,
		Errors:  err,
		Success: false,
		Data:    nil,
	}
}

func GetBadRequestResponse(message string, errors ...meta.FieldError) meta.ResponseDto {
	return GetErrorResponse(http.StatusBadRequest, message, errors...)
}

func GetInternalErrorResponse(message string, errors ...meta.FieldError) meta.ResponseDto {
	return GetErrorResponse(http.StatusInternalServerError, message, errors...)
}

func GetUnauthorizedResponse(message string, errors ...meta.FieldError) meta.ResponseDto {
	return GetErrorResponse(http.StatusUnauthorized, message, errors...)
}
