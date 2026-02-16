package models

import (
	"database/sql"
)

/*
- RBAC - three roles - Admin, Faculty, Learner
- examID = [year:4][semester:2][coursecode:9]

**LearnerExam** - Learners that are elgible for exams
    - studentID
    - examID e.g 2026S1ITCS5.100
    - status - ready, active, expired, closed, marked
*/
/*-- Table to store each learner's exam attempt, with a unique LearnerExamID for each attempt, and a foreign key reference to the Offerings table using ExamID
CREATE TABLE "Learnerexams" (
    "LearnerExamID" INTEGER,
    "StudentID"     VARCHAR NOT NULL,
    "ExamID"        VARCHAR(15) NOT NULL,
    "StartTime"     TIME,
    "EndTime"       TIME,
    "Status"        VARCHAR(6) NOT NULL DEFAULT 'ready',
    PRIMARY KEY("LearnerExamID" AUTOINCREMENT),
    FOREIGN KEY("ExamID") REFERENCES "Offerings"("ExamID"),
    CHECK (Status IN ('ready', 'active', 'expire', 'closed', 'marked'))
    --FOREIGN KEY("StudentID") REFERENCES "Learners"("StudentID")
);

*/

type LearnerExam struct {
	LearnerExamID int            `json:"learnerexamid"`
	StudentID     sql.NullString `json:"studentid"`
	ExamID        sql.NullString `json:"examid"` // [year:4][semester:2][coursecode:9]
	StartTime     sql.NullTime   `json:"starttime"`
	EndTime       sql.NullTime   `json:"endtime"`
	Status        sql.NullString `json:"status"` // ready, active, expired, closed, marked
	Grade         sql.NullInt32  `json:"grade"`
}

type LearnerExamDto struct {
	LearnerExamID string `json:"learnerexamid"`
	StudentID     string `json:"studentid"`
	ExamID        string `json:"examid"` //[year:4][semester:2][coursecode:*]
	StartTime     string `json:"starttime"`
	EndTime       string `json:"endtime"`
	Status        string `json:"status"` // ready, active, expired, closed, marked
	Grade         string `json:"grade"`
}

// structure for reading CSV files - used with the seeding function - database/seed.go
// this will also be used with the data import feature
type LearnerExamCSV struct {
	StudentID string `csv:"StudentID"`
	ExamID    string `csv:"ExamID"` //[year:4][semester:2][coursecode:*]
	StartTime string `csv:"StartTime"`
	EndTime   string `csv:"EndTime"`
	Status    string `csv:"Status"` // ready, active, expired, closed, marked
	Grade     int    `csv:"Grade"`
}
