package models

import (
	"time"

	"github.com/google/uuid"
)

type JwtToken struct {
	Id        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserId    uuid.UUID
}
