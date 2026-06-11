package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go-backend-task/config"
	"go-backend-task/internal/handler"
	"go-backend-task/internal/logger"
	"go-backend-task/internal/repository"
	"go-backend-task/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	// 1. Load configuration from environment variables
	cfg := config.LoadConfig()

	// 2. Initialize Zap Logger (Structured JSON Logging)
	logger.InitLogger(cfg.LogLevel)
	defer logger.Log.Sync()

	logger.Log.Info("Server is starting...", zap.String("port", cfg.Port))

	// 3. Establish connection pool to PostgreSQL using pgxpool (production-ready connection manager)
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Log.Fatal("Unable to establish connection pool to database", zap.Error(err))
	}
	defer dbpool.Close()

	// Perform ping to verify connection credentials and availability
	if err := dbpool.Ping(ctx); err != nil {
		logger.Log.Fatal("Database ping check failed", zap.Error(err))
	}
	logger.Log.Info("Successfully established connection pool to PostgreSQL")

	// 4. Initialize clean architecture layers (Repository, Controller/Handler)
	repo := repository.NewUserRepository(dbpool)
	userHandler := handler.NewUserHandler(repo)

	// 5. Initialize GoFiber web frame
	app := fiber.New(fiber.Config{
		// Custom global error handler ensuring JSON formatting on server crashes/panic recovery
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// 6. Register routes and associate middlewares
	routes.RegisterRoutes(app, userHandler)

	// 7. Handle Graceful Shutdown (SIGINT, SIGTERM)
	// This ensures in-flight requests complete before connection pool closes.
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-shutdownChan
		logger.Log.Info("Received termination signal. Shutting down server gracefully...")
		if err := app.Shutdown(); err != nil {
			logger.Log.Error("Failed to gracefully shutdown Fiber server", zap.Error(err))
		}
	}()

	// 8. Listen and Serve
	if err := app.Listen(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		logger.Log.Fatal("Fiber server terminated unexpectedly", zap.Error(err))
	}

	logger.Log.Info("Server stopped successfully.")
}
