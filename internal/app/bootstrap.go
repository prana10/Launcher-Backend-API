package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/lib/pq"
	"github.com/swaggo/fiber-swagger"
	_ "github.com/swaggo/fiber-swagger" 
	"go.uber.org/fx"

	"launcherbackend_api/internal/config"
	_ "launcherbackend_api/internal/delivery/http/docs" 
	"launcherbackend_api/internal/delivery/http/handle"
	"launcherbackend_api/internal/repository"
	"launcherbackend_api/internal/usecase"
)

var Module = fx.Options(
	fx.Provide(
		config.LoadConfig,
		NewFiberApp,
		ProvideDatabaseConnection,
		ProvideRepositories,
		ProvideUseCases,
		ProvideHandlers,
	),
	fx.Invoke(RegisterRoutes),
)

func NewFiberApp(cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Default error handler
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Configure middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	return app
}

func ProvideDatabaseConnection(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DBConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to database successfully")
	return db, nil
}

func ProvideRepositories(db *sql.DB) *repository.Repositories {
	return &repository.Repositories{
		OTA: repository.NewPostgresOTARepository(db),
	}
}

func ProvideUseCases(repos *repository.Repositories) *usecase.UseCases {
	return &usecase.UseCases{
		OTA: usecase.NewOTAUseCase(repos.OTA),
	}
}

func ProvideHandlers(useCases *usecase.UseCases) *handle.Handlers {
	return &handle.Handlers{
		OTA: handle.NewOTAHandler(useCases.OTA),
	}
}

func RegisterRoutes(app *fiber.App, handlers *handle.Handlers) {
	api := app.Group("/api/v1")
	handlers.OTA.RegisterRoutes(api)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})
	
	// Swagger UI endpoint
	app.Get("/swagger/*", fiberSwagger.WrapHandler)
}

func StartApp() {
	fx.New(
		Module,
		fx.Invoke(func(app *fiber.App, cfg *config.Config, lc fx.Lifecycle) {
			addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
			
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					log.Printf("Starting server on %s in %s mode", addr, cfg.Environment)
					log.Printf("Swagger documentation available at http://%s/swagger/index.html", addr)
					go func() {
						if err := app.Listen(addr); err != nil {
							log.Fatalf("Failed to start server: %v", err)
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					log.Println("Shutting down server")
					return app.Shutdown()
				},
			})
		}),
	).Run()
} 