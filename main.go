package main

import (
	"fmt"
	"go-auth/src/config"
	"go-auth/src/controllers"
)

func main() {
	app := controllers.New()

	fmt.Println("Server is listening on port", config.LISTEN_PORT)
	app.Listen(":" + config.LISTEN_PORT)
}
