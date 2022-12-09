package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	LISTEN_PORT                  string
	DB_CONNECTION                string
	JWT_ACCESS_TOKEN_EXPIRY_TIME int
	JWT_ACCESS_TOKEN_SECRET_KEY  string
)

func loadConfig() error {
	temp := ""
	temp = os.Getenv("LISTEN_PORT")
	if temp == "" {
		return errors.New("listen port is not found")
	}
	LISTEN_PORT = temp

	temp = os.Getenv("DB_CONNECTION")
	if temp == "" {
		return errors.New("DB Connection is not found")
	}
	DB_CONNECTION = temp

	temp = os.Getenv("JWT_ACCESS_TOKEN_SECRET_KEY")
	if temp == "" {
		return errors.New("JWT_ACCESS_TOKEN_SECRET_KEY is not found")
	}
	JWT_ACCESS_TOKEN_SECRET_KEY = temp

	temp = os.Getenv("JWT_ACCESS_TOKEN_EXPIRY_TIME")
	if temp == "" {
		return errors.New("JWT_ACCESS_TOKEN_EXPIRY_TIME is not found")
	}
	expiry, err := strconv.Atoi(temp)
	if err != nil {
		return err
	}
	JWT_ACCESS_TOKEN_EXPIRY_TIME = expiry

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
