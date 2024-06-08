package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

// Logger defines the interface for logging functionality.
type Logger interface {
	Info(msg string, kv ...interface{})
	Warn(msg string, kv ...interface{})
	Error(err error, msg string, kv ...interface{})
	Debug(msg string, kv ...interface{})
	LogAnything(data interface{})
	SetTransactionNumber(transactionNumber string)
}

// myLogger implements the Logger interface and wraps the slog.Logger
type myLogger struct {
	*slog.Logger
	transactionNumber string
}

func (l *myLogger) SetTransactionNumber(transactionNumber string) {
	l.transactionNumber = " "
	l.transactionNumber = transactionNumber
}

func (l *myLogger) Info(msg string, kv ...interface{}) {
	if l.transactionNumber != "" {
		kv = append(kv, "transactionNumber", l.transactionNumber)
	}
	l.Logger.Info(msg, kv...)
}

func (l *myLogger) Warn(msg string, kv ...interface{}) {
	if l.transactionNumber != "" {
		kv = append(kv, "transactionNumber", l.transactionNumber)
	}
	l.Logger.Warn(msg, kv...)
}

func (l *myLogger) Error(err error, msg string, kv ...interface{}) {
	// Get a clean string representation of the error
	errorStr := errToStr(err)

	// Construct the log message with additional key-value pairs
	logMessage := fmt.Sprintf("%s: %s", msg, errorStr)
	if l.transactionNumber != "" {
		kv = append(kv, "transactionNumber", l.transactionNumber)
	}
	l.Logger.Error(logMessage, kv...)
}

func (l *myLogger) Debug(msg string, kv ...interface{}) {
	// Check if debug logging is enabled before logging
	if l.transactionNumber != "" {
		kv = append(kv, "transactionNumber", l.transactionNumber)
	}
	l.Logger.Debug(msg, kv...)
}

func (l *myLogger) LogAnything(data interface{}) {
	// Convert the data to a string representation
	logData := convertDataToString(data)
	if l.transactionNumber != "" {
		logData = logData + " transactionNumber " + l.transactionNumber
	}
	// Log the data using the underlying slog.Logger
	l.Logger.Info(logData)
}

// Helper function to convert any data to a string
func convertDataToString(data interface{}) string {
	// Implement your logic here to convert any type of data to a string
	return fmt.Sprintf("%v", data)
}

// Helper function to convert an error to a clean string
func errToStr(err error) string {
	if err == nil {
		return "(no error)"
	}
	return err.Error()
}

var globalLogger *myLogger // Global logger variable

func InitLogger(serviceName string) {
	// Create a folder for logs if it doesn't exist
	logsDir := "logs"

	err := os.MkdirAll(logsDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}

	// Generate log file name with current date and service name
	logFileName := fmt.Sprintf("%s.log", serviceName)
	logFilePath := filepath.Join(logsDir, logFileName)

	// Create a file for logs
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err) // Handle file creation error more gracefully if needed
	}

	// Create a JSON handler that writes to the log file
	// handler := slog.NewJSONHandler(logFile, nil)
	handler := slog.NewTextHandler(logFile, nil)

	// Create a new logger with the JSON handler
	logger := slog.New(handler)

	// Wrap the logger in the myLogger struct (optional, for type safety)
	globalLogger = &myLogger{Logger: logger}
}

// GetLogger returns the global logger instance
func GetLogger() Logger {
	return globalLogger
}

func FrontendLogger(c *fiber.Ctx) error {

	type LogMessage struct {
		Message string `json:"message"`
		Level   string `json:"level"`
	}
	var message LogMessage
	err := c.BodyParser(&message)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid log request body from front end")
	}

	if message.Level == "info" {
		GetLogger().Info(message.Message)
	} else {
		GetLogger().LogAnything(message.Message)
	}

	return c.SendString("Log received")
}
