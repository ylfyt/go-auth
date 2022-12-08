package main

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ListenPort string
}

func (me *Config) loadConfig() error {
	listenPort := os.Getenv("LISTEN_PORT")
	if listenPort == "" {
		return errors.New("listen port is not found")
	}

	me.ListenPort = listenPort
	return nil
}

func (me *Config) Init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	err = me.loadConfig()
	if err != nil {
		panic(err)
	}
}
