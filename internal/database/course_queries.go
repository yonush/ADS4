package database

import (
	"ADS4/internal/models"
	"database/sql"
	_ "database/sql"
)

/*
-- Table to store courses created by admins, with a unique course code
CREATE TABLE IF NOT EXISTS "Courses" (
	"CourseCode"  VARCHAR(9),
	"Description" VARCHAR(255),
	"Level"  	  int, -- 1-9
    "Status"      VARCHAR(6) NOT NULL DEFAULT 'active', -- active, closed
	PRIMARY KEY("CourseCode"),
	UNIQUE("CourseCode"),
	CHECK (Status IN ('active','closed'))
	CHECK (Level IN (1,2,33,4,5,6,7,8,9)))
);
*/
// GetAllCourses retrieves all exam Courses from the database, with optional filtering by exam ID or status code. If no Courses are found, it returns an empty slice.
func (db *DB) GetAllCourses(coursecode string, statusCode string) ([]models.Courses, error) {
	var query string
	var args []interface{}
	var CourseExists bool

	// Check if the Course exists (if exam ID is provided)
	if coursecode != "" {
		CourseExistsQuery := `
		SELECT EXISTS (SELECT 1 FROM Courses o WHERE o.CourseCode = $1)`
		err := db.QueryRow(CourseExistsQuery, coursecode).Scan(&CourseExists)
		if err != nil && err == sql.ErrNoRows {
			return nil, err
		}
	}
	query = `SELECT o.CourseCode, o.Description, o.Level, o.Status
 		     FROM Courses o
	`

	// Add filtering by examid OR statuscode not both
	if statusCode != "" {
		query += `WHERE o.Status = $1`
		args = append(args, statusCode)
	} else if coursecode != "" {
		query += `WHERE o.CourseCode = $1`
		args = append(args, coursecode)

	}
	// Prepare and execute the query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Define the result slice
	var Courses []models.Courses

	// Scan the results
	for rows.Next() {
		var Course models.Courses
		err := rows.Scan(
			&Course.CourseCode,
			&Course.Description,
			&Course.Level,
			&Course.Status,
		)

		if err != nil {
			return nil, err
		}

		Courses = append(Courses, Course)
	}

	// Return empty slice if:
	// 1. no Courses are found
	if len(Courses) == 0 {
		return []models.Courses{}, nil
	}

	return Courses, nil
}

// GetCourseByID retrieves a specific exam Course from the database based on the provided exam ID. If the Course is not found, it returns nil.
func (db *DB) GetCourseByID(coursecode string) (*models.Courses, error) {

	var query string
	var CourseExists bool

	// Check if the Course exists (if Course ID is provided)
	if coursecode != "" {
		CourseExistsQuery := `
		SELECT EXISTS (SELECT 1 FROM Courses o WHERE o.CourseCode = $1)`
		err := db.QueryRow(CourseExistsQuery, coursecode).Scan(&CourseExists)
		if err != nil && err == sql.ErrNoRows {
			return nil, err
		}
	}

	query = `SELECT o.CourseCode, o.description, o.Level, o.Status
 		     FROM Courses o WHERE o.CourseCode = $1'`

	var Course models.Courses
	err := db.QueryRow(query, coursecode).Scan(
		&Course.CourseCode,
		&Course.Description,
		&Course.Level,
		&Course.Status,
	)

	if err != nil {
		return nil, err
	}

	return &Course, nil
}

// AddCourse adds a new exam Course to the database. It takes an Courses struct as input and returns an error if the operation fails.
func (db *DB) AddExamCourse(Course *models.Courses) error {
	query := "INSERT INTO Courses (CourseCode, Description, Level, Status) VALUES ($1, $2, $3, $4)"
	insertStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer insertStmt.Close()

	_, err = insertStmt.Exec(
		Course.CourseCode,
		Course.Description,
		Course.Level,
		Course.Status,
	)

	if err != nil {
		return err
	}

	return nil
}

// UpdateCourse updates an existing exam Course in the database based on the provided Courses struct. It returns an error if the operation fails.
func (db *DB) UpdateCourse(Course *models.Courses) error {
	//dont update the coursecode as it is the primary key and should not be changed
	query := "UPDATE Courses SET Description=$2, Level=$3, Status=$4 WHERE CourseCode=$1"
	updateStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer updateStmt.Close()

	_, err = updateStmt.Exec(
		Course.CourseCode,
		Course.Description,
		Course.Level,
		Course.Status)

	if err != nil {
		return err
	}

	return nil
}

// UpdateCourseStatus updates the status of an existing exam Course in the database based on the provided exam ID and status code. It returns an error if the operation fails.
func (db *DB) UpdateCourseStatus(coursecode string, statusCode string) error {
	query := "UPDATE Courses SET Status=$1 WHERE CourseCode=$2"
	updateStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer updateStmt.Close()

	_, err = updateStmt.Exec(statusCode, coursecode)

	if err != nil {
		return err
	}

	return nil
}

// DeleteCourse deletes an existing exam Course from the database based on the provided exam ID. It returns an error if the operation fails.
func (db *DB) DeleteCourse(coursecode string) error {
	query := "DELETE FROM Courses WHERE CourseCode = $1"
	deleteStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer deleteStmt.Close()

	_, err = deleteStmt.Exec(coursecode)

	if err != nil {
		return err
	}

	return nil
}
