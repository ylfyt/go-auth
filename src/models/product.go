package models

import "github.com/google/uuid"

type Product struct {
	Id          uuid.UUID `json:"id"`
	CreatedAt   string    `json:"createdAt"`
	UpdatedAt   string    `json:"updatedAt"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
}
