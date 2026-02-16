package database

import (
	"ADS4/internal/models"
	_ "database/sql"
)

/*
-- Table to store exam offerings created by admins, with a unique ExamID that follows the format [year:4][semester:2][coursecode:9]

CREATE TABLE "Offerings" (

	    "ExamID"      VARCHAR(15) NOT NULL,
	    "Year" 	      INTEGER,
	    "Semester"    VARCHAR(2) NOT NULL DEFAULT 'S1',
	    "CourseCode"  VARCHAR(9) NOT NULL,
	    "Password"    VARCHAR(8) NOT NULL,
	    "Status"      VARCHAR(6) NOT NULL DEFAULT 'active',
	    "Coordinator" INTEGER,
	    "OwnerID"     INTEGER,
	    "Duration"    INTEGER,
	    PRIMARY KEY("ExamID"),
	    UNIQUE("ExamID"),
	    FOREIGN KEY("Coordinator") REFERENCES "UserT"("UserID"),
	    FOREIGN KEY("OwnerID") REFERENCES "UserT"("UserID"),
	    FOREIGN KEY("CourseCode") REFERENCES "Courses"("CourseCode"),
		CHECK (Status IN ('active','closed')),
		CHECK ("Semester" IN ('S1','S2','S3')),
	    CHECK ("Duration" > 29 and "Duration" < 241)

);
*/
type Exam struct {
	ExamID      string `json:"ExamID"`
	CourseCode  string `json:"Code"`
	Description string `json:"Description"`
}

// used to hold a list of years for the offerings in the DB
type Year struct {
	Year string `json:"Year"`
}

// GetExamYears retrieves all exam offering years from the database. This is used by the admin interface filters
func (db *DB) GetExamYears() ([]Year, error) {
	var query string
	var args []any

	query = `SELECT distinct year FROM examMetrics;`

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

// GetAllOfferings retrieves all exam offerings from the database, with optional filtering by offering year. If no offerings are found, it returns an empty slice.
// TODO: add provision for semesters or alternately disable courses after a certain date
func (db *DB) GetActiveExams(filteryear string) ([]Exam, error) {
	var query string
	var args []any

	query = `SELECT o.examid AS ExamID, o.coursecode AS Code, c.description AS Description
 		     FROM Offerings o, Courses c 
			 WHERE o.coursecode = c.coursecode
			      AND o.status = 'active' AND c.status = 'active' 
				  `

	//e.g 2026
	// Add filtering by the current year or selected year
	if filteryear != "" {
		query += `AND o.year = $1`
		args = append(args, filteryear)
	}

	// Prepare and execute the query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Define the result slice
	var exams []Exam

	// Scan the results
	for rows.Next() {
		var exam Exam
		err := rows.Scan(
			&exam.ExamID,
			&exam.CourseCode,
			&exam.Description,
		)

		if err != nil {
			return nil, err
		}

		exams = append(exams, exam)
	}

	// Return empty slice if:
	// 1. no offerings are found
	if len(exams) == 0 {
		return []Exam{}, nil
	}

	return exams, nil
}

// GetAllOfferings retrieves all exam offerings from the database, with optional filtering by exam ID or status code. If no offerings are found, it returns an empty slice.
func (db *DB) GetAllOfferings(examID string, statusCode string) ([]models.Offerings, error) {
	var query string
	//var args []interface{}
	var args []any
	var offeringExists bool

	// Check if the offering exists (if exam ID is provided)
	if examID != "" {
		offeringExistsQuery := `
		SELECT EXISTS (SELECT 1 FROM Offerings o WHERE o.examid = $1)`
		err := db.QueryRow(offeringExistsQuery, examID).Scan(&offeringExists)
		if err != nil {
			return nil, err
		}
	}
	query = `SELECT o.examid, o.coursecode, o.year, o.semester, o.password, o.status, o.coordinator, o.ownerid,o.duration
 		     FROM Offerings o
	`

	// Add filtering by examid OR statuscode not both
	if statusCode != "" {
		query += `WHERE o.status = $1`
		args = append(args, statusCode)
	} else if examID != "" {
		query += `WHERE o.examid = $1`
		args = append(args, examID)

	}
	// Prepare and execute the query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Define the result slice
	var offerings []models.Offerings

	// Scan the results
	for rows.Next() {
		var offering models.Offerings
		err := rows.Scan(
			&offering.ExamID,
			&offering.Year,
			&offering.Semester,
			&offering.CourseCode,
			&offering.Password,
			&offering.Status,
			&offering.Coordinator,
			&offering.OwnerID,
			&offering.Duration,
		)

		if err != nil {
			return nil, err
		}

		offerings = append(offerings, offering)
	}

	// Return empty slice if:
	// 1. no offerings are found
	if len(offerings) == 0 {
		return []models.Offerings{}, nil
	}

	return offerings, nil
}

// GetOfferingByID retrieves a specific exam offering from the database based on the provided exam ID. If the offering is not found, it returns nil.
func (db *DB) GetOfferingByID(examID string) (*models.Offerings, error) {

	var query string
	var offeringExists bool

	// Check if the offering exists (if offering ID is provided)
	if examID != "" {
		offeringExistsQuery := `
		SELECT EXISTS (SELECT 1 FROM Offerings o WHERE o.examid = $1)`
		err := db.QueryRow(offeringExistsQuery, examID).Scan(&offeringExists)
		if err != nil {
			return nil, err
		}
	}

	query = `SELECT o.examid, o.coursecode, o.year, o.semester, o.password, o.status, o.coordinator, o.ownerid,o.duration
 		     FROM Offerings o WHERE o.examid = $1'`

	var offering models.Offerings
	err := db.QueryRow(query, examID).Scan(
		&offering.ExamID,
		&offering.Year,
		&offering.Semester,
		&offering.CourseCode,
		&offering.Password,
		&offering.Status,
		&offering.Coordinator,
		&offering.OwnerID,
		&offering.Duration,
	)

	if err != nil {
		return nil, err
	}

	return &offering, nil
}

// AddOffering adds a new exam offering to the database. It takes an Offerings struct as input and returns an error if the operation fails.
func (db *DB) AddExamOffering(offering *models.Offerings) error {
	query := `INSERT INTO Offerings (examID, coursecode, year, semester,  Password, Status, Coordinator, OwnerID, Duration ) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	insertStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer insertStmt.Close()

	_, err = insertStmt.Exec(
		offering.ExamID,
		offering.CourseCode,
		offering.Year,
		offering.Semester,
		offering.Password,
		offering.Status,
		offering.Coordinator,
		offering.OwnerID,
		offering.Duration,
	)

	if err != nil {
		return err
	}

	return nil
}

// UpdateOffering updates an existing exam offering in the database based on the provided Offerings struct. It returns an error if the operation fails.
func (db *DB) UpdateOffering(offering *models.Offerings) error {
	//dont update the examid as it is the primary key and should not be changed
	query := "UPDATE Offerings SET CourseCode=$2, Year=$3, Semester=$4, Password=$5, Status=$6, Coordinator=$7, OwnerID=$8 WHERE examID=$1"
	updateStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer updateStmt.Close()

	_, err = updateStmt.Exec(
		offering.ExamID,
		offering.CourseCode,
		offering.Year,
		offering.Semester,
		offering.Password,
		offering.Status,
		offering.Coordinator,
		offering.OwnerID,
		offering.Duration,
	)

	if err != nil {
		return err
	}

	return nil
}

// UpdateOfferingStatus updates the status of an existing exam offering in the database based on the provided exam ID and status code. It returns an error if the operation fails.
func (db *DB) UpdateOfferingStatus(examid string, statusCode string) error {
	query := "UPDATE Offerings SET Status=$1 WHERE examID=$2"
	updateStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer updateStmt.Close()

	_, err = updateStmt.Exec(statusCode, examid)

	if err != nil {
		return err
	}

	return nil
}

// DeleteOffering deletes an existing exam offering from the database based on the provided exam ID. It returns an error if the operation fails.
func (db *DB) DeleteOffering(examid string) error {
	query := "DELETE FROM Offerings WHERE examID = $1"
	deleteStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer deleteStmt.Close()

	_, err = deleteStmt.Exec(examid)

	if err != nil {
		return err
	}

	return nil
}
