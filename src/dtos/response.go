package dtos

import "go-auth/src/models"

type LoginResponse struct {
	User  models.User  `json:"user"`
	Token TokenPayload `json:"token"`
}
type FieldError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

type Response struct {
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Success bool         `json:"success"`
	Errors  []FieldError `json:"errors"`
	Data    any          `json:"data"`
}
