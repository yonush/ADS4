package app

import (
	"log"
	"net/http"
	"os"

	"ADS4/internal/config"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// AdminOnly middleware where Admin and Faculty are administrators
// Admin role allows for user account management
func (a *App) AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		role := claims["role"]

		if role != "Admin" && role != "Faculty" {
			return c.Redirect(http.StatusUnauthorized, "/dashboard?error=You%20do%20not%20have%20permission%20to%20access%20this%20page")
		}
		return next(c)
	}
}

// echo response for the keepalive/check if online route
func (a *App) HandeGetShutdown(c echo.Context) error {
	a.DB.Close()
	// Log the shutdown process
	log.Println("Shutting HTTP service down")
	os.Exit(0)
	return nil
}

// force the server to shutdown, for testing purposes - not for production use
func (a *App) HandeGetHello(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		return c.String(http.StatusOK, "OK")
	} else {
		return c.String(http.StatusOK, "")

	}
}

func (a *App) initRoutes() {
	secret := config.LoadConfig().JWTSecret
	// Public routes
	//user account public routes
	a.Router.GET("/", a.HandleGetIndex)
	a.Router.GET("/login", a.HandleGetLogin)
	a.Router.POST("/login", a.HandlePostLogin)
	a.Router.GET("/logout", a.HandleGetLogout)
	a.Router.GET("/register", a.HandleGetRegister)
	a.Router.POST("/register", a.HandlePostRegister)
	a.Router.GET("/forgot-password", a.HandleGetForgotPassword)
	a.Router.POST("/forgot-password", a.HandlePostForgotPassword)
	//other public routes
	a.Router.GET("/hello", a.HandeGetHello)
	a.Router.GET("/examlist", a.HandleGetExamList)
	a.Router.GET("/auth/:examid/:studentid", a.HandleGetStudentAuth)
	a.Router.GET("/exam/:examid/:password", a.HandleGetStudentExam)
	a.Router.POST("/examupload/:studentid/:examid/:password", a.HandlePostExamUpload)
	a.Router.GET("/yearlist", a.HandleGetYearList) //list of available years for the offerings

	//a.Router.GET("/shutdown", a.HandeGetShutdown) //admin only route to shutdown the server, for testing purposes

	// JWT middleware
	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(secret),
		TokenLookup: "cookie:token",
		ErrorHandler: func(c echo.Context, err error) error {
			return c.Redirect(http.StatusSeeOther, "/")
		},
	})

	// Protected routes
	protected := a.Router.Group("")
	protected.Use(jwtMiddleware)

	protected.GET("/dashboard", a.HandleGetDashboard) //index.html

	// Admin-only routes
	admin := protected.Group("")
	admin.Use(a.AdminOnly)
	admin.GET("/admin", a.HandleGetAdmin)
	admin.GET("/shutdown", a.HandeGetShutdown) //admin only route to shutdown the server, for testing purposes
	//admin.GET("/yearlist", a.HandleGetYearList) //list of available years for the offerings
	//Bulk data importer routers
	admin.POST("/importlearners", a.HandlePostImportLearner)
	admin.POST("/importcourses", a.HandlePostImportCourses)
	admin.POST("/importlearnerexams", a.HandlePostImportLearnerExams)
	admin.POST("/importofferings", a.HandlePostImportOfferings)

	// User management CRUD routes
	admin.POST("/api/user", a.HandlePostUser)
	admin.GET("/api/user", a.HandleGetAllUsers)
	admin.GET("/api/user/:username", a.HandleGetUserByUsername)
	admin.PUT("/api/user/:id", a.HandlePutUser)
	admin.DELETE("/api/user/:id", a.HandleDeleteUser)

	//exam offering management CRUD routes
	admin.POST("/api/offering", a.HandlePostOffering)
	admin.GET("/api/offering", a.HandleGetAllOfferings)
	admin.GET("/api/offering/:examid", a.HandleGetOfferingByID)
	admin.PUT("/api/offering/:examid", a.HandlePutOffering)
	admin.DELETE("/api/offering/:examid", a.HandleDeleteOffering)

	//learner exam management CRUD routes
	admin.POST("/api/learnerexam", a.HandlePostLearnerExam)
	admin.GET("/api/learnerexam", a.HandleGetAllLearnerExams)
	admin.GET("/api/learnerexam/:studentid/:examid", a.HandleGetLearnerExamByID)
	admin.PUT("/api/learnerexam/:studentid/:examid", a.HandlePutLearnerExam)
	admin.DELETE("/api/learnerexam/:studentid/:examid", a.HandleDeleteLearnerExam)

}
