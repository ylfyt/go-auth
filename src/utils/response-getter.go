package utils

import (
	"go-auth/src/dtos"
	"net/http"
)

func GetSuccessResponse(data interface{}) dtos.Response {
	return dtos.Response{
		Status:  http.StatusOK,
		Message: "",
		Success: true,
		Data:    data,
		Errors:  nil,
	}
}

func GetErrorResponse(statusCode int, message string, errors ...dtos.Error) dtos.Response {
	err := errors
	if err == nil {
		err = []dtos.Error{}
	}
	return dtos.Response{
		Status:  statusCode,
		Message: message,
		Errors:  err,
		Success: false,
		Data:    nil,
	}
}

func GetBadRequestResponse(message string, errors ...dtos.Error) dtos.Response {
	return GetErrorResponse(http.StatusBadRequest, message, errors...)
}

func GetInternalErrorResponse(message string, errors ...dtos.Error) dtos.Response {
	return GetErrorResponse(http.StatusInternalServerError, message, errors...)
}

func GetUnauthorizedResponse(message string, errors ...dtos.Error) dtos.Response {
	return GetErrorResponse(http.StatusUnauthorized, message, errors...)
}
