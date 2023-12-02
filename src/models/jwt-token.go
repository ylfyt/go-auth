package models

import (
	"time"
)

type JwtToken struct {
	Id        int64      `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	UserId    int        `db:"user_id"`
}
