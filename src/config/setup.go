package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

var (
	LISTEN_PORT   string
	DB_CONNECTION string
)

func loadConfig() error {
	listenPort := os.Getenv("LISTEN_PORT")
	if listenPort == "" {
		return errors.New("listen port is not found")
	}
	dbConn := os.Getenv("DB_CONNECTION")
	if dbConn == "" {
		return errors.New("DB Connection is not found")
	}

	LISTEN_PORT = listenPort
	DB_CONNECTION = dbConn
	return nil
}

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	err = loadConfig()
	if err != nil {
		panic(err)
	}
}
