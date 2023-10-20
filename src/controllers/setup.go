package controllers

import "github.com/ylfyt/meta/meta"

func New() *meta.App {
	app := meta.New(&meta.Config{
		BaseUrl: "/api",
	})

	for _, v := range appRoutes {
		app.AddEndPoint(v...)
	}

	return app
}
