package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBtype        string
	DBUser        string
	DBPassword    string
	DBName        string //also used as the name of the SQLite database file
	DBHost        string
	DBPort        int
	AdminPassword string
	JWTSecret     string
	AdminEmail    string
	DataDir       string
	ADSPORT       string
}

func LoadConfig() Config {
	// Try to load .env file, but don't fail if it doesn't exist
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, will attempt to use environment variables")
	}

	// List of required environment variables
	requiredEnvVars := map[string]string{
		"DB_TYPE":        "",
		"DB_USER":        "",
		"DB_PASSWORD":    "",
		"DB_NAME":        "",
		"DB_HOST":        "",
		"DB_PORT":        "",
		"ADMIN_PASSWORD": "",
		"JWT_SECRET":     "",
		"ADMIN_EMAIL":    "",
		"DATA_DIR":       "",
		"ADSPORT":        "",
	}

	// Check for missing environment variables
	var missingVars []string
	for env := range requiredEnvVars {
		if value := os.Getenv(env); value == "" {
			missingVars = append(missingVars, env)
		}
	}

	// If any required variables are missing, log them and exit
	if len(missingVars) > 0 {
		log.Fatalf("Missing environment variables: %v", missingVars)
	}

	// Get and validate DB_PORT
	dbPortStr := os.Getenv("DB_PORT")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("Invalid DB_PORT value: %v", err)
	}

	// Create and return the config
	return Config{
		DBtype:        os.Getenv("DB_TYPE"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        dbPort,
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
		AdminEmail:    os.Getenv("ADMIN_EMAIL"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		DataDir:       os.Getenv("DATA_DIR"),
		ADSPORT:       os.Getenv("ADSPORT"),
	}
}
