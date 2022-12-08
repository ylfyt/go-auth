package main

import (
	"fmt"
	"net/http"

	"go-auth/src/controllers"
)

var ServerConfig Config

func init() {
	ServerConfig.Init()
}

func main() {
	router := controllers.NewRouter()

	fmt.Println("Server is listening on port", ServerConfig.ListenPort)	
	http.ListenAndServe("127.0.0.1:"+ServerConfig.ListenPort, router)
}
