package app

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// HandleGetDashboard serves the dashboard page
func (a *App) HandlePostImportLearner(c echo.Context) error {
	// Check if request if a POST request
	if c.Request().Method != http.MethodPost {
		return c.Render(http.StatusMethodNotAllowed, "index.html", map[string]interface{}{
			"error": "Method not allowed",
		})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return c.Render(http.StatusOK, "admin.html", map[string]interface{}{
		"username":      claims["username"],
		"role":          claims["role"],
		"email":         claims["email"],
		"user_id":       claims["user_id"],
		"default_admin": claims["default_admin"],
	})
}

// HandleGetDashboard serves the dashboard page
func (a *App) HandlePostImportCourses(c echo.Context) error {
	// Check if request if a POST request
	if c.Request().Method != http.MethodPost {
		return c.Render(http.StatusMethodNotAllowed, "index.html", map[string]interface{}{
			"error": "Method not allowed",
		})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return c.Render(http.StatusOK, "admin.html", map[string]interface{}{
		"username":      claims["username"],
		"role":          claims["role"],
		"email":         claims["email"],
		"user_id":       claims["user_id"],
		"default_admin": claims["default_admin"],
	})
}

// HandleGetDashboard serves the dashboard page
func (a *App) HandlePostImportLearnerExams(c echo.Context) error {
	// Check if request if a POST request
	if c.Request().Method != http.MethodPost {
		return c.Render(http.StatusMethodNotAllowed, "index.html", map[string]interface{}{
			"error": "Method not allowed",
		})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return c.Render(http.StatusOK, "admin.html", map[string]interface{}{
		"username":      claims["username"],
		"role":          claims["role"],
		"email":         claims["email"],
		"user_id":       claims["user_id"],
		"default_admin": claims["default_admin"],
	})
}

// HandleGetDashboard serves the dashboard page
func (a *App) HandlePostImportOfferings(c echo.Context) error {
	// Check if request if a POST request
	if c.Request().Method != http.MethodPost {
		return c.Render(http.StatusMethodNotAllowed, "index.html", map[string]interface{}{
			"error": "Method not allowed",
		})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return c.Render(http.StatusOK, "admin.html", map[string]interface{}{
		"username":      claims["username"],
		"role":          claims["role"],
		"email":         claims["email"],
		"user_id":       claims["user_id"],
		"default_admin": claims["default_admin"],
	})
}
