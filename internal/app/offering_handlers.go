package app

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"ADS4/internal/models"
	"ADS4/internal/utils"

	"github.com/labstack/echo/v4"
)

//-------------------------------------------------------------------------------------------------------------
//The Handlers below are used for the leaner CRUD interfaces within the ADS and the Assessment Marker Tool

// HandleGetAllOfferings fetches all active exams from the database with optional filtering by status code
// and returns the results as JSON
func (a *App) HandleGetAllOfferings(c echo.Context) error {
	// Check if request if a GET request
	if c.Request().Method != http.MethodGet {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")
	}
	examID := c.QueryParam("examid")
	statusCode := c.QueryParam("status")

	examOfferings, err := a.DB.GetAllOfferings(examID, statusCode)
	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error fetching data", err)
	}

	// Return the results as JSON
	return c.JSON(http.StatusOK, examOfferings)
}

// HandleGetOfferingByID fetches a single exam offering by ID from the database and returns the result as JSON
func (a *App) HandleGetOfferingByID(c echo.Context) error {
	// Check if request if a POST request
	if c.Request().Method != http.MethodGet {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")
	}

	// Get the exam ID from the URL
	examID := c.Param("examid")
	//examID, err := strconv.Atoi(examIDStr)
	if examID != "" {
		return c.Redirect(http.StatusSeeOther, "Invalid exam ID")
	}

	// Fetch the exam from the database
	offering, err := a.DB.GetOfferingByID(examID)

	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error fetching data", err)
	}

	// Return the result as JSON
	return c.JSON(http.StatusOK, offering)
}

func (a *App) HandlePostOffering(c echo.Context) error {
	// Check if request if a GET request
	if c.Request().Method != http.MethodPost {
		return c.Render(http.StatusMethodNotAllowed, "dashboard.html", map[string]interface{}{
			"error": "Method not allowed",
		})
	}
	/*type OfferingsDto struct {
		ExamID        string `json:"examid"` //[year:4][semester:2][coursecode:9]
		Coursecode   string `json:"coursecode"`
		Password      string `json:"password"`
		Status        string `json:"status"` //active,close
		Coordinator   string `json:"coordinator"`
		OwnerID       string `json:"ownerid"`
	}
	*/
	// Parse form data
	examID := c.FormValue("examID")
	coursecode := c.FormValue("coursecode")
	year := c.FormValue("year")
	semester := c.FormValue("semester")
	password := c.FormValue("password")
	status := c.FormValue("status")
	coordinator := c.FormValue("coordinator")
	ownerid := c.FormValue("ownerid")
	duration := c.FormValue("duration")

	// Validate input
	offering, err := validateOffering(examID, coursecode, password, coordinator, ownerid, status, year, semester, duration)
	if err != nil {
		a.handleLogger("Error validating exam offering: " + err.Error())
		// Redirect to dashboard with error message
		return c.Redirect(http.StatusSeeOther, "/dashboard?error="+"Error validating exam offering: "+err.Error())
	}

	// Insert new exma offering exam
	err = a.DB.AddExamOffering(offering)
	if err != nil {
		a.handleLogger("Error adding exam offering: " + err.Error())
		return c.Redirect(http.StatusSeeOther, "/dashboard?error="+err.Error())
	}

	// Redirect to dashboard with success message
	return c.Redirect(http.StatusFound, "/dashboard?message=Exam offering added successfully")
}

func (a *App) HandlePutOffering(c echo.Context) error {
	// Check if request is not a PUT request
	if c.Request().Method != http.MethodPut {
		return c.Render(http.StatusMethodNotAllowed, "dashboard.html", map[string]interface{}{
			"error": "Method not allowed",
		})
	}

	examid := c.Param("examid")
	if utils.IsValidExamCode(examid) == false {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Invalid exam ID",
			"redirectURL": "/dashboard?error=Invalid exam ID",
		})
	}

	// Parse form data from the request body
	var offering models.OfferingsDto
	if err := c.Bind(&offering); err != nil {
		a.handleLogger("Error binding exam offering request body: " + err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Invalid exam offering request body",
			"redirectURL": "/dashboard?error=Invalid exam offering request body",
		})
	}

	// Validate input
	Offering, err := validateOffering(offering.ExamID, offering.CourseCode, offering.Year, offering.Semester,
		offering.Password, offering.Coordinator, offering.OwnerID, offering.Status, offering.Duration)
	if err != nil {
		a.handleLogger("Error validating exam offering: " + err.Error())
		// Redirect to dashboard with error message
		return c.Redirect(http.StatusSeeOther, "/dashboard?error="+"Error validating exam offering: "+err.Error())
	}

	// Update the exam in the database
	err = a.DB.UpdateOffering(Offering)
	if err != nil {
		a.handleLogger("Error updating exam offering: " + err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error updating exam offering: " + err.Error(),
			"redirectURL": "/dashboard?error=" + err.Error()})
	}

	// Redirect to dashboard with success message
	return c.JSON(http.StatusOK, map[string]string{"message": "exam updated successfully", "redirectURL": "/dashboard?message=exam updated successfully"})
}

func validateOffering(examid, coursecode, year, semester, password, coordinator, ownerid, status, duration string) (*models.Offerings, error) {
	const (
		ErrExamIDRequired          string = "exam ID is required"
		ErrCourseCodeRequired      string = "coursecode is required"
		ErrStatusRequired          string = "status code is required"
		ErrYearRequired            string = "year is required"
		ErrSemesterRequired        string = "semester is required"
		ErrYearInvalid             string = "invalid year - must be a number and not in the past"
		ErrDurationInvalid         string = "invalid duration - must be a number 30-240"
		ErrSemesterInvalid         string = "invalid semester - must be S1, S2, or S3"
		ErrExamID                  string = "invalid exam ID"
		ErrStatus                  string = "invalid status code - active/closed"
		ErrStatusRange             string = "status code should be 0 or 1"
		ErrCoordinatorDoesNotExist string = "coordinator does not exist"
		ErrOwneridDoesNotExist     string = "owner ID does not exist"
		ErrCourseCodeTooLong       string = "coursecode is too long, maximum 9 characters"
	)

	var offering models.Offerings

	if examid == "" {
		return &offering, errors.New(ErrExamIDRequired)
	}

	//TODO: db lookup for the coursecode
	if coursecode == "" {
		return &offering, errors.New(ErrCourseCodeRequired)
	}

	if status == "" {
		return &offering, errors.New(ErrStatusRequired)
	}

	if status != "active" && status != "closed" {
		return &offering, errors.New(ErrStatusRange)
	}

	if semester == "" {
		return &offering, errors.New(ErrSemesterRequired)
	}

	if semester != "S1" && semester != "S2" && semester != "S3" {
		return &offering, errors.New(ErrSemesterInvalid)
	}

	if year == "" {
		return &offering, errors.New(ErrYearRequired)
	}

	_year, err := strconv.Atoi(year)
	if err != nil {
		return &offering, errors.New(ErrYearInvalid)
	}

	_duration, err := strconv.Atoi(duration)
	if err != nil {
		return &offering, errors.New(ErrDurationInvalid)
	}

	if _duration < 30 || _duration > 240 {
		return &offering, errors.New(ErrDurationInvalid)
	}
	/*
		currentYear := time.Now().Year()
		if _year < currentYear {
			return &offering, errors.New(ErrYearInvalid)
		}
	*/
	//TODO: include a DB lookup for the faculty
	if coordinator == "" {
		return &offering, errors.New(ErrCoordinatorDoesNotExist)
	}

	//TODO: include a DB lookup for the faculty
	if ownerid == "" {
		return &offering, errors.New(ErrOwneridDoesNotExist)
	}

	//validate the examid
	if utils.IsValidExamCode(examid) == false {
		return &offering, errors.New(ErrExamID)
	}

	//e.g Computer System Architecture
	if len(coursecode) > 9 {
		return &offering, errors.New(ErrCourseCodeTooLong)
	}

	// Set the values of the exam model
	offering.ExamID = sql.NullString{String: examid, Valid: examid != ""}
	offering.Year = _year
	offering.Semester = sql.NullString{String: semester, Valid: semester != ""}
	offering.CourseCode = sql.NullString{String: coursecode, Valid: coursecode != ""}
	offering.Password = sql.NullString{String: password, Valid: password != ""}
	offering.Status = sql.NullString{String: status, Valid: status != ""}
	offering.Coordinator = sql.NullString{String: coordinator, Valid: coordinator != ""}
	offering.OwnerID = sql.NullString{String: ownerid, Valid: ownerid != ""}
	offering.Duration = _duration

	return &offering, nil
}

func (a *App) HandleDeleteOffering(c echo.Context) error {
	// Check if request is not a DELETE request
	if c.Request().Method != http.MethodDelete {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}

	examid := c.Param("examid")
	if utils.IsValidExamCode(examid) == false {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Invalid exam ID",
			"redirectURL": "/dashboard?error=Invalid exam ID",
		})
	}

	// Delete the exam from the database
	err := a.DB.DeleteOffering(examid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":       "Error deleting exam offering",
			"redirectURL": "/dashboard?error=Error deleting exam offering: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message":     "exam deleted successfully",
		"redirectURL": "/dashboard?message=exam deleted successfully",
	})
}

// HandlePutOfferingstatus

func (a *App) HandlePutOfferingStatus(c echo.Context) error {
	// Check if request is not a PUT request
	if c.Request().Method != http.MethodPut {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{
			"error":       "Method not allowed",
			"redirectURL": "/dashboard?error=Method not allowed"})
	}

	// Parse the exam ID from the URL parameter
	examid := c.Param("examid")
	if utils.IsValidExamCode(examid) == false {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Invalid exam ID",
			"redirectURL": "/dashboard?error=Invalid exam ID",
		})
	}

	// Create a struct to bind the JSON request body
	type StatusRequest struct {
		Status string `json:"status"`
	}

	// Bind the JSON request body to the struct
	var req StatusRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Invalid request body",
			"redirectURL": "/dashboard?error=Invalid request body"})
	}

	// Validate offering exists
	offering, err := a.DB.GetOfferingByID(examid)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error":       "Exam offering not found",
			"redirectURL": "/dashboard?error=Exam offering not found"})
	}

	// if no exam found, return 404
	if offering == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error":       "Exam offering not found",
			"redirectURL": "/dashboard?error=Exam offering not found"})
	}

	// Validate status
	if req.Status == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Status is required",
			"redirectURL": "/dashboard?error=Status is required"})
	}

	// Log the incoming data
	a.handleLogger("Exam ID: " + examid)
	a.handleLogger("Status: " + req.Status)

	// Update the exam status in the database
	err = a.DB.UpdateOfferingStatus(examid, req.Status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":       "Failed to update exam offering status",
			"redirectURL": "/dashboard?error=Failed to update exam offering status"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Exam offering status updated successfully"})
}
