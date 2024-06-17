package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DBConn *gorm.DB

// InitDatabase initializes the database connection for a given database name
func InitDatabase(dbname string) (*gorm.DB, error) {
	var dbHost, dbUser, dbPassword, dbPort, dbSslmode string

	switch dbname {
	case "PolicyProcessorDB":
		dbHost = os.Getenv("POLICY_PROCESSOR_DB_HOST")
		dbUser = os.Getenv("POLICY_PROCESSOR_DB_USER")
		dbPassword = os.Getenv("POLICY_PROCESSOR_DB_PASSWORD")
		dbPort = os.Getenv("POLICY_PROCESSOR_DB_PORT")
		dbSslmode = os.Getenv("POLICY_PROCESSOR_DB_SSLMODE")
	case "RateDB":
		dbHost = os.Getenv("RATE_DB_HOST")
		dbUser = os.Getenv("RATE_DB_USER")
		dbPassword = os.Getenv("RATE_DB_PASSWORD")
		dbPort = os.Getenv("RATE_DB_PORT")
		dbSslmode = os.Getenv("RATE_DB_SSLMODE")
	default:
		return nil, fmt.Errorf("no database configuration found for database: %s", dbname)
	}

	if dbHost == "" || dbUser == "" || dbPassword == "" || dbPort == "" || dbSslmode == "" {
		return nil, fmt.Errorf("database configuration environment variables are not set")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbHost,
		dbUser,
		dbPassword,
		dbname, // use the dbname parameter directly here
		dbPort,
		dbSslmode,
	)

	// Create a folder for logs if it doesn't exist
	logsDir := "logs"
	err := os.MkdirAll(logsDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}

	// Generate log file name with current date and time
	logFileName := fmt.Sprintf("database_%s.log", dbname)
	logFilePath := filepath.Join(logsDir, logFileName)

	// Create a file for logs
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err) // Handle file creation error more gracefully if needed
	}

	// Create a logger that writes to the log file
	loggerFile := log.New(logFile, "", log.LstdFlags)

	loggerDb := logger.New(
		loggerFile,
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	DBConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: loggerDb})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database %s: %w", dbname, err)
	}

	return DBConn, nil
}

// GetDBConn returns the current database connection
func GetDBConn() *gorm.DB {
	return DBConn
}

// SetDBConn sets the database connection (useful for testing)
func SetDBConn(db *gorm.DB) {
	DBConn = db
}
