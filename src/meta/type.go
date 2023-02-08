package meta

import "github.com/gofiber/fiber/v2"

type EndPoint struct {
	Method      string
	Path        string
	HandlerFunc interface{}
}

type Config struct {
	BaseUrl string
}

type App struct {
	fiberApp  *fiber.App
	router    fiber.Router
	config    *Config
	endPoints []EndPoint
}

type Error struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

type ResponseDto struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Errors  []Error     `json:"errors"`
	Data    interface{} `json:"data"`
}
