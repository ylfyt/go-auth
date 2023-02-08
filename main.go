package main

import (
	"fmt"
	"go-auth/src/meta"
)

func main() {
	app := meta.New(&meta.Config{
		BaseUrl: "/api",
	})

	app.Map("GET", "/", func(c *meta.Ctx) meta.ResponseDto {
		fmt.Printf("Data: %+v\n", c)
		return meta.ResponseDto{
			Status:  200,
			Success: true,
			Data:    "ok",
		}
	})

	app.Run(3000)
}
