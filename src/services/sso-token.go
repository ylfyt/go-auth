package services

import (
	"go-auth/src/dtos"

	"github.com/jmoiron/sqlx"
)

type IStorage struct {
}

type SsoTokenService struct {
	db *sqlx.DB
}

func NewSsoTokenService(_db *sqlx.DB) *SsoTokenService {
	return &SsoTokenService{
		db: _db,
	}
}

func (me *SsoTokenService) Exchange(exchangeToken string) (*dtos.TokenPayload, error) {
	var token *dtos.TokenPayload
	err := me.db.Get(&token, `SELECT * FROM tokens WHERE exchange = ?`, exchangeToken)
	return token, err
}

func (me *SsoTokenService) Store(exchangeToken int64, token dtos.TokenPayload) error {
	_, err := me.db.Exec(`
		INSERT INTO tokens (exchange, access, refresh, expired)
		VALUES (?, ?, ?, ?)
	`, exchangeToken, token.AccessToken, token.RefreshToken, token.ExpiredIn)
	if err != nil {
		return err
	}

	return nil
}
