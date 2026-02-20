package database

import (
	"ADS4/internal/models"
	_ "database/sql"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
)

/*
	purge - remove exisitng dataImportCourses
	overwrite - update exisitng data not insert

	purge overrides the overwrite flag - cannot update missing data
*/
// Import order 1
func (db *DB) ImportCourses(datafile string, purge, overwrite bool) error {
	src, err := os.Open(datafile)
	if err != nil {
		return err
	}
	defer src.Close()

	//read the CSV file in
	//TODO include some data checks before attempting to load
	var courses []*models.CoursesCSV
	if err := gocsv.UnmarshalFile(src, &courses); err != nil {
		return err
	}
	
	if purge {
		overwrite = false
		_, err := db.Exec(`DELETE FROM courses;`)
		if err != nil {
			return err
		}
	}

	var query string
	if overwrite {
		//perform updates
		query = `INSERT INTO Courses (CourseCode, Description, Level, Status)
			 	 VALUES ($1, $2, $3, $4)`
	} else {
		query = `INSERT INTO Courses (CourseCode, Description, Level, Status)
				 VALUES ($1, $2, $3, $4)`
		//assume the data is correct then perform the inserts
	}

	for _, course := range courses {
		_, err := db.Exec(query, course.CourseCode, course.Description, course.Level, course.Status)
		if err != nil {
			return err
		}
	}

	return nil
}

// Import order 1
func (db *DB) ImportLearners(datafile string, purge, overwrite bool) error {
	
	src, err := os.Open(datafile)
	if err != nil {
		return err
	}
	defer src.Close()

	//read the CSV file in
	//TODO include some data checks before attempting to load
	var learners []*models.LearnerCSV
	if err := gocsv.UnmarshalFile(src, &learners); err != nil {
		return err
	}
	
	if purge {
		overwrite = false
		_, err := db.Exec(`DELETE FROM learners;`)
		if err != nil {
			return err
		}
	}

	var query string
	if overwrite {
		//perform updates
		query = `INSERT INTO Learners (StudentID, Name, Status)
				VALUES ($1, $2, $3)`
	} else {
		query = `INSERT INTO Learners (StudentID, Name, Status)
				VALUES ($1, $2, $3)`
		//assume the data is correct then perform the inserts
	}

	for _, learner := range learners {
		_, err := db.Exec(query, learner.StudentID, learner.StudentName, learner.Status)
		if err != nil {
			return err
		}
	}

	return nil
}

// Import order 2


func (db *DB) ImportOfferings(datafile string, purge, overwrite bool) error {
	
	src, err := os.Open(datafile)
	if err != nil {
		return err
	}
	defer src.Close()


	//read the CSV file in
	//TODO include some data checks before attempting to load
	var offerings []*models.OfferingsCSV
	if err := gocsv.UnmarshalFile(src, &offerings); err != nil {
		return err
	}

	if purge {
		overwrite = false
		_, err := db.Exec(`DELETE FROM offerings;`)
		if err != nil {
			return err
		}
	}

	var query string
	if overwrite {
		//perform updates
		query = `INSERT INTO Offerings (ExamID, Year, Semester, CourseCode, Password, Status, Coordinator, OwnerID, Duration)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	} else {
		query = `INSERT INTO Offerings (ExamID, Year, Semester, CourseCode, Password, Status, Coordinator, OwnerID, Duration)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
		//assume the data is correct then perform the inserts
	}

	for _, offering := range offerings {
		_, err := db.Exec(query, offering.ExamID, offering.Year, offering.Semester, strings.TrimSpace(offering.CourseCode),
			offering.Password, offering.Status, offering.Coordinator, offering.OwnerID, offering.Duration)
		if err != nil {
			return err
		}
	}

	return nil
}

// Import order 3

func (db *DB) ImportLearnerExams(datafile string, purge, overwrite bool) error {
	
	src, err := os.Open(datafile)
	if err != nil {
		return err
	}
	defer src.Close()

	//read the CSV file in
	//TODO include some data checks before attempting to load
	var learnerexams []*models.LearnerExamCSV
	if err := gocsv.UnmarshalFile(src, &learnerexams); err != nil {
		return err
	}
	
	if purge {
		overwrite = false
		_, err := db.Exec(`DELETE FROM learnerexams;`)
		if err != nil {
			return err
		}
	}

	var query string
	if overwrite {
		//perform updates
		query = `INSERT INTO learnerexams (StudentID, ExamID, StartTime, EndTime,Status, Grade)
				VALUES ($1, $2, $3, $4, $5, $6)`
	} else {
		query = `INSERT INTO learnerexams (StudentID, ExamID, StartTime, EndTime,Status, Grade)
				VALUES ($1, $2, $3, $4, $5, $6)`
		//assume the data is correct then perform the inserts
	}

	for _, learnerexam := range learnerexams {
		_, err := db.Exec(query, learnerexam.StudentID, learnerexam.ExamID, learnerexam.StartTime, learnerexam.EndTime,
			learnerexam.Status, learnerexam.Grade)
		if err != nil {
			return err
		}
	}

	return nil
}
