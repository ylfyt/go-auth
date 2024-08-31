package dtos

import "go-auth/src/models"

type LoginResponse struct {
	User  models.User
	Token TokenPayload
}
type FieldError struct {
	Field string
	Tag   string
	Param string
}

type Response struct {
	Code    int
	Message string
	Success bool
	Errors  []FieldError
	Data    any
}
