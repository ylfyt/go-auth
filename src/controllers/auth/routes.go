package auth

import (
	"go-auth/src/middlewares"

	"github.com/ylfyt/meta/meta"
)

var Routes = []meta.EndPoint{
	{
		Method:      "POST",
		Path:        "/auth/register",
		HandlerFunc: register,
	},
	{
		Method:      "POST",
		Path:        "/auth/login",
		HandlerFunc: login,
	},
	{
		Method:      "POST",
		Path:        "/auth/refresh-token",
		HandlerFunc: refreshToken,
	},
	{
		Method:      "POST",
		Path:        "/auth/logout",
		HandlerFunc: logout,
	},
	{
		Method:      "POST",
		Path:        "/auth/logout-all",
		HandlerFunc: logoutAll,
	},
	{
		Method:      "GET",
		Path:        "/auth/users",
		HandlerFunc: getUsers,
		Middlewares: []any{
			middlewares.Authorization,
		},
	},
	{
		Method:      "GET",
		Path:        "/auth/users/:id",
		HandlerFunc: getUserById,
		Middlewares: []any{
			// middlewares.Authorization,
		},
	},
}
