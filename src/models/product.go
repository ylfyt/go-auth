package models

import (
	"time"
)

type Product struct {
	Id          int64
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	Name        string
	Description string
	Price       int64
}
