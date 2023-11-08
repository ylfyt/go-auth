package main

import (
	"fmt"
	"go-auth/src/controllers"
	"go-auth/src/structs"
	"go-auth/src/utils"
	"net/http"

	"github.com/caarlos0/env/v9"
	"github.com/ylfyt/go_db/go_db"
)

func main() {
	utils.LoadEnv()

	var config structs.EnvConf
	if err := env.ParseWithOptions(&config, env.Options{
		RequiredIfNoDef:       true,
		UseFieldNameByDefault: true,
	}); err != nil {
		panic(err)
	}
	fmt.Printf("Data: %+v\n", config)

	db, err := go_db.New(config.DbConnection, go_db.Option{
		MaxOpenConn:     50,
		MaxIdleConn:     10,
		MaxIdleLifeTime: 300,
	})
	if err != nil {
		panic(err)
	}

	app := controllers.New(db, &config)

	fmt.Println("Listening on port", config.ListenPort)
	http.ListenAndServe(fmt.Sprintf(":%d", config.ListenPort), app)
}
