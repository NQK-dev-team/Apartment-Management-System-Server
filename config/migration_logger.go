package config

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// CustomMigrationLogger implements gorm.logger.Interface
type CustomMigrationLogger struct {
	logDirectory string
	logLevel     logger.LogLevel
	fileHandles  map[string]*os.File // Cache for open log files, key is the full filename
	mu           sync.Mutex          // Mutex to protect concurrent access to fileHandles map
}

// NewCustomMigrationLogger creates a new instance of our custom logger
func NewCustomMigrationLogger(dir string, level logger.LogLevel) *CustomMigrationLogger {
	// Ensure the log directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}
	return &CustomMigrationLogger{
		logDirectory: dir,
		logLevel:     level,
		fileHandles:  make(map[string]*os.File),
	}
}

// LogMode sets the log level
func (l *CustomMigrationLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

// Info logs general information
func (l *CustomMigrationLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Info {
		log.Printf("[INFO] "+msg, data...)
	}
}

// Warn logs warnings
func (l *CustomMigrationLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Warn {
		log.Printf("[WARN] "+msg, data...)
	}
}

// Error logs errors
func (l *CustomMigrationLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Error {
		log.Printf("[ERROR] "+msg, data...)
	}
}

func (l *CustomMigrationLogger) isCreateOrDropStatement(sql string) bool {
	return strings.HasPrefix(strings.TrimSpace(sql), "CREATE") || strings.HasPrefix(strings.TrimSpace(sql), "DROP")
}

// Trace is the core method for logging SQL queries
func (l *CustomMigrationLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, _ := fc()

	// Only select CREATE or DROP statements
	if !l.isCreateOrDropStatement(sql) {
		return
	}

	dateStr := time.Now().Format("2006-01-02") // YYYY-MM-DD format

	finalDir := filepath.Join(l.logDirectory)
	if err := os.MkdirAll(finalDir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	filename := filepath.Join(finalDir, fmt.Sprintf("%s.log", dateStr))

	// Use mutex to handle concurrent requests safely
	l.mu.Lock()
	defer l.mu.Unlock()

	// Get or create the file handle
	file, ok := l.fileHandles[filename]
	if !ok {
		var fileErr error
		// Open file in append mode, create if it doesn't exist
		file, fileErr = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if fileErr != nil {
			log.Printf("Failed to open log file %s: %v", filename, fileErr)
			return
		}
		l.fileHandles[filename] = file
	}

	// Create a dedicated logger for the file
	fileLogger := log.New(file, "", log.LstdFlags)

	// Format and write the log message
	switch {
	case err != nil && l.logLevel >= logger.Error && !errors.Is(err, gorm.ErrRecordNotFound):
		fileLogger.Printf("[%.3fms] %s\n[ERROR] %v", float64(elapsed.Nanoseconds())/1e6, sql, err)
	case elapsed > 200*time.Millisecond && l.logLevel >= logger.Warn: // Log slow queries
		fileLogger.Printf("[%.3fms] %s\n[SLOW QUERY]", float64(elapsed.Nanoseconds())/1e6, sql)
	case l.logLevel >= logger.Info:
		fileLogger.Printf("[%.3fms] %s", float64(elapsed.Nanoseconds())/1e6, sql)
	}
}
