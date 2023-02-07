package auth

import (
	"go-auth/src/ctx"
	"go-auth/src/middlewares"
)

var Routes = []ctx.Route{
	{
		Name:        "AuthRegister",
		Method:      "POST",
		Pattern:     "/auth/register",
		HandlerFunc: register,
	},
	{
		Name:        "AuthLogin",
		Method:      "POST",
		Pattern:     "/auth/login",
		HandlerFunc: login,
	},
	{
		Name:        "AuthRefresh",
		Method:      "POST",
		Pattern:     "/auth/refresh-token",
		HandlerFunc: refreshToken,
	},
	{
		Name:        "AuthLogout",
		Method:      "POST",
		Pattern:     "/auth/logout",
		HandlerFunc: logout,
	},
	{
		Name:        "AuthLogoutAll",
		Method:      "POST",
		Pattern:     "/auth/logout-all",
		HandlerFunc: logoutAll,
	},
	{
		Name:        "AuthUsers",
		Method:      "GET",
		Pattern:     "/auth/users",
		HandlerFunc: getUsers,
		Middlewares: []interface{}{
			middlewares.Authorization,
		},
	},
	{
		Name:        "AuthGetUsersById",
		Method:      "GET",
		Pattern:     "/auth/users/{id}",
		HandlerFunc: getUserById,
		Middlewares: []interface{}{
			// middlewares.Authorization,
		},
	},
}
