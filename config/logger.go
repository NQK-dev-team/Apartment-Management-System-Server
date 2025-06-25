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

// ModelFileLogger implements gorm.logger.Interface
type ModelFileLogger struct {
	logDirectory string
	logLevel     logger.LogLevel
	fileHandles  map[string]*os.File // Cache for open log files, key is the full filename
	mu           sync.Mutex          // Mutex to protect concurrent access to fileHandles map
}

// NewModelFileLogger creates a new instance of our custom logger
func NewModelFileLogger(dir string, level logger.LogLevel) *ModelFileLogger {
	// Ensure the log directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}
	return &ModelFileLogger{
		logDirectory: dir,
		logLevel:     level,
		fileHandles:  make(map[string]*os.File),
	}
}

// Regular expressions to extract table names from SQL queries
var (
	// selectRegex = regexp.MustCompile(`(?i)SELECT\s+.*?\s+FROM\s+"?([a-zA-Z0-9_]+)"?`)
	intoRegex   = regexp.MustCompile(`(?i)INTO\s+"([^"]+)"`)
	updateRegex = regexp.MustCompile(`(?i)UPDATE\s+"([^"]+)"`)
	deleteRegex = regexp.MustCompile(`(?i)DELETE FROM\s+"([^"]+)"`)
)

// extractTableName parses SQL to find the table name
func extractTableName(sql string) string {
	var match []string

	// if match = selectRegex.FindStringSubmatch(sql); len(match) > 1 {
	// 	return match[1]
	// }
	if match = intoRegex.FindStringSubmatch(sql); len(match) > 1 {
		return match[1]
	}
	if match = updateRegex.FindStringSubmatch(sql); len(match) > 1 {
		return match[1]
	}
	if match = deleteRegex.FindStringSubmatch(sql); len(match) > 1 {
		return match[1]
	}
	return "gorm_default"
}

// LogMode sets the log level
func (l *ModelFileLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

// Info logs general information
func (l *ModelFileLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Info {
		log.Printf("[INFO] "+msg, data...)
	}
}

// Warn logs warnings
func (l *ModelFileLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Warn {
		log.Printf("[WARN] "+msg, data...)
	}
}

// Error logs errors
func (l *ModelFileLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Error {
		log.Printf("[ERROR] "+msg, data...)
	}
}

// Trace is the core method for logging SQL queries
func (l *ModelFileLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	tableName := extractTableName(sql)
	dateStr := time.Now().Format("2006-01-02") // YYYY-MM-DD format

	if tableName == "gorm_default" {
		return // Skip logging for default table
	}

	// Construct filename: logs/users-2025-06-25.log
	filename := filepath.Join(l.logDirectory, fmt.Sprintf("%s_%s.log", tableName, dateStr))

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
