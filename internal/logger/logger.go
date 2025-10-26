package logger

import (
	"fmt"
	"log"
	"os"
)

// Level represents log level
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

var (
	currentLevel Level = LevelInfo
	verbose      bool  = false
)

// Init initializes the logger
func Init(verboseMode bool) {
	verbose = verboseMode
	if verbose {
		currentLevel = LevelDebug
	}
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime)
}

// Debug logs a debug message
func Debug(format string, args ...interface{}) {
	if currentLevel <= LevelDebug {
		log.Printf("[DEBUG] "+format, args...)
	}
}

// Info logs an info message
func Info(format string, args ...interface{}) {
	if currentLevel <= LevelInfo {
		log.Printf("[INFO] "+format, args...)
	}
}

// Warn logs a warning message
func Warn(format string, args ...interface{}) {
	if currentLevel <= LevelWarn {
		log.Printf("[WARN] "+format, args...)
	}
}

// Error logs an error message
func Error(format string, args ...interface{}) {
	if currentLevel <= LevelError {
		log.Printf("[ERROR] "+format, args...)
	}
}

// Fatal logs a fatal message and exits
func Fatal(format string, args ...interface{}) {
	log.Fatalf("[FATAL] "+format, args...)
}

// Audit logs an audit message (always logged)
func Audit(action, user, resource string, success bool) {
	status := "SUCCESS"
	if !success {
		status = "FAILURE"
	}
	message := fmt.Sprintf("[AUDIT] action=%s user=%s resource=%s status=%s", 
		action, user, resource, status)
	log.Println(message)
	
	// TODO: Also write to audit log file
}
