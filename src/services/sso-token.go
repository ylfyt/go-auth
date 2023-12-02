package services

import (
	"fmt"
	"go-auth/src/dtos"
)

type IStorage struct {
}

type SsoTokenService struct {
	storage map[int64]dtos.TokenPayload
}

func NewSsoTokenService() *SsoTokenService {
	return &SsoTokenService{
		storage: make(map[int64]dtos.TokenPayload),
	}
}

func (me *SsoTokenService) Exchange(exchangeToken int64) (*dtos.TokenPayload, error) {
	token, exist := me.storage[exchangeToken]
	if !exist {
		return nil, fmt.Errorf("exchange %d not found", exchangeToken)
	}
	return &token, nil
}

func (me *SsoTokenService) Store(exchangeToken int64, token dtos.TokenPayload) error {
	me.storage[exchangeToken] = token
	return nil
}
