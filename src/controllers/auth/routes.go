package auth

import "go-auth/src/ctx"

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
		Pattern:     "/auth/refresh",
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
}
