package database

import (
	_ "database/sql"
	"errors"
	"strconv"
	"time"
)

/*
	Examination queries for all the non CRUD related queries
	used by:
	- dashboard - GetExamYears, GetExamByYearSemester
	- reporting - GetExaminations
	- AMT - GetExaminations
*/

// used to hold a list of years for the offerings in the DB
type Year struct {
	Year string `json:"Year"`
}

// GetExamYears retrieves all exam offering years from the database. This is used by the admin interface filters
func (db *DB) GetExamYears() ([]Year, error) {
	var query string
	var args []any

	query = `SELECT distinct year FROM examMetrics ORDER BY year ASC;`

	// Prepare and execute the query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Define the result slice
	var years []Year

	// Scan the results
	for rows.Next() {
		var year Year
		err := rows.Scan(
			&year.Year,
		)

		if err != nil {
			return nil, err
		}

		years = append(years, year)
	}

	// Return empty slice if:
	// 1. no offerings are found
	if len(years) == 0 {
		return []Year{}, nil
	}

	return years, nil
}

type ExamMetrics struct {
	CourseCode  string `json:"coursecode"`
	Description string `json:"description"`
	Password    string `json:"password"`
	ExamID      string `json:"examid"` //[year:4][semester:2][coursecode:*]
	Semester    string `json:"semester"`
	Year        string `json:"year"`
	Ready       string `json:"ready"`
	Active      string `json:"active"`
	Expired     string `json:"expired"`
	Closed      string `json:"closed"`
}

// query the exam offerings and metrics filtered by the offering year and semester

func (db *DB) GetExamByYearSemester(year, semester string) ([]ExamMetrics, error) {
	var query string

	query = `SELECT CourseCode,Description, Password, ExamID, Year, Semester,
			 Ready, Active, Expired, Closed
			 FROM examMetrics 
			 WHERE Year=$1 AND Semester = $2
			 ORDER BY coursecode DESC`

	// Prepare and execute the query
	rows, err := db.Query(query, year, semester)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Define the result slice
	var exammetrics []ExamMetrics

	// Scan the results
	for rows.Next() {
		var exammetric ExamMetrics
		err := rows.Scan(
			&exammetric.CourseCode,
			&exammetric.Description,
			&exammetric.Password,
			&exammetric.ExamID,
			&exammetric.Year,
			&exammetric.Semester,

			&exammetric.Ready,
			&exammetric.Active,
			&exammetric.Expired,
			&exammetric.Closed,
		)

		if err != nil {
			return nil, err
		}

		exammetrics = append(exammetrics, exammetric)
	}

	// Return empty slice if:
	// 1. no offerings are found
	if len(exammetrics) == 0 {
		return []ExamMetrics{}, nil
	}

	return exammetrics, nil
}

type ExamDetails struct {
	CourseCode  string `json:"coursecode"`
	StudentID   string `json:"studentid"`
	LearnerName string `json:"learnername"`
	ExamID      string `json:"examid"`
	Grade       string `json:"grade"`
}

func (db *DB) GetExaminations(field, value, semester string) ([]ExamDetails, error) {
	var query string
	var args []interface{}
	ErrNotFound := errors.New("query field not found")

	query = `SELECT CourseCode, StudentID, Name, ExamID, Grade
			 FROM ClosedExams `

	//default to the current year and add the semester
	year := strconv.Itoa(time.Now().Year())
	year = "2025"
	args = append(args, year)
	args = append(args, semester)
	query += `WHERE year=$1 AND semester=$2 AND `

	switch field {
	case "student":
		query += `studentid=$3 ORDER BY studentid`
		args = append(args, value)
	case "course":
		query += `coursecode=$3 ORDER BY coursecode`
		args = append(args, value)
	case "examid":
		query += `examid=$3 ORDER BY examid`
		args = append(args, value)
	default:
		return nil, ErrNotFound
	}

	// Prepare and execute the query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Define the result slice
	var exammetrics []ExamDetails

	// Scan the results
	for rows.Next() {
		var exammetric ExamDetails
		err := rows.Scan(
			&exammetric.CourseCode,
			&exammetric.StudentID,
			&exammetric.LearnerName,
			&exammetric.ExamID,
			&exammetric.Grade,
		)

		if err != nil {
			return nil, err
		}

		exammetrics = append(exammetrics, exammetric)
	}
	//fmt.Print(exammetrics)
	// Return empty slice if:
	// 1. no offerings are found
	if len(exammetrics) == 0 {
		return []ExamDetails{}, nil
	}

	return exammetrics, nil
}
