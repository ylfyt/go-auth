package services

import (
	"go-auth/src/config"
	"go-auth/src/models"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func createAccessToken(user models.User, refreshTokenId uuid.UUID) (string, error) {
	claims := models.AccessClaims{
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

func createRefreshToken() (string, *uuid.UUID, error) {
	jid := uuid.New()
	claims := models.RefreshClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(config.JWT_REFRESH_TOKEN_EXPIRY_TIME) * time.Minute).Unix(),
		},
		Jid: jid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(config.JWT_REFRESH_TOKEN_SECRET_KEY))
	if err != nil {
		return "", nil, err
	}

	return signed, &jid, nil
}

func CreateJwtToken(user models.User) (string, string, *uuid.UUID, error) {
	refreshToken, jid, err := createRefreshToken()
	if err != nil {
		return "", "", nil, err
	}
	accessToken, err := createAccessToken(user, *jid)
	if err != nil {
		return "", "", nil, err
	}

	return refreshToken, accessToken, jid, nil
}

func VerifyAccessToken(token string) *models.AccessClaims {
	claims := models.AccessClaims{}
	tkn, _ := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.JWT_ACCESS_TOKEN_SECRET_KEY), nil
	})

	if !tkn.Valid {
		return nil
	}

	return &claims
}

func VerifyRefreshToken(token string) (bool, *uuid.UUID) {
	claims := models.RefreshClaims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.JWT_REFRESH_TOKEN_SECRET_KEY), nil
	})

	if err == nil {
		return true, &claims.Jid
	}

	if err.Error() == jwt.ErrSignatureInvalid.Error() {
		return false, nil
	}

	if (claims == models.RefreshClaims{}) {
		return false, nil
	}

	return false, &claims.Jid
}
