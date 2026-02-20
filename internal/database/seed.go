package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	"ADS4/internal/config"

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

func SeedDatabase(db *DB) error {
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

func SeedData(db *DB) {
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

	// purge and import new data
	log.Println("Importing course & offerings data.")
	log.Printf("- Reading %v/courses.csv", datadir)
	err = db.ImportCourses(datadir+"/courses.csv", true, false)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("- Reading %v/learners.csv", datadir)
	err = db.ImportLearners(datadir+"/learners.csv", true, false)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("- Reading %v/offerings.csv", datadir)
	err = db.ImportOfferings(datadir+"/offerings.csv", true, false)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("- Reading %v/learnerexams.csv", datadir)
	err = db.ImportLearnerExams(datadir+"/learnerexams.csv", true, false)
	if err != nil {
		log.Fatal(err)
	}

	// Create a temp file in data/ directory
	tempFile, err := os.Create(datadir + "/seed_complete")
	if err != nil {
		log.Fatal(err)
	}
	tempFile.Close()

	// Set environment variable to indicate that data has been seeded
	os.Setenv("DATA_SEEDED", "true")

	log.Println("Seeding complete.")
}
