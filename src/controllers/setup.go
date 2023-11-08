package controllers

import (
	"go-auth/src/structs"
	"net/http"

	"github.com/go-chi/chi"
	jsoniter "github.com/json-iterator/go"
	"github.com/ylfyt/go_db/go_db"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type EndPoint struct {
	Method  string
	Path    string
	Handler func(http.ResponseWriter, *http.Request)
}

type Route struct {
	Base      string
	EndPoints []EndPoint
}

type ChiController struct {
	db     *go_db.DB
	config *structs.EnvConf
}

func New(_db *go_db.DB, _config *structs.EnvConf) *chi.Mux {
	controller := ChiController{
		db:     _db,
		config: _config,
	}
	authRoute := Route{
		Base: "/auth",
		EndPoints: []EndPoint{
			{
				Method:  "POST",
				Path:    "/login",
				Handler: controller.login,
			},
			{
				Method:  "POST",
				Path:    "/register",
				Handler: controller.register,
			},
			{
				Method:  "POST",
				Path:    "/refresh",
				Handler: controller.refreshToken,
			},
			{
				Method:  "POST",
				Path:    "/logout",
				Handler: controller.logout,
			},
			{
				Method:  "POST",
				Path:    "/logout-all",
				Handler: controller.logoutAll,
			},
		},
	}
	userRoute := Route{
		Base: "/user",
		EndPoints: []EndPoint{
			{
				Method:  "GET",
				Path:    "/",
				Handler: controller.getUsers,
			},
			{
				Method:  "GET",
				Path:    "/{id}",
				Handler: controller.getUserById,
			},
		},
	}
	homeRoute := Route{
		Base: "/",
		EndPoints: []EndPoint{
			{
				Method:  "GET",
				Path:    "/",
				Handler: controller.home,
			},
			{
				Method:  "GET",
				Path:    "/ping",
				Handler: controller.ping,
			},
		},
	}
	routes := []Route{
		authRoute,
		userRoute,
		homeRoute,
	}

	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		for _, route := range routes {
			r.Route(route.Base, func(r chi.Router) {
				for _, endpoint := range route.EndPoints {
					r.MethodFunc(endpoint.Method, endpoint.Path, endpoint.Handler)
				}
			})
		}
	})

	return r
}
