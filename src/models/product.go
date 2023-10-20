package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	Id          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Description string    `json:"description" col:"desc"`
	Price       int64     `json:"price"`
}
