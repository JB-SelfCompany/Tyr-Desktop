package events

import (
	"log"
	"time"

	"github.com/JB-SelfCompany/Tyr-Desktop/internal/core"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/models"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/yggmail"
)

// EventEmitter is a callback function that emits events to the frontend
type EventEmitter func(eventName string, data interface{})

// StatusUpdater is a callback function that updates the system tray status
type StatusUpdater func()

// StartEventMonitoring monitors backend event channels and forwards events to frontend
// This goroutine runs in the background and stops when shutdownChan is closed
// Returns a boolean channel that will be closed when monitoring stops
func StartEventMonitoring(
	sm *core.ServiceManager,
	emitFunc EventEmitter,
	updateStatusFunc StatusUpdater,
	shutdownChan <-chan struct{},
) {
	if sm == nil {
		log.Println("Service manager not initialized, skipping event monitoring")
		return
	}

	log.Println("Starting event monitoring...")

	// Get event channels from service manager
	eventChans := sm.GetEventChannels()
	if eventChans == nil {
		log.Println("Event channels not available")
		return
	}

	// Monitor events in a loop
	for {
		select {
		case <-shutdownChan:
			log.Println("Event monitoring stopped")
			return

		case logEvent, ok := <-eventChans.Log:
			if !ok {
				log.Println("Log event channel closed")
				return
			}
			dto := ConvertLogEvent(logEvent)
			emitFunc("service:log", dto)

		case mailEvent, ok := <-eventChans.Mail:
			if !ok {
				log.Println("Mail event channel closed")
				return
			}
			dto := ConvertMailEvent(mailEvent)
			emitFunc("service:mail", dto)

		case connEvent, ok := <-eventChans.Connection:
			if !ok {
				log.Println("Connection event channel closed")
				return
			}
			dto := ConvertConnectionEvent(connEvent)
			emitFunc("service:connection", dto)

			// Update system tray status when connection events occur
			if updateStatusFunc != nil {
				updateStatusFunc()
			}
		}
	}
}

// ConvertLogEvent converts yggmail.LogEvent to LogEventDTO
func ConvertLogEvent(event yggmail.LogEvent) models.LogEventDTO {
	return models.LogEventDTO{
		Timestamp: formatTimestamp(event.Timestamp),
		Level:     event.Level,
		Tag:       event.Tag,
		Message:   event.Message,
	}
}

// ConvertMailEvent converts yggmail.MailEvent to MailEventDTO
func ConvertMailEvent(event yggmail.MailEvent) models.MailEventDTO {
	return models.MailEventDTO{
		Timestamp:    formatTimestamp(event.Timestamp),
		Type:         event.Type,
		Mailbox:      event.Mailbox,
		From:         event.From,
		To:           event.To,
		Subject:      event.Subject,
		MailID:       event.MailID,
		ErrorMessage: event.ErrorMessage,
	}
}

// ConvertConnectionEvent converts yggmail.ConnectionEvent to ConnectionEventDTO
func ConvertConnectionEvent(event yggmail.ConnectionEvent) models.ConnectionEventDTO {
	return models.ConnectionEventDTO{
		Timestamp:    formatTimestamp(event.Timestamp),
		Type:         event.Type,
		Peer:         event.Peer,
		ErrorMessage: event.ErrorMessage,
	}
}

// formatTimestamp converts time.Time to RFC3339 string
func formatTimestamp(t time.Time) string {
	return t.Format(time.RFC3339)
}

// Event Names for Frontend:
// - "service:log"        -> LogEventDTO
// - "service:mail"       -> MailEventDTO
// - "service:connection" -> ConnectionEventDTO
// - "service:status"     -> string (status name)
//
// Frontend can subscribe to these events using:
// import { EventsOn } from '../wailsjs/runtime/runtime';
// EventsOn('service:log', (event: LogEventDTO) => { ... });
