package utils

import (
	"go-auth/src/meta"
	"net/http"
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

func GetErrorResponse(statusCode int, message string, errors ...meta.Error) meta.ResponseDto {
	err := errors
	if err == nil {
		err = []meta.Error{}
	}
	return meta.ResponseDto{
		Status:  statusCode,
		Message: message,
		Errors:  err,
		Success: false,
		Data:    nil,
	}
}

func GetBadRequestResponse(message string, errors ...meta.Error) meta.ResponseDto {
	return GetErrorResponse(http.StatusBadRequest, message, errors...)
}

func GetInternalErrorResponse(message string, errors ...meta.Error) meta.ResponseDto {
	return GetErrorResponse(http.StatusInternalServerError, message, errors...)
}

func GetUnauthorizedResponse(message string, errors ...meta.Error) meta.ResponseDto {
	return GetErrorResponse(http.StatusUnauthorized, message, errors...)
}
