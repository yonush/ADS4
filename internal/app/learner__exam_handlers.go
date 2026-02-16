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

// HandleGetAllLearnerExam fetches all active exams from the database with optional filtering by status code
// and returns the results as JSON
func (a *App) HandleGetAllLearnerExams(c echo.Context) error {
	// Check if request if a POST request
	if c.Request().Method != http.MethodGet {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")
	}
	learnerexamid := c.QueryParam("learnerexamID")
	statusCode := c.QueryParam("status")

	// Get the LearnerExam ID from the URL
	learnerexamID, err := strconv.Atoi(learnerexamid)
	if err != nil {
		a.handleLogger("Error converting leaner ID to integer: " + err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid learner exam ID",
			"redirectURL": "/dashboard?error=Invalid learner exam ID"})
	}

	LearnerExams, err := a.DB.GetAllLearnerExams(learnerexamID, statusCode)
	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error fetching leaner exam data", err)
	}

	// Return the results as JSON
	return c.JSON(http.StatusOK, LearnerExams)
}

// HandleGetOfferingByID fetches a single exam offering by ID from the database and returns the result as JSON
func (a *App) HandleGetLearnerExamByID(c echo.Context) error {
	// Check if request if a POST request
	if c.Request().Method != http.MethodGet {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")
	}

	// Get the LearnerExam ID from the URL
	id := c.Param("learnerexamID")
	learnerexamID, err := strconv.Atoi(id)
	if err != nil {
		a.handleLogger("Error converting leaner ID to integer: " + err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid learner exam ID",
			"redirectURL": "/dashboard?error=Invalid learner exam ID"})
	}

	// Fetch the learner from the database
	learnerexam, err := a.DB.GetLearnerExamByID(learnerexamID)

	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error fetching data", err)
	}

	// Return the result as JSON
	return c.JSON(http.StatusOK, learnerexam)
}

func (a *App) HandlePostLearnerExam(c echo.Context) error {
	// Check if request if a GET request
	if c.Request().Method != http.MethodPost {
		return c.Render(http.StatusMethodNotAllowed, "dashboard.html", map[string]interface{}{
			"error": "Method not allowed",
		})
	}

	learnerexamID := c.FormValue("learnerexamID")
	studentid := c.FormValue("studentid")
	examid := c.FormValue("examid")
	status := c.FormValue("status")

	a.handleLogger("Learner ID: " + learnerexamID)
	a.handleLogger("Student ID: " + studentid)
	a.handleLogger("Exam ID: " + examid)
	a.handleLogger("Status: " + status)

	// Validate input
	learnerexam, err := validateLearnerExam(learnerexamID, studentid, examid, status)
	if err != nil {
		a.handleLogger("Error validating leaner exam details " + err.Error())
		// Redirect to dashboard with error message
		return c.Redirect(http.StatusSeeOther, "/dashboard?error="+"Error validating learner exam details: "+err.Error())
	}

	// Insert new exma offering LearnerExam
	err = a.DB.AddLearnerExam(learnerexam)
	if err != nil {
		a.handleLogger("Error adding learner exam details: " + err.Error())
		return c.Redirect(http.StatusSeeOther, "/dashboard?error="+err.Error())
	}

	// Redirect to dashboard with success message
	return c.Redirect(http.StatusFound, "/dashboard?message=Learner exam details added successfully")
}

func (a *App) HandlePutLearnerExam(c echo.Context) error {
	// Check if request is not a PUT request
	if c.Request().Method != http.MethodPut {
		return c.Render(http.StatusMethodNotAllowed, "dashboard.html", map[string]interface{}{
			"error": "Method not allowed",
		})
	}

	// Get the LearnerExam ID from the URL
	id := c.Param("learnerexamID")
	_, err := strconv.Atoi(id)
	if err != nil {
		a.handleLogger("Error converting leaner ID to integer: " + err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid learner ID",
			"redirectURL": "/dashboard?error=Invalid learner ID"})
	}

	// Parse form data from the request body
	var learnerexam models.LearnerExamDto
	if err := c.Bind(&learnerexam); err != nil {
		a.handleLogger("Error binding learner exam details request body: " + err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Invalid learner exam details request body",
			"redirectURL": "/dashboard?error=Invalid learner exam details request body",
		})
	}

	// Log the incoming data
	a.handleLogger("Learner ID: " + learnerexam.LearnerExamID)
	a.handleLogger("Student ID: " + learnerexam.StudentID)
	a.handleLogger("Exam ID: " + learnerexam.ExamID)
	a.handleLogger("Status: " + learnerexam.Status)

	// Validate input
	learnerExam, err := validateLearnerExam(learnerexam.LearnerExamID, learnerexam.StudentID, learnerexam.ExamID, learnerexam.Status)
	if err != nil {
		a.handleLogger("Error validating learner exam details: " + err.Error())
		// Redirect to dashboard with error message
		return c.Redirect(http.StatusSeeOther, "/dashboard?error="+"Error validating learner exam details: "+err.Error())
	}

	// Update the LearnerExam in the database
	err = a.DB.UpdateLearnerExam(learnerExam)
	if err != nil {
		a.handleLogger("Error updating learner exam details: " + err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error updating learner exam details: " + err.Error(),
			"redirectURL": "/dashboard?error=" + err.Error()})
	}

	// Redirect to dashboard with success message
	return c.JSON(http.StatusOK, map[string]string{"message": "LearnerExam updated successfully", "redirectURL": "/dashboard?message=LearnerExam updated successfully"})
}

func validStatus(status string) bool {
	stat := utils.StatusSet{}

	stat.Add("available")
	stat.Add("active")
	stat.Add("expired")
	stat.Add("closed")
	stat.Add("marked")
	return stat.Has(status)
}

func validateLearnerExam(learnerexamID, studentid, examid, status string) (*models.LearnerExam, error) {
	const (
		ErrlearnerexamIDRequired string = "exam ID is required"
		ErrStudentIDRequired     string = "student ID is required"
		ErrStatusRequired        string = "status is required"
		ErrlearnerexamID         string = "invalid learner exam ID"
		ErrexamID                string = "invalid exam ID"
		ErrStatus                string = "invalid status code "
		ErrStudentIDTooLong      string = "student ID length exceeded"
	)

	var learnerexam models.LearnerExam

	if learnerexamID == "" {
		return &learnerexam, errors.New(ErrlearnerexamIDRequired)
	}

	// Get the LearnerExam ID from the URL
	learnerexamid, err := strconv.Atoi(learnerexamID)
	if err != nil {
		return &learnerexam, errors.New(ErrlearnerexamID)
	}

	if studentid == "" {
		return &learnerexam, errors.New(ErrStudentIDRequired)
	}

	if status == "" {
		return &learnerexam, errors.New(ErrStatusRequired)
	}

	if validStatus(status) == false {
		return &learnerexam, errors.New(ErrStatus)
	}

	if len(studentid) > 8 {
		return &learnerexam, errors.New(ErrStudentIDTooLong)
	}

	if utils.IsValidCourseCode(examid) == false {
		return &learnerexam, errors.New(ErrexamID)
	}

	// Set the values of the LearnerExam model
	// Initialize sql.NullString for optional fields
	learnerexam.LearnerExamID = learnerexamid
	learnerexam.StudentID = sql.NullString{String: studentid, Valid: true}
	learnerexam.ExamID = sql.NullString{String: examid, Valid: true}
	learnerexam.Status = sql.NullString{String: status, Valid: true}

	return &learnerexam, nil
}

// TODO: continue here ##
func (a *App) HandleDeleteLearnerExam(c echo.Context) error {
	// Check if request is not a DELETE request
	if c.Request().Method != http.MethodDelete {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}

	// Get the LearnerExam ID from the URL
	id := c.Param("learnerexamID")
	learnerexamID, err := strconv.Atoi(id)
	if err != nil {
		a.handleLogger("Error converting leaner ID to integer: " + err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid learner exam ID",
			"redirectURL": "/dashboard?error=Invalid learner exam ID"})
	}

	// Delete the LearnerExam from the database
	err = a.DB.DeleteLearnerExam(learnerexamID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":       "Error deleting learner exam",
			"redirectURL": "/dashboard?error=Error deleting learner exam " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message":     "Learner exam deleted successfully",
		"redirectURL": "/dashboard?message=Learner exam deleted successfully",
	})
}

func (a *App) HandlePutLearnerExamStatus(c echo.Context) error {
	// Check if request is not a PUT request
	if c.Request().Method != http.MethodPut {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{
			"error":       "Method not allowed",
			"redirectURL": "/dashboard?error=Method not allowed"})
	}

	// Get the LearnerExam ID from the URL
	learnerexamid := c.Param("learnerexamID")
	learnerexamID, err := strconv.Atoi(learnerexamid)
	if err != nil {
		a.handleLogger("Error converting leaner ID to integer: " + err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid learner exam ID",
			"redirectURL": "/dashboard?error=Invalid learner exam ID"})
	}

	// Create a struct to bind the JSON request body
	type StatusRequest struct {
		Status string `json:"status"`
	}

	// Bind the JSON request body to the struct
	var req StatusRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Invalid learner exam request body",
			"redirectURL": "/dashboard?error=Invalid learner exam request body"})
	}

	// Validate offering exists
	offering, err := a.DB.GetLearnerExamByID(learnerexamID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error":       "Learner exam not found",
			"redirectURL": "/dashboard?error=Learner exam not found"})
	}

	// if no LearnerExam found, return 404
	if offering == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error":       "Learner exam not found",
			"redirectURL": "/dashboard?error=Learner exam not found"})
	}

	// Validate status exists
	if req.Status == "" || validStatus(req.Status) == false {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Status code is required",
			"redirectURL": "/dashboard?error=Status code is required"})
	}

	// Validate status is valid value
	if validStatus(req.Status) == false {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Status code is invalid - must be one of available, active, expired, closed, marked",
			"redirectURL": "/dashboard?error=Status code is invalid - must be one of available, active, expired, closed, marked"})
	}

	// Log the incoming data
	a.handleLogger("Learner exam ID: " + learnerexamid)
	a.handleLogger("Status: " + req.Status)

	// Update the LearnerExam status in the database
	err = a.DB.UpdateLearnerExamStatus(learnerexamID, req.Status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":       "Failed to update the learner exam status",
			"redirectURL": "/dashboard?error=Failed to update the learner exam status"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Exam offering status updated successfully"})
}
