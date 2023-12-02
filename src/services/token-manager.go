package services

import (
	"go-auth/src/dtos"
	"go-auth/src/models"
	"go-auth/src/shared"
	"time"

	"github.com/golang-jwt/jwt"
)

func createAccessToken(config *shared.EnvConf, user models.User, refreshTokenId int64) (string, error) {
	claims := dtos.AccessClaims{
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

func createRefreshToken(config *shared.EnvConf) (string, int64, error) {
	jid := time.Now().Unix()
	claims := dtos.RefreshClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(config.JwtRefreshTokenExpiryTime) * time.Minute).Unix(),
		},
		Jid: jid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(config.JwtRefreshTokenSecretKey))
	if err != nil {
		return "", 0, err
	}

	return signed, jid, nil
}

func CreateJwtToken(config *shared.EnvConf, user models.User) (string, string, int64, error) {
	refreshToken, jid, err := createRefreshToken(config)
	if err != nil {
		return "", "", 0, err
	}
	accessToken, err := createAccessToken(config, user, jid)
	if err != nil {
		return "", "", 0, err
	}

	return refreshToken, accessToken, jid, nil
}

func VerifyAccessToken(config *shared.EnvConf, token string) (bool, dtos.AccessClaims) {
	claims := dtos.AccessClaims{}
	tkn, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.JwtAccessTokenSecretKey), nil
	})

	if err != nil {
		return false, dtos.AccessClaims{}
	}

	if !tkn.Valid {
		return false, dtos.AccessClaims{}
	}

	return true, claims
}

func VerifyRefreshToken(config *shared.EnvConf, token string) (bool, int64) {
	claims := dtos.RefreshClaims{}
	tkn, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.JwtRefreshTokenSecretKey), nil
	})

	if err != nil {
		return false, 0
	}

	if !tkn.Valid {
		return false, 0
	}

	return true, claims.Jid
}
