package models

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type AccessClaims struct {
	jwt.StandardClaims
	Username string    `json:"username"`
	Jid      uuid.UUID `json:"jid"`
}

type RefreshClaims struct {
	jwt.StandardClaims
	Jid uuid.UUID `json:"jid"`
}
