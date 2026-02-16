package database

import (
	"ADS4/internal/models"
	"database/sql"
	_ "database/sql"
)

/*
-- Table to store exam learner information, with a unique StudentID for each learner (NON-SYSTEM USER)
CREATE TABLE "Learners" (

	"StudentID"   VARCHAR(8) NOT NULL,
	"Name"        VARCHAR NOT NULL,
	"Status"      VARCHAR NOT NULL DEFAULT 'active',
	UNIQUE("StudentID"),
	PRIMARY KEY("StudentID"),
	CHECK (Status IN ('active','inactive'))

);
*/

func (db *DB) IsLearnerValid(studentid string) bool {
	var learnerExists bool

	if studentid == "" {
		return false
	}
	learnerExistsQuery := `
	SELECT EXISTS (SELECT 1 FROM Learners WHERE studentid = $1 AND status = 'active') as active`
	err := db.QueryRow(learnerExistsQuery, studentid).Scan(&learnerExists)
	if err != nil && err == sql.ErrNoRows {
		return false
	}

	return learnerExists
}

func (db *DB) GetAllLearners(studentid string, statusCode string) ([]models.Learner, error) {
	var query string
	//var args []interface{}
	var args []any
	var learnerExists bool

	// Check if the offering exists (if exam ID is provided)
	if studentid == "" {
		return []models.Learner{}, nil
	}

	learnerExistsQuery := `
		SELECT EXISTS (SELECT 1 FROM Learners WHERE studentid = $1)`
	err := db.QueryRow(learnerExistsQuery, studentid).Scan(&learnerExists)
	if err != nil && err == sql.ErrNoRows {
		return nil, err
	}
	query = `SELECT l.studentid, l.Name, l.Status FROM Learners l `

	// Add filtering by studentid OR statuscode not both
	if statusCode != "" {
		query += `WHERE o.status = $1`
		args = append(args, statusCode)
	} else if studentid != "" {
		query += `WHERE o.studentid = $1`
		args = append(args, studentid)

	}
	// Prepare and execute the query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Define the result slice
	var learners []models.Learner

	// Scan the results
	for rows.Next() {
		var learner models.Learner
		err := rows.Scan(
			&learner.StudentID,
			&learner.StudentName,
			&learner.Status,
		)

		if err != nil {
			return nil, err
		}

		learners = append(learners, learner)
	}

	// Return empty slice if:
	// 1. no Learners are found
	if len(learners) == 0 {
		return []models.Learner{}, nil
	}

	return learners, nil
}

func (db *DB) GetLearnerByID(studentid string) (*models.Learner, error) {
	var query string
	var learnerExists bool

	// Check if the offering exists (if offering ID is provided)
	if studentid != "" {
		learnerExistsQuery := `
		SELECT EXISTS (SELECT 1 FROM Learners WHERE studentid = $1)`
		err := db.QueryRow(learnerExistsQuery, studentid).Scan(&learnerExists)
		if err != nil && err == sql.ErrNoRows {
			return nil, err
		}
	}

	query = `SELECT l.studentid, l.Name, l.Status FROM Learners l WHERE l.studentid = $1'`

	var learner models.Learner
	err := db.QueryRow(query, studentid).Scan(
		&learner.StudentID,
		&learner.StudentName,
		&learner.Status,
	)

	if err != nil {
		return nil, err
	}

	return &learner, nil
}

func (db *DB) AddLearner(learner *models.Learner) error {

	query := "INSERT INTO Learners (studentid, name, status) VALUES ($1, $2, $3)"
	insertStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer insertStmt.Close()

	_, err = insertStmt.Exec(
		learner.StudentID,
		learner.StudentName,
		learner.Status,
	)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateLearner(learner *models.Learner) error {
	//dont update the studentid as it is the primary key and should not be changed
	query := "UPDATE Learners SET name=$2, status=$3 WHERE studentid=$1"
	updateStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer updateStmt.Close()

	_, err = updateStmt.Exec(
		learner.StudentID,
		learner.StudentName,
		learner.Status,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateLearnerStatus(studentid string, statusCode string) error {

	query := "UPDATE Learners SET Status=$1 WHERE studentid=$2"
	updateStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer updateStmt.Close()

	_, err = updateStmt.Exec(statusCode, studentid)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteLearner(studentid string) error {
	query := "DELETE FROM Learners WHERE studentid = $1"
	deleteStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer deleteStmt.Close()

	_, err = deleteStmt.Exec(studentid)

	if err != nil {
		return err
	}

	return nil
}
