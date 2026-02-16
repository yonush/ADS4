package models

import (
	"database/sql"
)

/*-- Table to store exam learner information, with a unique StudentID for each learner (NON-SYSTEM USER)
CREATE TABLE "Learners" (
    "StudentID"   VARCHAR NOT NULL,
    "Name"        VARCHAR NOT NULL,
    "Status"      VARCHAR NOT NULL DEFAULT 'active',
    "Grade"       INTEGER,
    UNIQUE("StudentID"),
    PRIMARY KEY("StudentID"),
    CHECK (Status IN ('active','inactive'))
);
*/

type Learner struct {
	StudentID   sql.NullString `json:"studentid"`
	StudentName sql.NullString `json:"studentname"`
	Status      sql.NullString `json:"status"` // active, inactive
}

type LearnerDto struct {
	StudentID   string `json:"studentid"`
	StudentName string `json:"studentname"`
	Status      string `json:"status"` // active, inactive
}

// structure for reading CSV files - used with the seeding function - database/seed.go
// this will also be used with the data import feature
type LearnerCSV struct {
	StudentID   string `csv:"StudentID"`
	StudentName string `csv:"StudentName"`
	Status      string `csv:"Status"` // active, inactive
}
