package server

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/internal/middlewares"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/pkg/helpers/http/response"

	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/pkg/log"
)

type HttpServer interface {
	Start(port string)
	MountMiddlewares()
	MountRoutes(db *sqlx.DB)
	GetApp() *fiber.App
}

type httpServer struct {
	app *fiber.App
}

func NewHttpServer() HttpServer {
	config := fiber.Config{
		CaseSensitive: true,
		AppName:       "Tutuplapak-API",
		ServerHeader:  "Tutuplapak",
		JSONEncoder:   sonic.Marshal,
		JSONDecoder:   sonic.Unmarshal,
		//ErrorHandler:  errorhandler.ErrorHandler,
	}
	app := fiber.New(config)
	return &httpServer{
		app: app,
	}
}

func (s httpServer) Start(port string) {
	if port[0] != ':' {
		port = ":" + port
	}

	err := s.app.Listen(port)

	if err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[SERVER][Start] failed to start server")
	}
}

func (s httpServer) MountMiddlewares() {
	s.app.Use(middlewares.LoggerConfig())
	s.app.Use(middlewares.Helmet())
	s.app.Use(middlewares.Cors())
	s.app.Use(middlewares.RecoverConfig())
}

func (s httpServer) MountRoutes(db *sqlx.DB) {
	s.app.Get("/", func(c *fiber.Ctx) error {
		return response.SendResponse(c, fiber.StatusOK, "Welcome to Tutuplapak API")
	})

	api := s.app.Group("/v1")

	api.Get("/", func(c *fiber.Ctx) error {
		return response.SendResponse(c, fiber.StatusOK, "TutupLapak API v1")
	})

	s.app.Use(func(c *fiber.Ctx) error {
		return c.SendFile("./web/not-found.html")
	})
}

func (s httpServer) GetApp() *fiber.App {
	return s.app
}
