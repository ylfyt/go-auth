package controllers

import (
	"go-auth/src/dtos"
	"go-auth/src/middlewares"
	"go-auth/src/services"
	"go-auth/src/shared"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
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

type Controller struct {
	db              *sqlx.DB
	config          *shared.EnvConf
	ssoTokenService *services.SsoTokenService
}

func New(_db *sqlx.DB, _config *shared.EnvConf, _ssoService *services.SsoTokenService) *chi.Mux {
	controller := Controller{
		db:              _db,
		config:          _config,
		ssoTokenService: _ssoService,
	}
	authRoute := Route{
		Base: "/auth",
		EndPoints: []EndPoint{
			{
				Method:  "POST",
				Path:    "/login",
				Handler: controller.login,
				Middlewares: []func(next http.Handler) http.Handler{
					middlewares.BodyParser[dtos.Register](),
				},
			},
			{
				Method:  "POST",
				Path:    "/register",
				Handler: controller.register,
				Middlewares: []func(next http.Handler) http.Handler{
					middlewares.BodyParser[dtos.Register](),
				},
			},
			{
				Method:  "POST",
				Path:    "/refresh",
				Handler: controller.refreshToken,
				Middlewares: []func(next http.Handler) http.Handler{
					middlewares.BodyParser[dtos.RefreshPayload](),
				},
			},
			{
				Method:  "POST",
				Path:    "/logout",
				Handler: controller.logout,
				Middlewares: []func(next http.Handler) http.Handler{
					middlewares.BodyParser[dtos.RefreshPayload](),
				},
			},
			{
				Method:  "POST",
				Path:    "/logout-all",
				Handler: controller.logoutAll,
				Middlewares: []func(next http.Handler) http.Handler{
					middlewares.BodyParser[dtos.Register](),
				},
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
				Middlewares: []func(next http.Handler) http.Handler{
					middlewares.Authorization(_config),
				},
			},
			{
				Method:  "GET",
				Path:    "/{id}",
				Handler: controller.getUserById,
				Middlewares: []func(next http.Handler) http.Handler{
					middlewares.Authorization(_config),
				},
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
	ssoRoute := Route{
		Base: "/sso",
		EndPoints: []EndPoint{
			{
				Method:  "POST",
				Path:    "/login",
				Handler: controller.ssoLogin,
				Middlewares: []func(next http.Handler) http.Handler{
					middlewares.BodyParser[dtos.SsoLoginPayload](),
				},
			},
			{
				Method:  "GET",
				Path:    "/client/{id}",
				Handler: controller.getSsoClient,
			},
		},
	}

	routes := []Route{
		authRoute,
		userRoute,
		homeRoute,
		ssoRoute,
	}

	r := chi.NewRouter()
	r.Use(middlewares.AccessLogger)
	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	}).Handler)

	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/login.html")
	})

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
