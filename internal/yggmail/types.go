package yggmail

import (
	"sync"
	"time"
)

// ServiceStatus represents the current state of the yggmail service
type ServiceStatus int

const (
	// StatusStopped indicates the service is not running
	StatusStopped ServiceStatus = iota
	// StatusStarting indicates the service is initializing
	StatusStarting
	// StatusRunning indicates the service is running normally
	StatusRunning
	// StatusStopping indicates the service is shutting down
	StatusStopping
	// StatusError indicates the service encountered an error
	StatusError
)

// String returns the string representation of ServiceStatus
func (s ServiceStatus) String() string {
	switch s {
	case StatusStopped:
		return "Stopped"
	case StatusStarting:
		return "Starting"
	case StatusRunning:
		return "Running"
	case StatusStopping:
		return "Stopping"
	case StatusError:
		return "Error"
	default:
		return "Unknown"
	}
}

// PeerInfo contains information about a connected peer
type PeerInfo struct {
	// Address is the peer's URI (e.g., "tls://example.com:12345")
	Address string
	// Latency is the round-trip time to the peer in milliseconds
	Latency int64
	// Status indicates if the peer is connected (Up=true) or disconnected (Up=false)
	Status bool
	// Inbound indicates if this is an inbound connection
	Inbound bool
	// Key is the peer's public key (hex encoded)
	Key string
	// Uptime is the connection duration in seconds
	Uptime int64
	// RXBytes is the total bytes received from this peer
	RXBytes int64
	// TXBytes is the total bytes transmitted to this peer
	TXBytes int64
	// RXRate is the current receive rate in bytes/sec
	RXRate int64
	// TXRate is the current transmit rate in bytes/sec
	TXRate int64
	// LastError contains the last error message if any
	LastError string
}

// LogEvent represents a log message from the yggmail service
type LogEvent struct {
	// Timestamp is when the log event occurred
	Timestamp time.Time
	// Level is the log level (e.g., "INFO", "ERROR", "DEBUG")
	Level string
	// Tag is the component that generated the log (e.g., "Yggmail", "Transport")
	Tag string
	// Message is the log message
	Message string
}

// MailEvent represents a mail-related event
type MailEvent struct {
	// Timestamp is when the event occurred
	Timestamp time.Time
	// Type is the event type (e.g., "new_mail", "sent", "error")
	Type string
	// Mailbox is the mailbox name (e.g., "INBOX")
	Mailbox string
	// From is the sender's email address
	From string
	// To is the recipient's email address
	To string
	// Subject is the email subject
	Subject string
	// MailID is the internal mail ID
	MailID int
	// ErrorMessage contains error details if Type is "error"
	ErrorMessage string
}

// ConnectionEvent represents a connection status change
type ConnectionEvent struct {
	// Timestamp is when the event occurred
	Timestamp time.Time
	// Type is the event type (e.g., "connected", "disconnected", "error")
	Type string
	// Peer is the peer address
	Peer string
	// ErrorMessage contains error details if Type is "error"
	ErrorMessage string
}

// EventChannels contains channels for receiving events from the service
// These channels are buffered to prevent blocking the service
type EventChannels struct {
	// Log receives log events from the service
	Log chan LogEvent
	// Mail receives mail-related events
	Mail chan MailEvent
	// Connection receives connection status changes
	Connection chan ConnectionEvent
	// closed indicates whether channels have been closed
	closed bool
	// mu protects the closed flag
	mu sync.RWMutex
}

// NewEventChannels creates a new EventChannels with buffered channels
// Buffer size is set to prevent blocking during high-frequency events
func NewEventChannels() *EventChannels {
	return &EventChannels{
		Log:        make(chan LogEvent, 100),
		Mail:       make(chan MailEvent, 50),
		Connection: make(chan ConnectionEvent, 50),
		closed:     false,
	}
}

// Close closes all event channels
// Call this when you're done receiving events to prevent goroutine leaks
func (ec *EventChannels) Close() {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	if ec.closed {
		return
	}
	ec.closed = true

	if ec.Log != nil {
		close(ec.Log)
	}
	if ec.Mail != nil {
		close(ec.Mail)
	}
	if ec.Connection != nil {
		close(ec.Connection)
	}
}

// IsClosed returns true if channels have been closed
func (ec *EventChannels) IsClosed() bool {
	ec.mu.RLock()
	defer ec.mu.RUnlock()
	return ec.closed
}

// MessageSizeLimitCheckResult represents the result of recipient message size limit check
type MessageSizeLimitCheckResult struct {
	// CanSend indicates whether message can be sent (size within limit)
	CanSend bool
	// ErrorMessage contains error message if CanSend is false
	ErrorMessage string
	// RecipientAddr is the recipient address that was checked
	RecipientAddr string
	// MessageSizeMB is the message size in megabytes
	MessageSizeMB float64
}
