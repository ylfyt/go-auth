package dtos

import "go-auth/src/models"

type LoginResponse struct {
	User  models.User  `json:"user"`
	Token TokenPayload `json:"token"`
}
