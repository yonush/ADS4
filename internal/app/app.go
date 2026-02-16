package app

import (
	"context"
	"log"
	"net/http"
	"os"

	"ADS4/internal/config"
	"ADS4/internal/database"
	"ADS4/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// App holds the application state including database and router
type App struct {
	DB      *database.DB
	Router  *echo.Echo
	Logger  *log.Logger
	Context context.Context
	DataDir string
}

// null route handler for testing
func (a *App) handleNULL(c echo.Context) error {
	a.Logger.Printf("\033[34mNull handler called\033[0m")
	return c.Render(http.StatusOK, "index.html", nil)
}

// handleError is a method of App for handling errors
func (a *App) handleError(c echo.Context, statusCode int, message string, err error) error {
	a.Logger.Printf("\033[31mError: %v\033[0m", err) // Use the logger in the App struct
	return c.JSON(statusCode, map[string]string{"error": message})
}

func (a *App) handleLogger(message string) {

	a.Logger.Printf("\033[34m%s\033[0m", message)
}

// NewApp creates a new instance of App
func NewApp(cfg config.Config) *App {
	// Initialize Echo
	router := echo.New()

	// Set up renderer
	renderer, err := utils.NewTemplateRenderer()
	if err != nil {
		// Handle the error, e.g.:
		panic(err)
	}

	router.Renderer = renderer
	utils.NewLogger() // new

	// Serve static files
	router.Static("/static", "static")

	router.Use(middleware.RequestLogger()) // Log requests
	router.Use(middleware.Recover())       // Recover from panics
	router.Use(middleware.CORS())          // Enable CORS
	router.Use(utils.LoggingMiddleware)

	// Initialize Database
	db, err := database.NewDB(cfg)
	if err != nil {
		panic(err)
	}

	// Seed database
	// Initialize database and seed data if needed
	if err := database.SeedDatabase(db.DB); err != nil {
		panic(err)
	}

	// Initialize Logger
	logger := log.New(os.Stdout, "\033[34mAPP: \033[0m", log.LstdFlags)

	app := &App{
		DB:      db,
		Router:  router,
		Logger:  logger,
		DataDir: cfg.DataDir,
	}

	// Initialize routes
	app.initRoutes()

	return app
}
