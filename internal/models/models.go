package models

import "time"

// DTOs for Wails v2 auto-generation of TypeScript types
// These structs are exported and will be converted to TypeScript interfaces
// Run `wails generate module` to generate TypeScript bindings

// ServiceStatusDTO contains the current service status and addresses
type ServiceStatusDTO struct {
	// Status is the current service state (Stopped, Starting, Running, Stopping, Error)
	Status string `json:"status"`
	// Running indicates if the service is currently running
	Running bool `json:"running"`
	// MailAddress is the yggmail address (e.g., "user@ygg")
	MailAddress string `json:"mailAddress"`
	// SMTPAddress is the local SMTP server address
	SMTPAddress string `json:"smtpAddress"`
	// IMAPAddress is the local IMAP server address
	IMAPAddress string `json:"imapAddress"`
	// DatabasePath is the path to the yggmail database file
	DatabasePath string `json:"databasePath"`
	// ErrorMessage contains error details if Status is "Error"
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// PeerInfoDTO contains information about a Yggdrasil network peer
type PeerInfoDTO struct {
	// Address is the peer URI (e.g., "tls://example.com:12345")
	Address string `json:"address"`
	// Enabled indicates if this peer is enabled in configuration
	Enabled bool `json:"enabled"`
	// Connected indicates if currently connected to this peer
	Connected bool `json:"connected"`
	// Latency is the round-trip time in milliseconds (0 if not connected)
	Latency int64 `json:"latency"`
	// Uptime is the connection duration in seconds (0 if not connected)
	Uptime int64 `json:"uptime"`
	// RXBytes is the total bytes received from this peer
	RXBytes int64 `json:"rxBytes"`
	// TXBytes is the total bytes transmitted to this peer
	TXBytes int64 `json:"txBytes"`
	// RXRate is the current receive rate in bytes/sec
	RXRate int64 `json:"rxRate"`
	// TXRate is the current transmit rate in bytes/sec
	TXRate int64 `json:"txRate"`
	// LastError contains the last error message if any
	LastError string `json:"lastError,omitempty"`
}

// ConfigDTO contains the application configuration
type ConfigDTO struct {
	// OnboardingComplete indicates if initial setup is complete
	OnboardingComplete bool `json:"onboardingComplete"`
	// Peers is the list of Yggdrasil network peers
	Peers []PeerConfigDTO `json:"peers"`
	// Language is the UI language ("en", "ru")
	Language string `json:"language"`
	// Theme is the UI theme ("light", "dark", "system")
	Theme string `json:"theme"`
	// AutoStart indicates if service starts on system boot
	AutoStart bool `json:"autoStart"`
	// SMTPAddress is the local SMTP server address
	SMTPAddress string `json:"smtpAddress"`
	// IMAPAddress is the local IMAP server address
	IMAPAddress string `json:"imapAddress"`
	// DatabasePath is the path to the yggmail database file
	DatabasePath string `json:"databasePath"`
}

// PeerConfigDTO represents a peer configuration
type PeerConfigDTO struct {
	// Address is the peer URI
	Address string `json:"address"`
	// Enabled indicates if this peer is enabled
	Enabled bool `json:"enabled"`
}

// LogEventDTO represents a log message event
type LogEventDTO struct {
	// Timestamp is when the event occurred (RFC3339 format)
	Timestamp string `json:"timestamp"`
	// Level is the log level (INFO, WARN, ERROR, DEBUG)
	Level string `json:"level"`
	// Tag is the component that generated the log
	Tag string `json:"tag"`
	// Message is the log message
	Message string `json:"message"`
}

// MailEventDTO represents a mail-related event
type MailEventDTO struct {
	// Timestamp is when the event occurred (RFC3339 format)
	Timestamp string `json:"timestamp"`
	// Type is the event type (new_mail, sent, error)
	Type string `json:"type"`
	// Mailbox is the mailbox name (e.g., "INBOX")
	Mailbox string `json:"mailbox"`
	// From is the sender's email address
	From string `json:"from"`
	// To is the recipient's email address
	To string `json:"to"`
	// Subject is the email subject
	Subject string `json:"subject"`
	// MailID is the internal mail ID
	MailID int `json:"mailId"`
	// ErrorMessage contains error details if Type is "error"
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// ConnectionEventDTO represents a connection status change
type ConnectionEventDTO struct {
	// Timestamp is when the event occurred (RFC3339 format)
	Timestamp string `json:"timestamp"`
	// Type is the event type (connected, disconnected, error)
	Type string `json:"type"`
	// Peer is the peer address
	Peer string `json:"peer"`
	// ErrorMessage contains error details if Type is "error"
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// BackupOptionsDTO contains options for creating a backup
type BackupOptionsDTO struct {
	// BackupPath is the path where the backup file will be saved
	BackupPath string `json:"backupPath"`
	// IncludeDatabase indicates if yggmail.db should be included
	IncludeDatabase bool `json:"includeDatabase"`
	// Password is the encryption password
	Password string `json:"password"`
}

// RestoreOptionsDTO contains options for restoring from a backup
type RestoreOptionsDTO struct {
	// BackupPath is the path to the backup file
	BackupPath string `json:"backupPath"`
	// Password is the decryption password
	Password string `json:"password"`
}

// ResultDTO represents a generic operation result
type ResultDTO struct {
	// Success indicates if the operation succeeded
	Success bool `json:"success"`
	// Message contains additional information or error message
	Message string `json:"message,omitempty"`
	// Data contains optional result data (JSON-encoded)
	Data string `json:"data,omitempty"`
}

// MessageSizeLimitCheckResultDTO represents the result of recipient message size limit check
type MessageSizeLimitCheckResultDTO struct {
	// CanSend indicates whether message can be sent (size within limit)
	CanSend bool `json:"canSend"`
	// ErrorMessage contains error message if CanSend is false
	ErrorMessage string `json:"errorMessage,omitempty"`
	// RecipientAddr is the recipient address that was checked
	RecipientAddr string `json:"recipientAddr"`
	// MessageSizeMB is the message size in megabytes
	MessageSizeMB float64 `json:"messageSizeMB"`
}

// StorageStatsDTO contains information about storage usage
type StorageStatsDTO struct {
	// DatabaseSizeMB is the size of the yggmail.db file in megabytes
	DatabaseSizeMB float64 `json:"databaseSizeMB"`
	// FilesSizeMB is the total size of all stored message files in megabytes
	FilesSizeMB float64 `json:"filesSizeMB"`
	// TotalSizeMB is the total storage usage in megabytes
	TotalSizeMB float64 `json:"totalSizeMB"`
	// MaxMessageSizeMB is the current maximum message size limit in megabytes
	MaxMessageSizeMB int64 `json:"maxMessageSizeMB"`
}

// Helper functions to convert internal types to DTOs

// formatTimestamp converts time.Time to RFC3339 string
func formatTimestamp(t time.Time) string {
	return t.Format(time.RFC3339)
}
