package controllers

import (
	"fmt"
	"go-auth/src/middlewares"
	"go-auth/src/structs"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ylfyt/go_db/go_db"
)

type EndPoint struct {
	Method      string
	Path        string
	Handler     func(http.ResponseWriter, *http.Request)
	Middlewares []func(next http.Handler) http.Handler
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
				Middlewares: []func(next http.Handler) http.Handler{
					middlewares.Authorization(_config),
				},
			},
			{
				Method:  "GET",
				Path:    "/ping",
				Handler: controller.ping,
				Middlewares: []func(next http.Handler) http.Handler{
					func(next http.Handler) http.Handler {
						return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							fmt.Println("PING")
							next.ServeHTTP(w, r)
						})
					},
				},
			},
		},
	}
	routes := []Route{
		authRoute,
		userRoute,
		homeRoute,
	}
	r := chi.NewRouter()
	r.Use(middlewares.AccessLogger)
	r.Route("/api", func(r chi.Router) {
		for _, route := range routes {
			r.Route(route.Base, func(r chi.Router) {
				for _, endpoint := range route.EndPoints {
					sub := r.Group(nil)
					for _, mid := range endpoint.Middlewares {
						sub.Use(mid)
					}
					sub.MethodFunc(endpoint.Method, endpoint.Path, endpoint.Handler)
				}
			})
		}
	})

	return r
}
