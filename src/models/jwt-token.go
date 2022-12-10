package models

import "github.com/google/uuid"

type JwtToken struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt string    `json:"updatedAt"`
	UserId    string    `json:"userId"`
}
