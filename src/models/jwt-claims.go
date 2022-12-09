package models

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JwtClaims struct {
	jwt.StandardClaims
	Username string    `json:"username"`
	Jid      uuid.UUID `json:"jid"`
}
