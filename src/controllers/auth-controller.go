package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go-auth/src/dtos"
	"go-auth/src/models"
	"io"
	"net/http"
)

func Register(r *http.Request) ResponseDTO {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return getErrorResponse(http.StatusBadRequest, "Failed to get request body")
	}

	var data dtos.Register
	err = json.Unmarshal(body, &data)
	if err != nil {
		return getErrorResponse(http.StatusBadRequest, "Failed to get payload")
	}
	user := models.User{
		Id:       uuid.New(),
		Username: data.Username,
		Password: data.Password,
	}
	fmt.Printf("Data: %+v\n", user)
	return getSuccessResponse(user)
}
