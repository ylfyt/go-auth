package main

import "go-auth/src/controllers"

func main() {
	app := controllers.New()

	app.Run(3000)
}
