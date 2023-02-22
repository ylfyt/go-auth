package main

import (
	"database/sql"
	"go-auth/src/config"
	"go-auth/src/controllers"
	"go-auth/src/meta"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", config.DB_CONNECTION)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	
	app := controllers.New()
	meta.AddService(app, db)

	app.Run(3000)
}
