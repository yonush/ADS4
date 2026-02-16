package app

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

//----------------------------------------------------------------------------------------------------------
// The handler below are used by the Z2A Assessment tool

// HandleGetAllOfferings fetches a curated list of all active exams from the database and returns the results as JSON
// list includes the course description from a join with the courses table
func (a *App) HandleGetYearList(c echo.Context) error {
	// Check if request if a GET request
	if c.Request().Method != http.MethodGet {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")
	}

	//retrieve the list of years in the database
	examYears, err := a.DB.GetExamYears()
	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error fetching data", err)
	}

	// Return the results as JSON
	return c.JSON(http.StatusOK, examYears)
	
}

// HandleGetAllOfferings fetches a curated list of all active exams from the database and returns the results as JSON
// list includes the course description from a join with the courses table
func (a *App) HandleGetExamList(c echo.Context) error {
	// Check if request if a GET request
	if c.Request().Method != http.MethodGet {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")
	}

	//retrieve the list of active exams for the current year only
	currentyear := strconv.Itoa(time.Now().Year())
	examOfferings, err := a.DB.GetActiveExams(currentyear)
	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error fetching data", err)
	}

	// Return the results as JSON
	return c.JSON(http.StatusOK, examOfferings)
}

// POST /examupload/:studentid/:examid/:password
// HandleGetStudentAuth checks if the learner is permitted to engage in the chosen exam identified by examid
// returns the exam password if correct
// TODO: include an auth token and update the learner exam status after saving the exam
func (a *App) HandlePostExamUpload(c echo.Context) error {
	// Check if request if a POST request
	if c.Request().Method != http.MethodPost {
		return c.JSON(http.StatusMethodNotAllowed, map[string]interface{}{
			"error": "Method not allowed",
		})
	}
	password := c.Param("password")
	examid := c.Param("examid")
	studentid := c.Param("studentid")

	//TODO do a time check to set the exam state to expired if the learner does not save within the time

	pass, err := a.DB.GetExamPassword(examid, studentid)
	if pass == "" || pass != password || err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"Status": "Error"})
	}

	// check if the exam is still valid by checking the state and elapsed time
	//return if the time has expired, we assume the exam was started
	if a.DB.CheckIfTime(examid, studentid) == false {
		//make sure the exam is recorded as expired and not closed then return
		if a.DB.CloseLearnerExam(studentid, examid, true) != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{"Status": "Error", "Message": "Unable to set the exam status"})
		}
		return c.JSON(http.StatusBadRequest, map[string]any{"Status": "Error", "Message": "Exam has expired"})
	}

	//retrieve the uploaded exam file
	examfile, err := c.FormFile("exam")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"Status": "Error", "Message": "Unable to access the source exam file"})
	}
	src, err := examfile.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	//copy and save the exam file
	filepath := a.DataDir + "/learners/" + examfile.Filename
	dst, err := os.Create(filepath)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"Status": "Error", "Message": "Unable to create the destination exam file"})
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"Status": "Error", "Message": "Unable to write the exam file"})
	}

	//close off the exam if need be
	final := c.FormValue("final")
	if final == "closed" {
		a.DB.CloseLearnerExam(studentid, examid, false)
	}
	return c.JSON(http.StatusOK, map[string]any{"Status": "OK"})
}

// GET /auth/{examid}/{studentid}
// HandleGetStudentAuth checks if the learner is permitted to engage in the chosen exam identified by examid
// returns the exam password if correct
// TODO: include an auth token to authorise an exam
func (a *App) HandleGetStudentAuth(c echo.Context) error {
	// Check if request if a GET request
	if c.Request().Method != http.MethodGet {
		return c.JSON(http.StatusMethodNotAllowed, map[string]interface{}{
			"error": "Method not allowed",
		})
	}
	examid := c.Param("examid")
	studentid := c.Param("studentid")

	//learner must exist and be an active learner in the system
	if a.DB.IsLearnerValid(studentid) == false {
		return c.JSON(http.StatusBadRequest, map[string]any{"Status": "Error", "Message": "Invalid or inactive student ID"})
	}

	// chgeck if the exam is still open. A closed/expired exam cannot be authorised
	if a.DB.CheckExamClosed(examid, studentid) {
		//make sure the exam is recorded as expired and not closed then return
		//if a.DB.CloseLearnerExam(studentid, examid, true) != nil {
		//	return c.JSON(http.StatusBadRequest, map[string]any{"Status": "Error", "Message": "Exam has expired, status set to expire"})
		//}
		return c.JSON(http.StatusBadRequest, map[string]any{"Status": "Error", "Message": "Exam has expired or been closed"})
	}

	//check if the learner is allocated to the exam
	password, err := a.DB.GetExamPassword(examid, studentid)
	if password == "" || err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"Status": "Error", "Message": "Unable to get an authorised password"})
	}
	//set the exam active and start time once the learner has bene authorised
	err = a.DB.StartLearnerExam(studentid, examid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"Status": "Error", "Message": "Unable to initiate the exam"})
	}

	return c.JSON(http.StatusOK, map[string]any{"Status": "OK", "examid": examid, "studentid": studentid, "password": password})
}

// GET /exam/{examid}/{password}
// HandleGetStudentExam retrieves an exam based on the ID and password
// TODO: use the learners details and an auth token to enforce the security
func (a *App) HandleGetStudentExam(c echo.Context) error {
	// Check if request if a GET request
	if c.Request().Method != http.MethodGet {
		return c.JSON(http.StatusMethodNotAllowed, map[string]interface{}{
			"error": "Method not allowed",
		})
	}
	examid := c.Param("examid")
	password := c.Param("password")
	isvalid := a.DB.IsExamActive(examid, password)
	if isvalid == false {
		return c.JSON(http.StatusBadRequest, map[string]any{"success": false, "Message": "Exam retrieval unauthorised"})
	}
	//read the entire exam file into memory - around 50KB of text
	filepath := a.DataDir + "/exams/" + strings.Replace(examid, ".", "_", 1) + ".json"
	data, err := os.ReadFile(filepath)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"success": false, "Message": "Unable to retrieve the exam file"})
	}
	//send the file back
	return c.JSONBlob(http.StatusOK, data)
}
