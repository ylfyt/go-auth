package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
	Password  string
}
