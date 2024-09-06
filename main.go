package main

import (
	"errors"
	"fmt"
	"go-auth/src/controllers"
	"go-auth/src/logger"
	"go-auth/src/services"
	"go-auth/src/shared"
	"go-auth/src/utils"
	"net/http"
	"os"

	"github.com/caarlos0/env/v9"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/json-iterator/go/extra"
	_ "modernc.org/sqlite"
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
	// f, err := os.OpenFile("./test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	panic(err)
	// }
	l := logger.NewLogger(os.Stdout, logger.LOG_INFO)
	l.If("Config: %+v", config)

	db, err := sqlx.Open("sqlite", config.DbConnection)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if driver, err := sqlite.WithInstance(db.DB, &sqlite.Config{}); err != nil {
		panic(err)
	} else if m, err := migrate.NewWithDatabaseInstance("file://./migrations", "sqlite", driver); err != nil {
		panic(err)
	} else if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			panic(fmt.Errorf("cannot perform database migration: %+v", err))
		}
	}

	ssoTokenService := services.NewSsoTokenService()
	app := controllers.New(db, &config, ssoTokenService)

	fmt.Println("Listening on port", config.ListenPort)
	http.ListenAndServe(fmt.Sprintf(":%d", config.ListenPort), app)
}
