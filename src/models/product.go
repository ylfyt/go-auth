package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	Id          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	Description string
	Price       int64
}
