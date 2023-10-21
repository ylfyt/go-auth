package services

import (
	"go-auth/src/models"
	"go-auth/src/structs"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func createAccessToken(config *structs.EnvConf, user models.User, refreshTokenId uuid.UUID) (string, error) {
	claims := models.AccessClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(config.JwtAccessTokenExpiryTime) * time.Second).Unix(),
		},
		Username: user.Username,
		Jid:      refreshTokenId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(config.JwtAccessTokenSecretKey))
	if err != nil {
		return "", err
	}
	return signed, err
}

func createRefreshToken(config *structs.EnvConf) (string, *uuid.UUID, error) {
	jid := uuid.New()
	claims := models.RefreshClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(config.JwtRefreshTokenExpiryTime) * time.Minute).Unix(),
		},
		Jid: jid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(config.JwtRefreshTokenSecretKey))
	if err != nil {
		return "", nil, err
	}

	return signed, &jid, nil
}

func CreateJwtToken(config *structs.EnvConf, user models.User) (string, string, *uuid.UUID, error) {
	refreshToken, jid, err := createRefreshToken(config)
	if err != nil {
		return "", "", nil, err
	}
	accessToken, err := createAccessToken(config, user, *jid)
	if err != nil {
		return "", "", nil, err
	}

	return refreshToken, accessToken, jid, nil
}

func VerifyAccessToken(config *structs.EnvConf, token string) (bool, models.AccessClaims) {
	claims := models.AccessClaims{}
	tkn, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.JwtAccessTokenSecretKey), nil
	})

	if err != nil {
		return false, models.AccessClaims{}
	}

	if !tkn.Valid {
		return false, models.AccessClaims{}
	}

	return true, claims
}

func VerifyRefreshToken(config *structs.EnvConf, token string) (bool, *uuid.UUID) {
	claims := models.RefreshClaims{}
	tkn, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.JwtRefreshTokenSecretKey), nil
	})

	if err != nil {
		return false, nil
	}

	if !tkn.Valid {
		return false, nil
	}

	return true, &claims.Jid
}
