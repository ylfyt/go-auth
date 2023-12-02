package models

import (
	"github.com/golang-jwt/jwt"
)

type AccessClaims struct {
	jwt.StandardClaims
	Username string `json:"un"`
	Jid      int64  `json:"jid"`
}

type RefreshClaims struct {
	jwt.StandardClaims
	Jid int64 `json:"jid"`
}
