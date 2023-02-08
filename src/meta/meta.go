package meta

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func New(config *Config) *App {
	fiberApp := fiber.New(fiber.Config{
		StrictRouting: true,
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
	})

	router := fiberApp.Group(config.BaseUrl)

	return &App{
		fiberApp: fiberApp,
		config:   config,
		router:   router,
	}
}

func (app *App) Map(method string, path string, handler interface{}) {
	err := app.validateHandler(handler)
	if err != nil {
		panic(err)
	}
	app.endPoints = append(app.endPoints, EndPoint{
		Method:      method,
		Path:        path,
		HandlerFunc: handler,
	})
}

func (app *App) Run(port int) {
	app.setup()
	fmt.Println("App running on port", port)
	app.fiberApp.Listen(fmt.Sprintf(":%d", port))
}
