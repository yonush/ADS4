package app

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

// retrieve a file uploaded from a web form and save it into the data folder
func (a *App) TransferFile(c echo.Context, destfile string) error {

	//retrieve the uploaded data file - this assumes the form variable datafile
	inf, err := c.FormFile("datafile")
	if err != nil {
		a.handleLogger("Error accessing import file data: " + err.Error())
		return err
	}

	src, err := inf.Open()
	if err != nil {
		a.handleLogger("Error opening srouce import file: " + err.Error())
		return err
	}
	defer src.Close()

	//copy and save the exam file
	dst, err := os.Create(destfile)
	if err != nil {
		a.handleLogger("Error creating tagert import file: " + err.Error())
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		a.handleLogger("Error copying source import file to target location: " + err.Error())
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
		a.handleLogger("Transfer error with file import: " + err.Error())
		return err
	}
	var ErrNotFound = errors.New("import handler not found")

	switch target {
	case "learner":
		err = a.DB.ImportLearners(datafile, purge, overwrite)
	case "course":
		err = a.DB.ImportCourses(datafile, purge, overwrite)
	case "learnerexam":
		err = a.DB.ImportLearnerExams(datafile, purge, overwrite)
	case "offering":
		err = a.DB.ImportOfferings(datafile, purge, overwrite)
	default:
		return ErrNotFound
	}

	if err != nil {
		a.handleLogger("Error importing data into database: " + err.Error())
		return err
	}

	return nil
}

func (a *App) HandlePostImport(c echo.Context) error {
	// Check if request is a POST request
	if c.Request().Method != http.MethodPost {
		return c.String(http.StatusSeeOther, "Invalid handler method")
		//return c.Render(http.StatusMethodNotAllowed, "index.html", map[string]interface{}{"error": "Method not allowed"})
	}

	//check of the target file is correct
	target := c.Param("target")
	if target != "learner" && target != "course" && target != "offering" && target != "learnerexam" {
		//return c.Redirect(http.StatusSeeOther, "/admin?error="+"Invalid import target.")
		//return c.JSON(http.StatusSeeOther, map[string]string{"error": "Invalid import import target."})
		return c.String(http.StatusSeeOther, "Invalid import target: "+target)
	}

	err := a.ProcessImportFile(c, target)
	if err != nil {
		//return c.Render(http.StatusSeeOther, "admin.html", map[string]interface{}{"error": "Error processing import file: " + err.Error()})
		//return c.Redirect(http.StatusSeeOther, "/admin?error=Error processing import file: "+err.Error())
		//return c.JSON(http.StatusSeeOther, map[string]string{"error": "Error processing import file: " + err.Error()})
		return c.String(http.StatusSeeOther, "Error processing import file: "+err.Error())
	}
	return c.String(http.StatusOK, "Successful import for : "+target)
}
