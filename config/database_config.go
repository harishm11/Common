package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DBConn *gorm.DB

// InitDatabase initializes the database connection for a given database name
func InitDatabase(dbname string) (*gorm.DB, error) {

	env := viper.GetString("env")
	dbConfig := viper.Sub(env + "." + dbname)

	if dbConfig == nil {
		return nil, fmt.Errorf("no database configuration found for database: %s in environment: %s", dbname, env)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		dbConfig.GetString("host"),
		dbConfig.GetString("user"),
		dbConfig.GetString("password"),
		dbname, // use the dbname parameter directly here
		dbConfig.GetInt("port"),
		dbConfig.GetString("sslmode"),
	)

	// Create a folder for logs if it doesn't exist
	logsDir := "logs"
	err := os.MkdirAll(logsDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}

	// Generate log file name with current date and time
	//currentTime := time.Now().Format("2006-01-02_15-04-05")
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

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: loggerDb})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database %s: %w", dbname, err)
	}

	return db, nil
}
