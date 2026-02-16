package database

import (
	"ADS4/internal/models"
	"database/sql"
	_ "database/sql"
)

/*
-- Table to store each Learnerexam's exam attempt, with a unique LearnerExamID for each attempt, and a foreign key reference to the Offerings table using ExamID

CREATE TABLE "Learnerexams" (

	"LearnerExamID" INTEGER,
	"StudentID"     VARCHAR(8) NOT NULL,
	"ExamID"        VARCHAR(15) NOT NULL,
	"StartTime"     TIME,
	"EndTime"       TIME,
	"Status"        VARCHAR(6) NOT NULL DEFAULT 'ready',
	"Grade"         INTEGER,
	PRIMARY KEY("LearnerExamID" AUTOINCREMENT),
	FOREIGN KEY("ExamID") REFERENCES "Offerings"("ExamID"),
	CHECK (Status IN ('ready', 'active', 'expire', 'closed', 'marked'))
	FOREIGN KEY("StudentID") REFERENCES "Learners"("StudentID")

);
*/

func (db *DB) StartLearnerExam(studentid, examid string) error {
	//use the SQL inbuilt time function - only sets time with no date or TZ info
	query := "UPDATE Learnerexams SET starttime=time('now','localtime'), Status='active' WHERE studentid=$1 AND examid=$2"
	updateStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer updateStmt.Close()

	_, err = updateStmt.Exec(studentid, examid)

	if err != nil {
		return err
	}

	return nil
}

// close the learner exam and set the status - expired or closed
func (db *DB) CloseLearnerExam(studentid, examid string, expired bool) error {
	var query string

	if expired {
		query = "UPDATE Learnerexams SET Status='expire',EndTime=time('now','localtime') WHERE studentid=$1 AND examid=$2"
	} else {
		query = "UPDATE Learnerexams SET Status='closed',EndTime=time('now','localtime') WHERE studentid=$1 AND examid=$2"
	}
	updateStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer updateStmt.Close()

	_, err = updateStmt.Exec(studentid, examid)

	if err != nil {
		return err
	}

	return nil
}

// check if the exam is still valid by checking the start and end time against the duration
// returns true
// this requires a join on the learnerexam and offerings tables
func (db *DB) CheckIfTime(examid, studentid string) bool {
	var HasTimeLeft bool

	// Check if the offering exists (if exam ID is provided)
	if studentid != "" && examid != "" {
		Query := `SELECT (substr(timediff(l.starttime,time('now','localtime')),13,2)*60+substr(timediff(l.starttime,time('now','localtime')),16,2)) < o.duration as isTimeLeft   
				  FROM Offerings o,  Learnerexams l
				  WHERE l.studentid=$1 AND l.examid=$2
	  				   AND l.status = 'active' AND o.examid = l.examid`
		err := db.QueryRow(Query, studentid, examid).Scan(&HasTimeLeft)
		if err != nil || err == sql.ErrNoRows {
			return false
		}
	}
	return HasTimeLeft
}

// check if the exam is still valid by checking the start and end time against the duration
// returns true
// this requires a join on the learnerexam and offerings tables
func (db *DB) CheckExamClosed(examid, studentid string) bool {
	var isValid bool

	// Check if the offering exists (if exam ID is provided)
	if studentid != "" && examid != "" {
		Query := `SELECT EXIST (SELECT 1   
				  FROM Offerings o,  Learnerexams l
				  WHERE l.studentid=$1 AND l.examid=$2
	  				   AND l.status IN ('active','ready') AND o.examid = l.examid) AS isValid`
		err := db.QueryRow(Query, studentid, examid).Scan(&isValid)
		if err != nil || err == sql.ErrNoRows {
			return false
		}
	}
	return isValid
}

// check if a learner is allowed to engage in an exam. Learner exam status must be ready or active
// returns the exam password if the learner is authorised
// this requires a join on the learnerexam and offerings tables
func (db *DB) GetExamPassword(examid, studentid string) (string, error) {
	var password string

	// Check if the offering exists (if exam ID is provided)
	if studentid != "" && examid != "" {
		Query := `SELECT o.password
				  FROM Offerings o,  Learnerexams l
				  WHERE l.studentid=$1
					AND l.examid=$2
					AND l.status IN ('ready','active') 
					AND o.examid = l.examid`
		err := db.QueryRow(Query, studentid, examid).Scan(&password)
		if err != nil || err == sql.ErrNoRows {
			return "", err
		}
	}

	return password, nil
}

// retrieves a learner exam. Learner exam status must be ready or active
// returns the exam if the learner is authorised

func (db *DB) IsExamActive(examid, password string) bool {
	var isExamAuth string

	// Check if the offering exists (if exam ID is provided)
	if password != "" && examid != "" {
		Query := `SELECT EXISTS (SELECT 1 FROM Offerings
				  WHERE examid=$1 AND password=$2 AND status = 'active')`
		err := db.QueryRow(Query, examid, password).Scan(&isExamAuth)
		if err != nil || err == sql.ErrNoRows {
			return false
		}
	}

	return true
}

func (db *DB) GetAllLearnerExams(learnerexamid int, statusCode string) ([]models.LearnerExam, error) {
	var query string
	//var args []interface{}
	var args []any
	var LearnerexamExists bool

	// Check if the offering exists (if exam ID is provided)
	if learnerexamid > 0 {
		LearnerexamExistsQuery := `
		SELECT EXISTS (SELECT 1 FROM Learnerexams WHERE learnerexamid = $1)`
		err := db.QueryRow(LearnerexamExistsQuery, learnerexamid).Scan(&LearnerexamExists)
		if err != nil {
			return nil, err
		}
	}
	query = `SELECT l.learnerexamid, l.studentid, l.examid, l.starttime, l.endtime, l.status, l.grade FROM Learnerexams l `

	// Add filtering by learnerexamid OR statuscode not both
	if statusCode != "" {
		query += `WHERE o.status = $1`
		args = append(args, statusCode)
	} else if learnerexamid > 0 {
		query += `WHERE o.learnerexamid = $1`
		args = append(args, learnerexamid)

	}
	// Prepare and execute the query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Define the result slice
	var Learnerexams []models.LearnerExam

	// Scan the results
	for rows.Next() {
		var Learnerexam models.LearnerExam
		err := rows.Scan(
			&Learnerexam.LearnerExamID,
			&Learnerexam.StudentID,
			&Learnerexam.ExamID,
			&Learnerexam.StartTime,
			&Learnerexam.EndTime,
			&Learnerexam.Grade,
			&Learnerexam.Status,
		)

		if err != nil {
			return nil, err
		}

		Learnerexams = append(Learnerexams, Learnerexam)
	}

	// Return empty slice if:
	// 1. no Learnerexams are found
	if len(Learnerexams) == 0 {
		return []models.LearnerExam{}, nil
	}

	return Learnerexams, nil
}

func (db *DB) GetLearnerExamByID(learnerexamid int) (*models.LearnerExam, error) {
	var query string
	var LearnerexamExists bool

	// Check if the offering exists (if offering ID is provided)
	if learnerexamid > 0 {
		LearnerexamExistsQuery := `
		SELECT EXISTS (SELECT 1 FROM Learnerexams WHERE learnerexamid = $1)`
		err := db.QueryRow(LearnerexamExistsQuery, learnerexamid).Scan(&LearnerexamExists)
		if err != nil {
			return nil, err
		}
	}

	query = `SELECT l.learnerexamid, l.studentid, l.examid, l.starttime, l.endtime, l.status, l.grade FROM Learnerexams l WHERE l.learnerexamid = $1`
	var Learnerexam models.LearnerExam
	err := db.QueryRow(query, learnerexamid).Scan(
		&Learnerexam.LearnerExamID,
		&Learnerexam.StudentID,
		&Learnerexam.ExamID,
		&Learnerexam.StartTime,
		&Learnerexam.EndTime,
		&Learnerexam.Grade,
		&Learnerexam.Status,
	)

	if err != nil {
		return nil, err
	}

	return &Learnerexam, nil
}

func (db *DB) AddLearnerExam(Learnerexam *models.LearnerExam) error {
	//query = `SELECT l.learnerexamid, l.studentid, l.examid, l.starttime, l.endtime, l.status, l.grade FROM Learnerexams l `
	query := `INSERT INTO Learnerexams (learnerexamid, studentid, examid, starttime, endtime, status, grade) 
			  VALUES ($1, $2, $3,$4, $5, $6, $7)`
	insertStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer insertStmt.Close()

	_, err = insertStmt.Exec(
		Learnerexam.LearnerExamID,
		Learnerexam.StudentID,
		Learnerexam.ExamID,
		Learnerexam.StartTime,
		Learnerexam.EndTime,
		Learnerexam.Grade,
		Learnerexam.Status,
	)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateLearnerExam(Learnerexam *models.LearnerExam) error {
	query := `UPDATE Learnerexams SET studentid=$2, examid=$3, starttime=$4, endtime=$5, status=$6, grade=$7 WHERE learnerexamid=$1`
	updateStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer updateStmt.Close()

	_, err = updateStmt.Exec(
		Learnerexam.LearnerExamID,
		Learnerexam.StudentID,
		Learnerexam.ExamID,
		Learnerexam.StartTime,
		Learnerexam.EndTime,
		Learnerexam.Grade,
		Learnerexam.Status,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateLearnerExamStatus(learnerexamid int, statusCode string) error {
	query := "UPDATE Learnerexams SET Status=$1 WHERE learnerexamid=$2"
	updateStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer updateStmt.Close()

	_, err = updateStmt.Exec(statusCode, learnerexamid)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateLearnerExamGrade(learnerexamid int, grade int) error {

	query := "UPDATE Learnerexams SET grade=$1 WHERE learnerexamid=$2"
	updateStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer updateStmt.Close()

	_, err = updateStmt.Exec(grade, learnerexamid)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteLearnerExam(learnerexamid int) error {
	query := "DELETE FROM Learnerexams WHERE learnerexamid = $1"
	deleteStmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer deleteStmt.Close()

	_, err = deleteStmt.Exec(learnerexamid)

	if err != nil {
		return err
	}

	return nil
}
