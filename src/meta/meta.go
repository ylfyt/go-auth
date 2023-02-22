package meta

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func New(config *Config) *App {
	fiberApp := fiber.New(fiber.Config{
		StrictRouting: false,
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

func (app *App) Map(method string, path string, handler interface{}, middlewares ...func(c *fiber.Ctx) error) {
	err := app.validateHandler(handler)
	if err != nil {
		panic(err)
	}
	app.endPoints = append(app.endPoints, EndPoint{
		Method:      method,
		Path:        path,
		HandlerFunc: handler,
		Middlewares: middlewares,
	})
}

func (app *App) AddEndPoint(endPoints ...EndPoint) {
	for _, v := range endPoints {
		err := app.validateHandler(v.HandlerFunc)
		if err != nil {
			fmt.Print(v.Path, " ")
			panic(err)
		}
	}
	app.endPoints = append(app.endPoints, endPoints...)
}

func (app *App) Run(port int) {
	app.setup()
	fmt.Println("App running on port", port)
	app.fiberApp.Listen(fmt.Sprintf("0.0.0.0:%d", port))
}
