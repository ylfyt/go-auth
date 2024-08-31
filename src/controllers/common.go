package controllers

import (
	"fmt"
	"go-auth/src/dtos"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

func GetErrorResponse(statusCode int, message string, errors ...dtos.FieldError) dtos.Response {
	return dtos.Response{
		Code:  statusCode,
		Message: message,
		Errors:  errors,
		Success: false,
		Data:    nil,
	}
}

func sendResponse(w http.ResponseWriter, response dtos.Response) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(response.Code)
	err := jsoniter.ConfigCompatibleWithStandardLibrary.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println("ERR", err)
	}
}

func sendBadRequestResponse(w http.ResponseWriter, msg string) {
	response := GetErrorResponse(http.StatusBadRequest, msg)
	sendResponse(w, response)
}

func sendInternalErrorResponse(w http.ResponseWriter, msg string) {
	response := GetErrorResponse(http.StatusInternalServerError, msg)
	sendResponse(w, response)
}

func sendDefaultInternalErrorResponse(w http.ResponseWriter) {
	sendInternalErrorResponse(w, "Internal Server Error")
}

func GetSuccessResponse(data any) dtos.Response {
	return dtos.Response{
		Code:  http.StatusOK,
		Message: "",
		Success: true,
		Data:    data,
		Errors:  nil,
	}
}

func sendSuccessResponse(w http.ResponseWriter, data ...any) {
	if len(data) == 0 {
		sendResponse(w, GetSuccessResponse(nil))
		return
	}
	sendResponse(w, GetSuccessResponse(data[0]))
}
