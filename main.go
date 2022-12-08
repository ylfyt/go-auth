package main

import (
	"fmt"
	"net/http"

	"go-auth/src/config"
	"go-auth/src/controllers"
)

func main() {
	router := controllers.NewRouter()

	fmt.Println("Server is listening on port", config.LISTEN_PORT)
	http.ListenAndServe("127.0.0.1:"+config.LISTEN_PORT, router)
}
