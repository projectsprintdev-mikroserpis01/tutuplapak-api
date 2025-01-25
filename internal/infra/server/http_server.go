package server

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	authController "github.com/projectsprintdev-mikroserpis01/tutuplapak-api/internal/app/auth/controller"
	authRepo "github.com/projectsprintdev-mikroserpis01/tutuplapak-api/internal/app/auth/repository"
	authSvc "github.com/projectsprintdev-mikroserpis01/tutuplapak-api/internal/app/auth/service"
	userController "github.com/projectsprintdev-mikroserpis01/tutuplapak-api/internal/app/user/controller"
	userRepo "github.com/projectsprintdev-mikroserpis01/tutuplapak-api/internal/app/user/repository"
	userSvc "github.com/projectsprintdev-mikroserpis01/tutuplapak-api/internal/app/user/service"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/internal/middlewares"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/pkg/bcrypt"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/pkg/helpers/http/response"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/pkg/jwt"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/pkg/validator"

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
	validator := validator.Validator
	jwt := jwt.Jwt
	bcrypt := bcrypt.Bcrypt
	middleware := middlewares.NewMiddleware(jwt)

	s.app.Get("/", func(c *fiber.Ctx) error {
		return response.SendResponse(c, fiber.StatusOK, "Welcome to Tutuplapak API")
	})

	v1 := s.app.Group("/v1")

	authRepository := authRepo.NewAuthRepository(db)
	userRepository := userRepo.NewUserRepository(db)

	authService := authSvc.NewAuthService(authRepository, validator, bcrypt, jwt)
	userService := userSvc.NewUserService(userRepository, validator)

	authController.InitAuthController(v1, authService)
	userController.InitUserController(v1, userService, middleware)

	s.app.Use(func(c *fiber.Ctx) error {
		return c.SendFile("./web/not-found.html")
	})
}

func (s httpServer) GetApp() *fiber.App {
	return s.app
}
