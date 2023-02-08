package main

import (
	"fmt"
	"go-auth/src/meta"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := meta.New(&meta.Config{
		BaseUrl: "/api",
	})

	endPoint := []meta.EndPoint{
		{
			Method: "GET",
			Path:   "/ping",
			HandlerFunc: func(c *fiber.Ctx) meta.ResponseDto {
				return meta.ResponseDto{
					Status:  200,
					Success: true,
					Data:    "pong",
				}
			},
		},
		{
			Method: "GET",
			Path:   "/hello",
			HandlerFunc: func(c *fiber.Ctx) meta.ResponseDto {
				return meta.ResponseDto{
					Status:  400,
					Success: true,
					Data:    "world",
				}
			},
		},
	}

	app.Map("GET", "/", func(c *fiber.Ctx) meta.ResponseDto {
		fmt.Printf("Data: %+v\n", c)
		return meta.ResponseDto{
			Status:  200,
			Success: true,
			Data:    "ok",
		}
	})

	app.AddEndPoint(endPoint...)

	app.Run(3000)
}
