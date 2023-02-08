package main

import "go-auth/src/meta"

func main() {
	app := meta.New(&meta.Config{
		BaseUrl: "/api",
	})

	app.Map("GET", "/", func() meta.ResponseDto {
		return meta.ResponseDto{
			Status:  200,
			Success: true,
			Data:    "ok",
		}
	})

	app.Run(3000)
}
