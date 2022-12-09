package services

import (
	"go-auth/src/config"
	"go-auth/src/models"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func CreateAccessToken(user models.User, refreshTokenId uuid.UUID) (string, error) {
	claims := models.JwtClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(config.JWT_ACCESS_TOKEN_EXPIRY_TIME) * time.Second).Unix(),
		},
		Username: user.Username,
		Jid:      refreshTokenId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(config.JWT_ACCESS_TOKEN_SECRET_KEY))
	if err != nil {
		return "", err
	}
	return signed, err
}

func VerifyAccessToken(token string) (*models.JwtClaims, error) {
	claims := models.JwtClaims{}
	tkn, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.JWT_ACCESS_TOKEN_SECRET_KEY), nil
	})

	if !tkn.Valid {
		return nil, err
	}
	return &claims, nil
}
