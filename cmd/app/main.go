package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/internal/infra/database"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/internal/infra/env"
	"github.com/projectsprintdev-mikroserpis01/tutuplapak-api/internal/infra/server"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// Define Prometheus metrics
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "route", "instance"},
	)
	httpDurationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_duration_seconds",
			Help:    "Histogram of HTTP request durations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route", "instance"},
	)
)

func init() {
	// Register metrics with Prometheus
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpDurationHistogram)
}

// Prometheus Middleware to track HTTP request metrics
func PrometheusMiddleware(c *fiber.Ctx) error {
	route := c.Path()
	start := time.Now()

	// Get the app instance dynamically from the environment variable
	instance := os.Getenv("HOSTNAME") // This gives you the unique hostname for each container (e.g., app_1, app_2, etc.)

	err := c.Next() // Call the next handler in the chain

	// Track request metrics
	httpRequestsTotal.WithLabelValues(c.Method(), route, instance).Inc()
	duration := time.Since(start).Seconds()
	httpDurationHistogram.WithLabelValues(c.Method(), route, instance).Observe(duration)

	return err
}

// Helper function to convert http.Handler to fiber.Handler
func adaptor(h http.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		handler := fasthttpadaptor.NewFastHTTPHandler(h)
		handler(c.Context())
		return nil
	}
}

func main() {
	// Initialize server and database
	server := server.NewHttpServer()
	psqlDB := database.NewPgsqlConn()
	defer psqlDB.Close()
	app := server.GetApp()

	// Expose the /metrics endpoint for Prometheus
	// must be defined before MountRoutes
	// must be defined before Prometheus middleware
	app.Get("/metrics", adaptor(promhttp.Handler()))

	// Apply Prometheus middleware before MountMiddlewares
	app.Use(PrometheusMiddleware) // Apply Prometheus middleware

	// Mount middlewares and routes
	server.MountMiddlewares()
	server.MountRoutes(psqlDB)

	routes := app.GetRoutes()

	// Log available routes when initialized
	for _, route := range routes {
		fmt.Printf("%s -> '%s'\n", route.Method, route.Path)
	}

	// Start the server
	server.Start(env.AppEnv.AppPort)
}
