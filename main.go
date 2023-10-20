package main

import (
	"go-auth/src/config"
	"go-auth/src/controllers"
	"go-auth/src/middlewares"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/ylfyt/go_db/go_db"
)

func main() {
	config.Print()

	db, err := go_db.New(config.DB_CONNECTION, go_db.Option{
		MaxOpenConn:     50,
		MaxIdleConn:     10,
		MaxIdleLifeTime: 300,
	})
	if err != nil {
		panic(err)
	}

	app := controllers.New()
	app.Use(middlewares.AccessLogger)
	app.AddService(db)

	port, err := strconv.Atoi(config.LISTEN_PORT)
	if err != nil {
		port = 3000
	}
	app.Run(port)
}
