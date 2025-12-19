package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp time.Time
	Level     string
	Message   string
}

var (
	// Global log buffer
	logBuffer     []LogEntry
	logBufferLock sync.Mutex
	maxLogBuffer  = 1000
)

// AddLogEntry adds a new log entry to the global buffer
func AddLogEntry(level, message string) {
	logBufferLock.Lock()
	defer logBufferLock.Unlock()

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
	}

	logBuffer = append(logBuffer, entry)

	// Trim buffer if it exceeds max size
	if len(logBuffer) > maxLogBuffer {
		logBuffer = logBuffer[len(logBuffer)-maxLogBuffer:]
	}
}

// GetLogEntries returns a copy of the log buffer
func GetLogEntries() []LogEntry {
	logBufferLock.Lock()
	defer logBufferLock.Unlock()

	entries := make([]LogEntry, len(logBuffer))
	copy(entries, logBuffer)
	return entries
}

// ClearLogs clears all log entries
func ClearLogs() {
	logBufferLock.Lock()
	defer logBufferLock.Unlock()
	logBuffer = []LogEntry{}
}

// LogWriter is a custom writer that captures log output
type LogWriter struct {
	originalWriter io.Writer
}

// NewLogWriter creates a new log writer that captures logs
func NewLogWriter(originalWriter io.Writer) *LogWriter {
	return &LogWriter{
		originalWriter: originalWriter,
	}
}

// Write implements io.Writer interface
func (lw *LogWriter) Write(p []byte) (n int, err error) {
	// Write to original output (console)
	n, err = lw.originalWriter.Write(p)

	// Parse log message to extract level and message
	message := string(p)
	level, msg := parseLogMessage(message)

	// Add to log buffer for UI display
	AddLogEntry(level, msg)

	return n, err
}

// parseLogMessage extracts log level and message from log output
func parseLogMessage(logMsg string) (level, message string) {
	// Default level
	level = "INFO"

	// Clean up the message
	message = strings.TrimSpace(logMsg)

	// Try to extract explicit log level from format [LEVEL]
	if strings.Contains(message, "[ERROR]") {
		level = "ERROR"
		message = strings.ReplaceAll(message, "[ERROR]", "")
	} else if strings.Contains(message, "[WARN]") {
		level = "WARN"
		message = strings.ReplaceAll(message, "[WARN]", "")
	} else if strings.Contains(message, "[INFO]") {
		level = "INFO"
		message = strings.ReplaceAll(message, "[INFO]", "")
	} else if strings.Contains(message, "[DEBUG]") {
		level = "DEBUG"
		message = strings.ReplaceAll(message, "[DEBUG]", "")
	} else {
		// Try to detect log level from message content (fallback)
		lowerMsg := strings.ToLower(message)
		if strings.Contains(lowerMsg, "error") || strings.Contains(lowerMsg, "failed") || strings.Contains(lowerMsg, "fatal") {
			level = "ERROR"
		} else if strings.Contains(lowerMsg, "warning") || strings.Contains(lowerMsg, "warn") {
			level = "WARN"
		} else if strings.Contains(lowerMsg, "debug") {
			level = "DEBUG"
		}
	}

	// Clean up message
	message = strings.TrimSpace(message)

	return level, message
}

// SetupLogging configures the application logging to use our custom logger
func SetupLogging() {
	// Get the current log output (usually os.Stderr)
	originalOutput := os.Stderr

	// Create custom log writer
	logWriter := NewLogWriter(originalOutput)

	// Set the log output to our custom writer
	log.SetOutput(logWriter)
	log.SetFlags(log.Ldate | log.Ltime)

	// Add initial log entry
	log.Println("Tyr logging system initialized")
}

// LogInfo logs an informational message
func LogInfo(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log.Printf("[INFO] %s", msg)
}

// LogError logs an error message
func LogError(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log.Printf("[ERROR] %s", msg)
}

// LogWarn logs a warning message
func LogWarn(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log.Printf("[WARN] %s", msg)
}

// LogDebug logs a debug message
func LogDebug(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log.Printf("[DEBUG] %s", msg)
}
