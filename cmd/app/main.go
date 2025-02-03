package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
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
	// Dynamic instance label
	instance string

	// HTTP metrics
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "route", "status", "container_id", "env"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route", "container_id", "env"},
	)

	// Application metrics
	goroutinesGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "go_goroutines_current",
			Help: "Current number of goroutines",
		},
		[]string{"container_id", "env"},
	)

	threadsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "go_threads_current",
			Help: "Current number of OS threads",
		},
		[]string{"container_id", "env"},
	)

	// Memory metrics
	memAllocGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "go_memory_alloc_bytes",
			Help: "Current memory usage in bytes",
		},
		[]string{"container_id", "env"},
	)

	memTotalAllocGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "go_memory_total_alloc_bytes",
			Help: "Total allocated memory in bytes",
		},
		[]string{"container_id", "env"},
	)

	// GC metrics
	gcPauseGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "go_gc_pause_ns",
			Help: "Last GC pause time in nanoseconds",
		},
		[]string{"container_id", "env"},
	)
)

func init() {
	// Get the instance name (fallback to hostname if ENV is not set)
	instance = os.Getenv("HOSTNAME")
	if instance == "" {
		instance = "unknown"
	}

	// Register custom metrics
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration, goroutinesGauge, threadsGauge, memAllocGauge, memTotalAllocGauge, gcPauseGauge)
}

func PrometheusRuntimeMetrics() {
	var memStats runtime.MemStats
	for {
		runtime.ReadMemStats(&memStats)

		goroutinesGauge.WithLabelValues(instance, env.AppEnv.AppEnv).Set(float64(runtime.NumGoroutine()))
		threadsGauge.WithLabelValues(instance, env.AppEnv.AppEnv).Set(float64(runtime.GOMAXPROCS(0)))
		memAllocGauge.WithLabelValues(instance, env.AppEnv.AppEnv).Set(float64(memStats.Alloc))
		memTotalAllocGauge.WithLabelValues(instance, env.AppEnv.AppEnv).Set(float64(memStats.TotalAlloc))
		gcPauseGauge.WithLabelValues(instance, env.AppEnv.AppEnv).Set(float64(memStats.PauseNs[(memStats.NumGC+255)%256]))

		time.Sleep(time.Second)
	}
}

// Prometheus Middleware to track HTTP request metrics
func PrometheusMiddleware(c *fiber.Ctx) error {
	start := time.Now()
	method := c.Method()
	route := c.Path()

	err := c.Next() // Call the next handler in the chain

	// Record metrics after request completes
	duration := time.Since(start).Seconds()
	status := c.Response().StatusCode()

	// Track request metrics
	httpRequestsTotal.WithLabelValues(method, route, fmt.Sprint(status), instance, env.AppEnv.AppEnv).Inc()
	httpRequestDuration.WithLabelValues(method, route, instance, env.AppEnv.AppEnv).Observe(duration)

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

	// Start runtime metrics collection
	go PrometheusRuntimeMetrics()

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
