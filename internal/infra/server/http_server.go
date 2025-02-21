package server

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	authController "github.com/projectsprintdev-mikroserpis01/tutuplapak-api/internal/app/auth/controller"
	authRepo "github.com/projectsprintdev-mikroserpis01/tutuplapak-api/internal/app/auth/repository"
	authSvc "github.com/projectsprintdev-mikroserpis01/tutuplapak-api/internal/app/auth/service"
	productSvc "github.com/projectsprintdev-mikroserpis01/tutuplapak-api/internal/app/product/service"
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

	api := s.app.Group("/v1")

	authRepository := authRepo.NewAuthRepository(db)
	userRepository := userRepo.NewUserRepository(db)

	authService := authSvc.NewAuthService(authRepository, validator, bcrypt, jwt)
	userService := userSvc.NewUserService(userRepository, validator)

	authController.InitAuthController(api, authService)
	userController.InitUserController(api, userService, middleware)

	api.Get("/", func(c *fiber.Ctx) error {
		return response.SendResponse(c, fiber.StatusOK, "TutupLapak API v1")
	})

	// Product Routes
	products := api.Group("/products")

	// Create Product
	products.Post("/", func(c *fiber.Ctx) error {
		var data json.RawMessage
		if err := c.BodyParser(&data); err != nil {
			return response.SendResponse(c, fiber.StatusBadRequest, "Invalid request")
		}
		id := GenerateUUID() // Generate a random product ID
		productSvc.CreateProduct(c.Context(), id, data)
		return response.SendResponse(c, fiber.StatusCreated, fiber.Map{"id": id, "message": "Product created"})
	})

	// Get Products (Mock Response)
	products.Get("/", func(c *fiber.Ctx) error {
		// Simulate returning some products (Replace with real DB query)
		mockProducts := []fiber.Map{
			{"id": GenerateUUID(), "name": "Product 1", "price": 10000},
			{"id": GenerateUUID(), "name": "Product 2", "price": 20000},
		}
		return response.SendResponse(c, fiber.StatusOK, mockProducts)
	})

	// Update Product
	products.Put("/:id", func(c *fiber.Ctx) error {
		var data json.RawMessage
		id := c.Params("id")
		if err := c.BodyParser(&data); err != nil {
			return response.SendResponse(c, fiber.StatusBadRequest, "Invalid request")
		}
		productSvc.UpdateProduct(c.Context(), id, data)
		return response.SendResponse(c, fiber.StatusOK, fiber.Map{"id": id, "message": "Product updated"})
	})

	// Delete Product
	products.Delete("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		productSvc.DeleteProduct(c.Context(), id, nil)
		return response.SendResponse(c, fiber.StatusOK, fiber.Map{"id": id, "message": "Product deleted"})
	})

	s.app.Use(func(c *fiber.Ctx) error {
		return c.SendFile("./web/not-found.html")
	})
}

func (s httpServer) GetApp() *fiber.App {
	return s.app
}

// GenerateUUID creates a random 16-byte hex string (UUID-like)
func GenerateUUID() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "random-id" // Fallback if random fails
	}
	return hex.EncodeToString(bytes)
}
