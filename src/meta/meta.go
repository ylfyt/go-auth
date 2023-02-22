package meta

import (
	"encoding/json"
	"fmt"
	"reflect"

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
		fiberApp:     fiberApp,
		config:       config,
		router:       router,
		dependencies: make(map[string]any),
	}
}

func (app *App) Map(method string, path string, handler interface{}, middlewares ...func(c *fiber.Ctx) error) {
	app.endPoints = append(app.endPoints, EndPoint{
		Method:      method,
		Path:        path,
		HandlerFunc: handler,
		Middlewares: middlewares,
	})
}

func (app *App) AddEndPoint(endPoints ...EndPoint) {
	app.endPoints = append(app.endPoints, endPoints...)
}

func (app *App) Run(port int) {
	for _, v := range app.endPoints {
		err := app.validateHandler(v.HandlerFunc)
		if err != nil {
			fmt.Print(v.Path, " ")
			panic(err)
		}
	}
	app.setup()
	fmt.Println("App running on port", port)
	app.fiberApp.Listen(fmt.Sprintf("0.0.0.0:%d", port))
}

func AddService[T any](app *App, service *T) {
	app.dependencies[reflect.TypeOf(service).String()] = service
}

func (app *App) Use(path string, handler func(c *fiber.Ctx) error) {
	if path == "" {
		app.fiberApp.Use(handler)
		return
	}
	app.fiberApp.Use(path, handler)
}
