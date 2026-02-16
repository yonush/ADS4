package models

import "database/sql"

/* Courses
   - CourseCode e.g ITCS5.100
   - Level - 1-9
   - description
   - status - active,closed
*/
type Courses struct {
	CourseCode  sql.NullString `json:"coursecode"`
	Description sql.NullString `json:"description"`
	Level       int            `json:"level"`
	Status      sql.NullString `json:"status"` // active, closed
}

type CoursesDto struct {
	CourseCode  string `json:"coursecode"`
	Description string `json:"description"`
	Level       string `json:"level"`
	Status      string `json:"status"` // active, closed
}

//structure for reading CSV files - used with the seeding function - database/seed.go
//this will also be used with the data import feature
type CoursesCSV struct {
	CourseCode  string `csv:"CourseCode"`
	Description string `csv:"Description"`
	Level       int    `csv:"Level"`
	Status      string `csv:"Status"` // active, closed
}
