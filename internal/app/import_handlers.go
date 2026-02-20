package app

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// retrieve a file uploaded from a web form and save it into the data folder
func (a *App) TransferFile(c echo.Context, destfile string) error {

	//retrieve the uploaded data file - this assumes the form variable datafile
	inf, err := c.FormFile("datafile")
	if err != nil {
		return err
	}

	src, err := inf.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	//copy and save the exam file
	dst, err := os.Create(destfile)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

func (a *App) ProcessImportFile(c echo.Context, target string) error {
	purge := c.FormValue("purge") == "true"
	overwrite := c.FormValue("overwrite") == "true"

	//fmt.Printf("## purge %s, overwrite %s ##", c.FormValue("purge"), c.FormValue("overwrite"))
	//retrieve the uploaded file and transfer it to the data folder
	datafile := a.DataDir + "/" + target + ".csv"
	err := a.TransferFile(c, datafile)
	if err != nil {
		return err
	}
	var ErrNotFound = errors.New("handler not found")

	switch target {
	case "learners":
		err = a.DB.ImportLearners(datafile, purge, overwrite)
	case "courses":
		err = a.DB.ImportCourses(datafile, purge, overwrite)
	case "learnerexams":
		err = a.DB.ImportLearnerExams(datafile, purge, overwrite)
	case "offerings":
		err = a.DB.ImportOfferings(datafile, purge, overwrite)
	default:
		return ErrNotFound
	}

	if err != nil {
		return err
	}

	return nil
}

func (a *App) HandlePostImportLearner(c echo.Context) error {
	// Check if request is a POST request
	if c.Request().Method != http.MethodPost {
		return c.Render(http.StatusMethodNotAllowed, "index.html", map[string]interface{}{
			"error": "Method not allowed"})
	}

	err := a.ProcessImportFile(c, "learners")
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin?error="+"Error processing import file: "+err.Error())
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return c.Render(http.StatusOK, "admin.html", map[string]interface{}{
		"message":       "Learners successfully imported.",
		"username":      claims["username"],
		"role":          claims["role"],
		"email":         claims["email"],
		"user_id":       claims["user_id"],
		"default_admin": claims["default_admin"],
	})
}

func (a *App) HandlePostImportCourses(c echo.Context) error {

	// Check if request is a POST request
	if c.Request().Method != http.MethodPost {
		return c.Render(http.StatusMethodNotAllowed, "index.html", map[string]interface{}{
			"error": "Method not allowed"})

	}

	err := a.ProcessImportFile(c, "courses")
	if err != nil {
		return c.Render(http.StatusSeeOther, "admin.html", map[string]interface{}{
			"error": "Error processing import file: " + err.Error()})

		//return c.Redirect(http.StatusSeeOther, "/admin?error=Error processing import file: "+err.Error())
		//return c.JSON(http.StatusSeeOther, map[string]string{"error": "Error processing import file: " + err.Error()})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return c.Render(http.StatusOK, "admin.html", map[string]interface{}{
		"message":       "Courses successfully imported.",
		"username":      claims["username"],
		"role":          claims["role"],
		"email":         claims["email"],
		"user_id":       claims["user_id"],
		"default_admin": claims["default_admin"],
	})
}

func (a *App) HandlePostImportLearnerExams(c echo.Context) error {
	// Check if request is a POST request
	if c.Request().Method != http.MethodPost {
		return c.Render(http.StatusMethodNotAllowed, "index.html", map[string]interface{}{
			"error": "Method not allowed"})

	}

	err := a.ProcessImportFile(c, "learnerexams")
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin?error="+"Error processing import file: "+err.Error())
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return c.Render(http.StatusOK, "admin.html", map[string]interface{}{
		"message":       "Learner exams successfully imported.",
		"username":      claims["username"],
		"role":          claims["role"],
		"email":         claims["email"],
		"user_id":       claims["user_id"],
		"default_admin": claims["default_admin"],
	})
}

func (a *App) HandlePostImportOfferings(c echo.Context) error {
	// Check if request is a POST request
	if c.Request().Method != http.MethodPost {
		return c.Render(http.StatusMethodNotAllowed, "index.html", map[string]interface{}{
			"error": "Method not allowed"})

	}

	err := a.ProcessImportFile(c, "offerings")
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin?error="+"Error processing import file: "+err.Error())
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return c.Render(http.StatusOK, "admin.html", map[string]interface{}{
		"message":       "Exam offerings successfully imported.",
		"username":      claims["username"],
		"role":          claims["role"],
		"email":         claims["email"],
		"user_id":       claims["user_id"],
		"default_admin": claims["default_admin"],
	})
}
