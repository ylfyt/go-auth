package ctx

import (
	"go-auth/src/models"
	"net/http"
)

func Get[T any](r *http.Request, k key) T {
	val := r.Context().Value(k)
	if val == nil {
		var temp T
		return temp
	}
	return val.(T)
}

func GetUserClaimsCtx(r *http.Request) models.AccessClaims {
	return Get[models.AccessClaims](r, UserClaimsCtxKey)
}

func GetReqIdCtx(r *http.Request) string {
	return Get[string](r, ReqIdCtxKey)
}
