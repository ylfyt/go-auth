package main

import (
	"fmt"
	"go-auth/src/controllers"
	"go-auth/src/services"
	"go-auth/src/structs"
	"go-auth/src/utils"
	"net/http"

	"github.com/caarlos0/env/v9"
	"github.com/jmoiron/sqlx"
	"github.com/json-iterator/go/extra"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ylfyt/go_db/go_db"
)

func init() {
	extra.SetNamingStrategy(func(s string) string {
		if len(s) < 1 {
			return s
		}
		first := s[0]
		if 'A' <= first && first <= 'Z' {
			first += 'a' - 'A'
		}
		strBytes := []byte(s)
		strBytes[0] = first
		return string(strBytes)
	})
	
	utils.LoadEnv()
}

func main() {
	var config structs.EnvConf
	if err := env.ParseWithOptions(&config, env.Options{
		RequiredIfNoDef:       true,
		UseFieldNameByDefault: true,
	}); err != nil {
		panic(err)
	}
	fmt.Printf("Data: %+v\n", config)

	db2, err := sqlx.Open("sqlite3", "example.db")
	if err != nil {
		panic(err)
	}
	defer db2.Close()
	ssoTokenService := services.NewSsoTokenService(db2)

	db, err := go_db.New(config.DbConnection, go_db.Option{
		MaxOpenConn:     50,
		MaxIdleConn:     10,
		MaxIdleLifeTime: 300,
	})
	if err != nil {
		panic(err)
	}

	app := controllers.New(db, &config, ssoTokenService)

	fmt.Println("Listening on port", config.ListenPort)
	http.ListenAndServe(fmt.Sprintf(":%d", config.ListenPort), app)
}
