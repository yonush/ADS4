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

/* Courses
   - CourseCode e.g ITCS5.100
   - Level - 1-9
   - description
   - status - active,closed
*/

// HandleGetAllCourses fetches all active exams from the database with optional filtering by status code
// and returns the results as JSON
func (a *App) HandleGetAllCourses(c echo.Context) error {
	// Check if request if a POST request
	if c.Request().Method != http.MethodGet {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")
	}
	coursecode := c.QueryParam("coursecode")
	statusCode := c.QueryParam("status")

	examCourses, err := a.DB.GetAllCourses(coursecode, statusCode)
	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error fetching data", err)
	}

	// Return the results as JSON
	return c.JSON(http.StatusOK, examCourses)
}

// HandleGetCourseByID fetches a single exam Course by ID from the database and returns the result as JSON
func (a *App) HandleGetCourseByID(c echo.Context) error {
	// Check if request if a POST request
	if c.Request().Method != http.MethodGet {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")
	}

	// Get the exam ID from the URL
	coursecode := c.Param("coursecode")
	//coursecode, err := strconv.Atoi(coursecodeStr)
	if coursecode != "" {
		return c.Redirect(http.StatusSeeOther, "Invalid exam ID")
	}

	// Fetch the exam from the database
	Course, err := a.DB.GetCourseByID(coursecode)

	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error fetching data", err)
	}

	// Return the result as JSON
	return c.JSON(http.StatusOK, Course)
}

func (a *App) HandlePostCourse(c echo.Context) error {
	// Check if request if a GET request
	if c.Request().Method != http.MethodPost {
		return c.Render(http.StatusMethodNotAllowed, "dashboard.html", map[string]interface{}{
			"error": "Method not allowed",
		})
	}
	/*type CoursesDto struct {
	CourseCode  string `json:"coursecode"`
	Description string `json:"description"`
	Level       string    `json:"level"`
	Status      string    `json:"status"` //True,False or 1,0
	}
	*/
	// Parse form data
	coursecode := c.FormValue("coursecode")
	level := c.FormValue("level")
	description := c.FormValue("description")
	status := c.FormValue("status")

	a.handleLogger("Course code: " + coursecode)
	a.handleLogger("Description: " + description)
	a.handleLogger("Level: " + level)
	a.handleLogger("Status: " + status)

	// Validate input
	Course, err := validateCourse(coursecode, description, level, status)
	if err != nil {
		a.handleLogger("Error validating exam Course: " + err.Error())
		// Redirect to dashboard with error message
		return c.Redirect(http.StatusSeeOther, "/dashboard?error="+"Error validating course: "+err.Error())
	}

	// Insert new exma Course exam
	err = a.DB.AddExamCourse(Course)
	if err != nil {
		a.handleLogger("Error adding course: " + err.Error())
		return c.Redirect(http.StatusSeeOther, "/dashboard?error="+err.Error())
	}

	// Redirect to dashboard with success message
	return c.Redirect(http.StatusFound, "/dashboard?message=Course added successfully")
}

func (a *App) HandlePutCourse(c echo.Context) error {
	// Check if request is not a PUT request
	if c.Request().Method != http.MethodPut {
		return c.Render(http.StatusMethodNotAllowed, "dashboard.html", map[string]interface{}{
			"error": "Method not allowed",
		})
	}

	coursecode := c.Param("coursecode")
	if utils.IsValidCourseCode(coursecode) == false {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Invalid course code",
			"redirectURL": "/dashboard?error=Invalid course code",
		})
	}

	// Parse form data from the request body
	var coursedto models.CoursesDto
	if err := c.Bind(&coursedto); err != nil {
		a.handleLogger("Error binding course request body: " + err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Invalid course request body",
			"redirectURL": "/dashboard?error=Invalid exam Course request body",
		})
	}

	a.handleLogger("Course code: " + coursedto.CourseCode)
	a.handleLogger("Description: " + coursedto.Description)
	a.handleLogger("Level: " + coursedto.Level)
	a.handleLogger("Status: " + coursedto.Status)

	// Validate input
	course, err := validateCourse(coursedto.CourseCode, coursedto.Description, coursedto.Level, coursedto.Status)
	if err != nil {
		a.handleLogger("Error validating Course: " + err.Error())
		// Redirect to dashboard with error message
		return c.Redirect(http.StatusSeeOther, "/dashboard?error="+"Error validating Course: "+err.Error())
	}

	// Update the exam in the database
	err = a.DB.UpdateCourse(course)
	if err != nil {
		a.handleLogger("Error updating exam Course: " + err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error updating exam Course: " + err.Error(),
			"redirectURL": "/dashboard?error=" + err.Error()})
	}

	// Redirect to dashboard with success message
	return c.JSON(http.StatusOK, map[string]string{"message": "exam updated successfully", "redirectURL": "/dashboard?message=exam updated successfully"})
}

func validateCourse(coursecode, description, level, status string) (*models.Courses, error) {
	const (
		ErrCourseCodeRequired      string = "coursecode is required"
		ErrLevelRequired           string = "level code is required"
		ErrDescriptionRequired     string = "description is required"
		Errcoursecode              string = "invalid course code"
		ErrStatusRequired          string = "status code is required"
		ErrStatusRange             string = "invalid status code - active/close"
		ErrLevelInvalid            string = "invalid level type"
		ErrLevelRange              string = "invalid level range - 1 to 9"
		ErrCoordinatorDoesNotExist string = "coordinator does not exist"
		ErrOwneridDoesNotExist     string = "owner ID does not exist"
		ErrCourseCodeTooLong       string = "coursecode is too long, maximum 9 characters"
		ErrDescriptionTooLong      string = "description is too long, maximum 255 characters"
	)

	var course models.Courses

	//TODO: db lookup for the coursecode
	if coursecode == "" {
		return &course, errors.New(ErrCourseCodeRequired)
	}

	//validate the coursecode - this function should be implemented in the utils package and should be more robust
	if utils.IsValidCourseCode(coursecode) == false {
		return &course, errors.New(Errcoursecode)
	}

	if len(coursecode) > 9 {
		return &course, errors.New(ErrCourseCodeTooLong)
	}

	if description == "" {
		return &course, errors.New(ErrDescriptionRequired)
	}

	if len(description) > 255 {
		return &course, errors.New(ErrDescriptionTooLong)
	}

	if status == "" {
		return &course, errors.New(ErrStatusRequired)
	}

	if status != "active" && status != "closed" {
		return &course, errors.New(ErrStatusRange)
	}

	if level == "" {
		return &course, errors.New(ErrLevelRequired)
	}

	lvl, err := strconv.Atoi(level)
	if err != nil {
		return &course, errors.New(ErrLevelInvalid)
	}

	if lvl < 0 || lvl > 9 {
		return &course, errors.New(ErrLevelRange)
	}

	// Set the values of the exam model
	course.CourseCode = sql.NullString{String: coursecode, Valid: coursecode != ""}
	course.Description = sql.NullString{String: description, Valid: description != ""}
	course.Status = sql.NullString{String: status, Valid: status != ""}
	course.Level = lvl

	return &course, nil
}

func (a *App) HandleDeleteCourse(c echo.Context) error {
	// Check if request is not a DELETE request
	if c.Request().Method != http.MethodDelete {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}

	coursecode := c.Param("coursecode")
	if utils.IsValidCourseCode(coursecode) == false {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Invalid exam ID",
			"redirectURL": "/dashboard?error=Invalid exam ID",
		})
	}

	// Delete the exam from the database
	err := a.DB.DeleteCourse(coursecode)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":       "Error deleting exam Course",
			"redirectURL": "/dashboard?error=Error deleting exam Course: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message":     "exam deleted successfully",
		"redirectURL": "/dashboard?message=exam deleted successfully",
	})
}

// HandlePutCoursestatus

func (a *App) HandlePutCourseStatus(c echo.Context) error {
	// Check if request is not a PUT request
	if c.Request().Method != http.MethodPut {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{
			"error":       "Method not allowed",
			"redirectURL": "/dashboard?error=Method not allowed"})
	}

	// Parse the exam ID from the URL parameter
	coursecode := c.Param("coursecode")
	if utils.IsValidCourseCode(coursecode) == false {
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

	// Validate Course exists
	Course, err := a.DB.GetCourseByID(coursecode)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error":       "Exam Course not found",
			"redirectURL": "/dashboard?error=Exam Course not found"})
	}

	// if no exam found, return 404
	if Course == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error":       "Exam Course not found",
			"redirectURL": "/dashboard?error=Exam Course not found"})
	}

	// Validate status
	if req.Status == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Status is required",
			"redirectURL": "/dashboard?error=Status is required"})
	}

	// Log the incoming data
	a.handleLogger("Exam ID: " + coursecode)
	a.handleLogger("Status: " + req.Status)

	// Update the exam status in the database
	err = a.DB.UpdateCourseStatus(coursecode, req.Status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":       "Failed to update exam Course status",
			"redirectURL": "/dashboard?error=Failed to update exam Course status"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Exam Course status updated successfully"})
}
