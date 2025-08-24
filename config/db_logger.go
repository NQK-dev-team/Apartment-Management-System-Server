package config

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// CustomDBLogger implements gorm.logger.Interface
type CustomDBLogger struct {
	logDirectory string
	logLevel     logger.LogLevel
	fileHandles  map[string]*os.File // Cache for open log files, key is the full filename
	mu           sync.Mutex          // Mutex to protect concurrent access to fileHandles map
	selectRegex  *regexp.Regexp
	updateRegex  *regexp.Regexp
	deleteRegex  *regexp.Regexp
	intoRegex    *regexp.Regexp
}

// NewCustomDBLogger creates a new instance of our custom logger
func NewCustomDBLogger(dir string, level logger.LogLevel) *CustomDBLogger {
	// Ensure the log directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}
	return &CustomDBLogger{
		logDirectory: dir,
		logLevel:     level,
		fileHandles:  make(map[string]*os.File),
		selectRegex:  regexp.MustCompile(`(?i)SELECT\s+.*?\s+FROM\s+"?([a-zA-Z0-9_]+)"?`),
		updateRegex:  regexp.MustCompile(`(?i)UPDATE\s+"([^"]+)"`),
		deleteRegex:  regexp.MustCompile(`(?i)DELETE FROM\s+"([^"]+)"`),
		intoRegex:    regexp.MustCompile(`(?i)INTO\s+"([^"]+)"`),
	}
}

// extractTableName parses SQL to find the table name
func (l *CustomDBLogger) extractTableName(sql string) string {
	var match []string

	if GetEnv("APP_ENV") == "development" {
		if match = l.selectRegex.FindStringSubmatch(sql); len(match) > 1 {
			return match[1]
		}
	}
	if match = l.intoRegex.FindStringSubmatch(sql); len(match) > 1 {
		return match[1]
	}
	if match = l.updateRegex.FindStringSubmatch(sql); len(match) > 1 {
		return match[1]
	}
	if match = l.deleteRegex.FindStringSubmatch(sql); len(match) > 1 {
		return match[1]
	}
	return "gorm_default"
}

// LogMode sets the log level
func (l *CustomDBLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

// Info logs general information
func (l *CustomDBLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Info {
		log.Printf("[INFO] "+msg, data...)
	}
}

// Warn logs warnings
func (l *CustomDBLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Warn {
		log.Printf("[WARN] "+msg, data...)
	}
}

// Error logs errors
func (l *CustomDBLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Error {
		log.Printf("[ERROR] "+msg, data...)
	}
}

// Trace is the core method for logging SQL queries
func (l *CustomDBLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	tableName := l.extractTableName(sql)
	dateStr := time.Now().Format("2006-01-02") // YYYY-MM-DD format

	if tableName == "gorm_default" {
		return // Skip logging for default table
	}

	// Construct filename: logs/2025/2025-06/2025-06-25/users_2025-06-25.log
	currentDate := time.Now().Format("2006-01-02")
	currentYear := time.Now().Format("2006")
	currentMonth := time.Now().Format("2006-01")
	finalDir := filepath.Join(l.logDirectory, currentYear, currentMonth, currentDate)
	if err := os.MkdirAll(finalDir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	filename := filepath.Join(finalDir, fmt.Sprintf("%s_%s.log", tableName, dateStr))

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
		fileLogger.Printf("[%.3fms] [rows:%d] %s\n[ERROR] %v", float64(elapsed.Nanoseconds())/1e6, rows, sql, err)
	case elapsed > 200*time.Millisecond && l.logLevel >= logger.Warn: // Log slow queries
		fileLogger.Printf("[%.3fms] [rows:%d] %s\n[SLOW QUERY]", float64(elapsed.Nanoseconds())/1e6, rows, sql)
	case l.logLevel >= logger.Info:
		fileLogger.Printf("[%.3fms] [rows:%d] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
	}
}
