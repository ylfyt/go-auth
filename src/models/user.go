package models

import (
	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt string    `json:"updatedAt"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
}
