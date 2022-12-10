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
	}
}

func GetErrorResponse(statusCode int, message string) dtos.Response {
	return dtos.Response{
		Status:  statusCode,
		Message: message,
		Success: false,
		Data:    nil,
	}
}