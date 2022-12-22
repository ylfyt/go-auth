package main

import (
	"fmt"
	"go-auth/src/config"
	"go-auth/src/controllers"
	"net/http"
)

func main() {
	router := controllers.NewRouter()

	fmt.Println("Server is listening on port", config.LISTEN_PORT)
	http.ListenAndServe("127.0.0.1:"+config.LISTEN_PORT, router)
}
