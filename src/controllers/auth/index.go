package auth

import (
	"go-auth/src/middlewares"
	"go-auth/src/structs"

	"github.com/ylfyt/go_db/go_db"
	"github.com/ylfyt/meta/meta"
)

type AuthController struct {
	db     *go_db.DB
	config *structs.EnvConf
}

func (me *AuthController) Setup(db *go_db.DB, config *structs.EnvConf) []meta.EndPoint {
	me.db = db
	me.config = config

	return []meta.EndPoint{
		{
			Method:      "POST",
			Path:        "/auth/register",
			HandlerFunc: me.register,
		},
		{
			Method:      "POST",
			Path:        "/auth/login",
			HandlerFunc: me.login,
		},
		{
			Method:      "POST",
			Path:        "/auth/refresh-token",
			HandlerFunc: me.refreshToken,
		},
		{
			Method:      "POST",
			Path:        "/auth/logout",
			HandlerFunc: me.logout,
		},
		{
			Method:      "POST",
			Path:        "/auth/logout-all",
			HandlerFunc: me.logoutAll,
		},
		{
			Method:      "GET",
			Path:        "/auth/users",
			HandlerFunc: me.getUsers,
			Middlewares: []any{
				middlewares.Authorization,
			},
		},
		{
			Method:      "GET",
			Path:        "/auth/users/:id",
			HandlerFunc: me.getUserById,
			Middlewares: []any{
				// middlewares.Authorization,
			},
		},
	}
}
