package models

import "github.com/google/uuid"

type JwtToken struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	UserId    string    `json:"user_id"`
}
