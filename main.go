package main

import (
	"fmt"
	"go-auth/src/controllers"
	"go-auth/src/services"
	"go-auth/src/shared"
	"go-auth/src/utils"
	"net/http"

	"github.com/caarlos0/env/v9"
	"github.com/jmoiron/sqlx"
	"github.com/json-iterator/go/extra"
	_ "github.com/mattn/go-sqlite3"
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
	var config shared.EnvConf
	if err := env.ParseWithOptions(&config, env.Options{
		RequiredIfNoDef:       true,
		UseFieldNameByDefault: true,
	}); err != nil {
		panic(err)
	}
	fmt.Printf("Data: %+v\n", config)

	db, err := sqlx.Open("sqlite3", config.DbConnection)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ssoTokenService := services.NewSsoTokenService()
	app := controllers.New(db, &config, ssoTokenService)

	fmt.Println("Listening on port", config.ListenPort)
	http.ListenAndServe(fmt.Sprintf(":%d", config.ListenPort), app)
}
