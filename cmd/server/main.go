package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"money-transfer/config"
	"money-transfer/internal/api/handlers"
	"money-transfer/internal/api/router"
	"money-transfer/internal/service/bank"
	"money-transfer/internal/storage/postgres"

	"github.com/gin-gonic/gin"
)

// @title Money Transfer API
// @version 1.0
// @description API for money transfers between accounts

// @host localhost:8080
// @BasePath /api/v1

// @schemes http
// @produce json
// @consumes json
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)

	// Initialize storage
	store, err := postgres.NewStore(cfg.Database.GetDSN())
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Account().InitializeTestData(context.Background()); err != nil {
		log.Fatal(err)
	}

	// Initialize services
	bankService := bank.NewService(store)

	// Create handlers using factory
	handlersFactory := handlers.NewFactory(bankService)
	appHandlers := handlersFactory.CreateHandlers()

	// Initialize router
	r := router.NewRouter(appHandlers)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Channel for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for termination signal
	<-quit
	log.Println("Shutting down server...")

	// Context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited properly")
}
