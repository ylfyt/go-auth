package models

import (
	"time"
)

type User struct {
	Id        int        `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	Username  string     `db:"username"`
	Password  string     `db:"password" json:"-"`
}
