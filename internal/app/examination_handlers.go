package app

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

/*
	Handlers for all the non CRUD examination related actions
	used by:
	- dashboard - HandleGetYearList, HandleExamMetrics
	- reporting
	- AMT
*/

func (a *App) HandleGetYearList(c echo.Context) error {
	// Check if request if a GET request
	if c.Request().Method != http.MethodGet {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")
	}

	//retrieve the list of years in the database
	examyears, err := a.DB.GetExamYears()
	if err != nil {
		a.handleLogger("Error fetching exam year data: " + err.Error())
		return a.handleError(c, http.StatusInternalServerError, "Error fetching data", err)
	}

	// Return the results as JSON
	return c.JSON(http.StatusOK, examyears)
}

// e.g. /exammetrics?year=2025&semester=S1
func (a *App) HandleExamMetrics(c echo.Context) error {
	// Check if request if a GET request
	if c.Request().Method != http.MethodGet {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")
	}

	//year := c.QueryParam("year") //passed as ?year=
	//year := c.Param("year") //path parameter

	//assume the current year else use the argument
	year := strconv.Itoa(time.Now().Year())
	_year := c.QueryParam("year")
	if _year != "" {
		year = _year
	}

	//assume the first semester else use the argument
	semester := "S1"
	_semester := c.QueryParam("semester")
	if _semester != "" {
		semester = _semester
	}

	//retrieve the list of active exam offerings with current metrics
	metrics, err := a.DB.GetExamByYearSemester(year, semester)
	if err != nil {
		a.handleLogger("Error fetching exam data: " + err.Error())
		return a.handleError(c, http.StatusInternalServerError, "Error fetching data", err)
	}

	// Return the results as JSON
	return c.JSON(http.StatusOK, metrics)

}

// e.g. /closedexams/:field/:value/:semester
func (a *App) HandleClosedExams(c echo.Context) error {
	// Check if request if a GET request
	if c.Request().Method != http.MethodGet {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")
	}

	//year := c.QueryParam("year") //passed as ?year=
	//year := c.Param("year") //path parameter
	field := c.Param("field")
	value := c.Param("value")
	semester := c.Param("semester")

	if semester != "S1" && semester != "S2" && semester != "S3" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid or missing semester value",
		})
	}
	if field != "student" && field != "course" && field != "examid" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid or missing field",
		})
	}

	if value == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid or missing data value",
		})
	}

	//fmt.Printf("%s, %s, %s\n", field, value, semester)
	metrics, err := a.DB.GetExaminations(field, value, semester)
	if err != nil {
		a.handleLogger("Error fetchign examination data: ")
		return a.handleError(c, http.StatusInternalServerError, "Error fetching data", err)
	}

	// Return the results as JSON
	return c.JSON(http.StatusOK, metrics)

}
