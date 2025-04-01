package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"github.com/spanwalla/url-shortener/config"
	_ "github.com/spanwalla/url-shortener/docs"
	httpApi "github.com/spanwalla/url-shortener/internal/controller/http"
	"github.com/spanwalla/url-shortener/internal/repository"
	"github.com/spanwalla/url-shortener/internal/service"
	"github.com/spanwalla/url-shortener/pkg/encoder"
	"github.com/spanwalla/url-shortener/pkg/httpserver"
	"github.com/spanwalla/url-shortener/pkg/postgres"
	"github.com/spanwalla/url-shortener/pkg/validator"
)

// @title URL Shortener
// @version 1.0

// @host localhost:8080
// @BasePath /

// Run creates objects via constructors
func Run() {
	// Config
	configPath, ok := os.LookupEnv("CONFIG_PATH")
	if !ok || len(configPath) == 0 {
		log.Fatal("app - os.LookupEnv: CONFIG_PATH is empty")
	}

	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatal(fmt.Errorf("app - config.New: %w", err))
	}

	// Logger
	initLogger(cfg.Level)
	log.Info("Config read")

	// Postgres
	log.Info("Connecting to postgres...")
	pg, err := postgres.New(cfg.URL, postgres.MaxPoolSize(cfg.PoolMax))
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Services and repos
	log.Info("Initializing services and repos...")
	services := service.New(service.Dependencies{
		Repos:               repository.NewPostgresRepositories(pg),
		Encoder:             encoder.NewRandom(cfg.Alphabet, time.Now().UnixNano()),
		AliasLength:         cfg.Length,
		AttemptsOnCollision: cfg.Attempts,
	})

	// Echo handler
	log.Info("Initializing handlers and routes...")
	handler := echo.New()
	handler.Validator = validator.NewCustomValidator()
	httpApi.ConfigureRouter(handler, services)

	// HTTP Server
	log.Info("Starting HTTP server...")
	log.Debugf("Server port: %s", cfg.Port)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.Port))

	log.Info("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Errorf("app - Run - httpServer.Notify: %v", err)
	}

	// Graceful shutdown
	log.Info("Shutting down...")

	err = httpServer.Shutdown()
	if err != nil {
		log.Errorf("app - Run - httpServer.Shutdown: %v", err)
	}
}
