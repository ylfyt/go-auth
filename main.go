package main

import (
	"fmt"
	"go-auth/src/controllers"
	"go-auth/src/structs"
	"go-auth/src/utils"
	"net/http"

	"github.com/caarlos0/env/v9"
	"github.com/json-iterator/go/extra"
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
}

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
