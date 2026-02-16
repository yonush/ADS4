package models

import "database/sql"

/* Examination Offerings
   - examID e.g 2026S1ITCS5.100
   - year e.g 2026
   - semester e.g S1,S2,S3
   - coursecode e.g ITCS5.100
   - password
   - status - active,closed
   - PC - program coodinator
   - owner - lecturer/s of the exam
*/

/* -- Table to store exam offerings created by admins, with a unique ExamID that follows the format [year:4][semester:2][coursecode:9]
CREATE TABLE "Offerings" (
    "ExamID"      VARCHAR(15) NOT NULL,
    "Year" 	      INTEGER,
    "Semester"    VARCHAR(2) NOT NULL DEFAULT 'S1',
    "CourseCode"  VARCHAR(9) NOT NULL,
    "Password"    VARCHAR(8) NOT NULL,
    "Status"      VARCHAR(6) NOT NULL DEFAULT 'active',
    "Coordinator" INTEGER,
    "OwnerID"     INTEGER,
    PRIMARY KEY("ExamID"),
    UNIQUE("ExamID"),
    FOREIGN KEY("Coordinator") REFERENCES "UserT"("UserID"),
    FOREIGN KEY("OwnerID") REFERENCES "UserT"("UserID"),
    FOREIGN KEY("CourseCode") REFERENCES "Courses"("CourseCode"),
	CHECK (Status IN ('active','closed'))
	CHECK (Status IN ('S1','S2','S3'))

);
*/

type Offerings struct {
	ExamID      sql.NullString `json:"examid"` // [year:4][semester:2][coursecode:*]
	CourseCode  sql.NullString `json:"coursecode"`
	Year        int            `json:"year"`
	Semester    sql.NullString `json:"semester"`
	Password    sql.NullString `json:"password"`
	Status      sql.NullString `json:"status"`
	Coordinator sql.NullString `json:"coordinator"`
	OwnerID     sql.NullString `json:"ownerid"`
	Duration    int            `json:"duration"`
}

type OfferingsDto struct {
	ExamID      string `json:"examid"` //[year:4][semester:2][coursecode:*]
	CourseCode  string `json:"coursecode"`
	Year        string `json:"year"`
	Semester    string `json:"semester"`
	Password    string `json:"password"`
	Status      string `json:"status"`
	Coordinator string `json:"coordinator"`
	OwnerID     string `json:"ownerid"`
	Duration    string `json:"duration"`
}

// structure for reading CSV files - used with the seeding function - database/seed.go
//this will also be used with the data import feature
type OfferingsCSV struct {
	ExamID      string `csv:"ExamID"`
	CourseCode  string `csv:"CourseCode"`
	Year        int    `csv:"Year"`
	Semester    string `csv:"Semester"`
	Password    string `csv:"Password"`
	Status      string `csv:"Status"`
	Coordinator int    `csv:"Coordinator"`
	OwnerID     int    `csv:"OwnerID"`
	Duration    int    `csv:"Duration"`
}
