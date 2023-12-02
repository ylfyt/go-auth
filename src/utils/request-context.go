package utils

import (
	"context"
	"go-auth/src/dtos"
	"net/http"
)

type CtxKeyType int

const (
	CTX_AUTH_CLAIM_KEY CtxKeyType = iota + 0
	CTX_BODY_KEY
	CTX_REQ_ID_KEY
)

func GetContext[T any](r *http.Request, key CtxKeyType) *T {
	val := r.Context().Value(key)
	if val == nil {
		return nil
	}

	if z, ok := val.(T); ok {
		return &z
	}

	return nil
}

func SetContext(r *http.Request, key CtxKeyType, val any) *http.Request {
	ctx := context.WithValue(r.Context(), key, val)
	return r.WithContext(ctx)
}

func GetBodyContext[T any](r *http.Request) *T {
	return GetContext[T](r, CTX_BODY_KEY)
}

func GetAuthClaimContext(r *http.Request) *dtos.AccessClaims {
	return GetContext[dtos.AccessClaims](r, CTX_AUTH_CLAIM_KEY)
}
