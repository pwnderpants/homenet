package logger

import (
	"fmt"
	"log"
	"time"
)

// LogLevel represents the logging level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var levelNames = map[LogLevel]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
}

// Logger provides structured logging functionality
type Logger struct {
	level LogLevel
}

// New creates a new logger instance
func New(level LogLevel) *Logger {
	return &Logger{level: level}
}

// NewDefault creates a new logger with INFO level
func NewDefault() *Logger {
	return New(INFO)
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// shouldLog checks if the given level should be logged
func (l *Logger) shouldLog(level LogLevel) bool {
	return level >= l.level
}

// formatMessage formats the log message with timestamp and level
func (l *Logger) formatMessage(level LogLevel, message string, args ...interface{}) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelName := levelNames[level]

	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	return fmt.Sprintf("[%s] %s: %s", timestamp, levelName, message)
}

// Debug logs a debug message
func (l *Logger) Debug(message string, args ...interface{}) {
	if l.shouldLog(DEBUG) {
		log.Println(l.formatMessage(DEBUG, message, args...))
	}
}

// Info logs an info message
func (l *Logger) Info(message string, args ...interface{}) {
	if l.shouldLog(INFO) {
		log.Println(l.formatMessage(INFO, message, args...))
	}
}

// Warn logs a warning message
func (l *Logger) Warn(message string, args ...interface{}) {
	if l.shouldLog(WARN) {
		log.Println(l.formatMessage(WARN, message, args...))
	}
}

// Error logs an error message
func (l *Logger) Error(message string, args ...interface{}) {
	if l.shouldLog(ERROR) {
		log.Println(l.formatMessage(ERROR, message, args...))
	}
}

// ErrorWithErr logs an error message with an error
func (l *Logger) ErrorWithErr(message string, err error, args ...interface{}) {
	if l.shouldLog(ERROR) {
		if len(args) > 0 {
			message = fmt.Sprintf(message, args...)
		}
		if err != nil {
			message = fmt.Sprintf("%s: %v", message, err)
		}
		log.Println(l.formatMessage(ERROR, message))
	}
}

// Global logger instance
var defaultLogger = NewDefault()

// SetGlobalLevel sets the global logger level
func SetGlobalLevel(level LogLevel) {
	defaultLogger.SetLevel(level)
}

// Debug logs a debug message using the global logger
func Debug(message string, args ...interface{}) {
	defaultLogger.Debug(message, args...)
}

// Info logs an info message using the global logger
func Info(message string, args ...interface{}) {
	defaultLogger.Info(message, args...)
}

// Warn logs a warning message using the global logger
func Warn(message string, args ...interface{}) {
	defaultLogger.Warn(message, args...)
}

// Error logs an error message using the global logger
func Error(message string, args ...interface{}) {
	defaultLogger.Error(message, args...)
}

// ErrorWithErr logs an error message with an error using the global logger
func ErrorWithErr(message string, err error, args ...interface{}) {
	defaultLogger.ErrorWithErr(message, err, args...)
}
