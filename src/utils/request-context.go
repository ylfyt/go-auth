package utils

import (
	"context"
	"net/http"
)

type CtxKeyType int

const (
	CTX_CLAIM_KEY CtxKeyType = iota + 0
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
