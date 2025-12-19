package main

import (
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/models"
)

// Re-export DTOs from internal/models for Wails bindings
// Wails generates TypeScript types from these exported types in package main

// ServiceStatusDTO contains the current service status and addresses
type ServiceStatusDTO = models.ServiceStatusDTO

// PeerInfoDTO contains information about a Yggdrasil network peer
type PeerInfoDTO = models.PeerInfoDTO

// ConfigDTO contains the application configuration
type ConfigDTO = models.ConfigDTO

// PeerConfigDTO represents a peer configuration
type PeerConfigDTO = models.PeerConfigDTO

// LogEventDTO represents a log message event
type LogEventDTO = models.LogEventDTO

// MailEventDTO represents a mail-related event
type MailEventDTO = models.MailEventDTO

// ConnectionEventDTO represents a connection status change
type ConnectionEventDTO = models.ConnectionEventDTO

// BackupOptionsDTO contains options for creating a backup
type BackupOptionsDTO = models.BackupOptionsDTO

// RestoreOptionsDTO contains options for restoring from a backup
type RestoreOptionsDTO = models.RestoreOptionsDTO

// ResultDTO represents a generic operation result
type ResultDTO = models.ResultDTO
