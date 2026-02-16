package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/joho/godotenv"

	"ADS4/internal/config"
	"ADS4/internal/models"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
	"golang.org/x/crypto/bcrypt"
)

// SeedStatus tracks the seeding state
type SeedStatus struct {
	TempFilePath string
	EnvVarName   string
}

// NewSeedStatus creates a new SeedStatus instance
func NewSeedStatus() *SeedStatus {
	datadir := config.LoadConfig().DataDir
	return &SeedStatus{
		TempFilePath: filepath.Join(datadir, "seed_complete"),
		EnvVarName:   "DATA_SEEDED",
	}
}

// IsDataSeeded checks if data has been seeded using both temp file and env var
func (s *SeedStatus) IsDataSeeded() bool {

	_, err := os.Stat(s.TempFilePath)
	fileExists := !os.IsNotExist(err)

	// Check environment variable
	envVal := os.Getenv(s.EnvVarName)
	envExists := envVal == "true"

	// Load from .env file if environment variable is not set
	if !envExists {
		if err := godotenv.Load(); err == nil {
			envVal = os.Getenv(s.EnvVarName)
			envExists = envVal == "true"
		}
	}

	// Return true if either check passes
	return fileExists || envExists
}

// MarkDataAsSeeded marks the data as seeded
func (s *SeedStatus) MarkDataAsSeeded() error {
	// Create temp file
	err := os.MkdirAll(filepath.Dir(s.TempFilePath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	tempFile, err := os.Create(s.TempFilePath)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	tempFile.Close()

	// Set environment variable
	err = os.Setenv(s.EnvVarName, "true")
	if err != nil {
		return fmt.Errorf("failed to set environment variable: %v", err)
	}

	return nil
}

func SeedDatabase(db *sql.DB) error {
	seedStatus := NewSeedStatus()

	if !seedStatus.IsDataSeeded() {
		log.Println("Seeding database...")
		SeedData(db) // Corrected here to remove error handling
		if err := seedStatus.MarkDataAsSeeded(); err != nil {
			return fmt.Errorf("failed to mark data as seeded: %v", err)
		}
		log.Println("Database seeding completed successfully")
	} else {
		log.Println("Database already seeded")
	}

	return nil
}

func SeedData(db *sql.DB) {
	// Get admin password from .env
	adminPassword := config.LoadConfig().AdminPassword
	datadir := config.LoadConfig().DataDir

	if adminPassword == "" {
		log.Fatal("ADMIN_PASSWORD not set in .env file")
	}

	userPassword := "Pa$$w0rd" //default password for seeded user

	// Generate hash for password
	adminHash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error generating hash for admin password")
		log.Fatal(err)
	}

	userHash, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)

	if err != nil {
		fmt.Println("Error generating hash for user password")
		log.Fatal(err)
	}

	log.Println("Seeding data...")

	// Insert Users - user names must be at least 6 characters
	_, err = db.Exec(`
		INSERT INTO UserT (username, password, role, email, defaultadmin)
		VALUES ('adminx', $1, 'Admin', 'admin@email.com', 1)`, adminHash)
	if err != nil {
		log.Printf("- adminx user exists, skipping admin user creation")
	}

	_, err = db.Exec(`
		INSERT INTO UserT (username, password, role, email, defaultadmin)
		VALUES ('bobbyx', $1, 'Faculty', 'bobbyx@email.com',0)`, userHash)
	if err != nil {
		log.Printf("- bobbyx user exists, skipping faculty user creation")
	}

	log.Println("Importing course & offerings data.")
	log.Printf("- Purging current demo data")
	_, err = db.Exec(`DELETE FROM learnerexams; DELETE FROM learners; DELETE FROM offerings; DELETE FROM courses;`)
	if err != nil {
		log.Fatal(err)
	}
	//important note: the CSV files must be in the same order as the database tables to avoid foreign key constraint errors when seeding
	log.Printf("- Reading %v/courses.csv", datadir)
	f, err := os.Open(datadir + "/courses.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var courses []*models.CoursesCSV
	if err = gocsv.UnmarshalFile(f, &courses); err != nil {
		log.Fatal(err)
	}

	for _, course := range courses {
		_, err = db.Exec(`
			INSERT INTO Courses (CourseCode, Description, Level, Status)
			VALUES ($1, $2, $3, $4)`, course.CourseCode, course.Description, course.Level, course.Status)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("- Reading %v/offerings.csv", datadir)
	f, err = os.Open(datadir + "/offerings.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var offerings []*models.OfferingsCSV
	if err = gocsv.UnmarshalFile(f, &offerings); err != nil {
		log.Fatal(err)
	}

	for _, offering := range offerings {
		_, err = db.Exec(`
				INSERT INTO Offerings (ExamID, Year, Semester, CourseCode, Password, Status, Coordinator, OwnerID, Duration)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			offering.ExamID, offering.Year, offering.Semester, strings.TrimSpace(offering.CourseCode),
			offering.Password, offering.Status, offering.Coordinator, offering.OwnerID, offering.Duration)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("- Reading %v/learners.csv", datadir)
	f, err = os.Open(datadir + "/learners.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var learners []*models.LearnerCSV
	if err = gocsv.UnmarshalFile(f, &learners); err != nil {
		log.Fatal(err)
	}

	for _, learner := range learners {
		_, err = db.Exec(`
				INSERT INTO Learners (StudentID, Name, Status)
				VALUES ($1, $2, $3)`, learner.StudentID, learner.StudentName, learner.Status)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("- Reading %v/learnerexams.csv", datadir)
	f, err = os.Open(datadir + "/learnerexam.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var learnerexams []*models.LearnerExamCSV
	if err = gocsv.UnmarshalFile(f, &learnerexams); err != nil {
		log.Fatal(err)
	}

	for _, learnerexam := range learnerexams {
		_, err = db.Exec(`
				INSERT INTO learnerexams (StudentID, ExamID, StartTime, EndTime,Status, Grade)
				VALUES ($1, $2, $3, $4, $5, $6)`,
			learnerexam.StudentID, learnerexam.ExamID, learnerexam.StartTime, learnerexam.EndTime,
			learnerexam.Status, learnerexam.Grade)
		if err != nil {
			log.Fatal(err)
		}
	}
	// Create a temp file in .internal/ directory

	tempFile, err := os.Create(datadir + "/seed_complete")
	if err != nil {
		log.Fatal(err)
	}
	tempFile.Close()

	// Set environment variable to indicate that data has been seeded
	os.Setenv("DATA_SEEDED", "true")

	log.Println("Seeding complete.")
}
