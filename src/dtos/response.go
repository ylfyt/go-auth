package dtos

import "go-auth/src/models"

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Errors  []Error     `json:"errors"`
	Data    interface{} `json:"data"`
}

type LoginResponse struct {
	User  models.User  `json:"user"`
	Token TokenPayload `json:"token"`
}
