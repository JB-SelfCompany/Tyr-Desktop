package autoconfig

import (
	"encoding/xml"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
)

type Server struct {
	// Configuration
	mailDomain   string // e.g., "yggmail" for addresses like user@yggmail
	smtpHost     string // e.g., "127.0.0.1"
	smtpPort     string // e.g., "1025"
	imapHost     string // e.g., "127.0.0.1"
	imapPort     string // e.g., "1143"
	listenAddr   string // e.g., "127.0.0.1:8080"
	displayName  string // e.g., "Yggmail"
	shortName    string // e.g., "Yggmail"

	// HTTP server
	server   *http.Server
	listener net.Listener

	// State management
	mu      sync.RWMutex
	running bool
}

// ServerConfig contains configuration for the autoconfig server
type ServerConfig struct {
	MailDomain  string // Email domain (e.g., "yggmail")
	SMTPHost    string // SMTP server hostname
	SMTPPort    string // SMTP server port
	IMAPHost    string // IMAP server hostname
	IMAPPort    string // IMAP server port
	ListenAddr  string // Address to listen on (e.g., "127.0.0.1:8080")
	DisplayName string // Display name for email provider
	ShortName   string // Short name for email provider
}

// ClientConfig represents the XML structure for Thunderbird/DeltaChat autoconfiguration
type ClientConfig struct {
	XMLName       xml.Name      `xml:"clientConfig"`
	Version       string        `xml:"version,attr"`
	EmailProvider EmailProvider `xml:"emailProvider"`
}

// EmailProvider contains the email provider configuration
type EmailProvider struct {
	ID              string           `xml:"id,attr"`
	Domains         []string         `xml:"domain"`
	DisplayName     string           `xml:"displayName"`
	DisplayShortName string          `xml:"displayShortName"`
	IncomingServer  IncomingServer   `xml:"incomingServer"`
	OutgoingServer  OutgoingServer   `xml:"outgoingServer"`
}

// IncomingServer contains IMAP server configuration
type IncomingServer struct {
	Type           string `xml:"type,attr"`
	Hostname       string `xml:"hostname"`
	Port           string `xml:"port"`
	SocketType     string `xml:"socketType"`
	Authentication string `xml:"authentication"`
	Username       string `xml:"username"`
}

// OutgoingServer contains SMTP server configuration
type OutgoingServer struct {
	Type           string `xml:"type,attr"`
	Hostname       string `xml:"hostname"`
	Port           string `xml:"port"`
	SocketType     string `xml:"socketType"`
	Authentication string `xml:"authentication"`
	Username       string `xml:"username"`
}

// NewServer creates a new autoconfiguration server
func NewServer(config ServerConfig) (*Server, error) {
	// Validate configuration
	if config.MailDomain == "" {
		return nil, fmt.Errorf("mail domain cannot be empty")
	}
	if config.SMTPHost == "" {
		return nil, fmt.Errorf("SMTP host cannot be empty")
	}
	if config.SMTPPort == "" {
		return nil, fmt.Errorf("SMTP port cannot be empty")
	}
	if config.IMAPHost == "" {
		return nil, fmt.Errorf("IMAP host cannot be empty")
	}
	if config.IMAPPort == "" {
		return nil, fmt.Errorf("IMAP port cannot be empty")
	}
	if config.ListenAddr == "" {
		config.ListenAddr = "127.0.0.1:8080"
	}
	if config.DisplayName == "" {
		config.DisplayName = "Yggmail"
	}
	if config.ShortName == "" {
		config.ShortName = "Yggmail"
	}

	s := &Server{
		mailDomain:  config.MailDomain,
		smtpHost:    config.SMTPHost,
		smtpPort:    config.SMTPPort,
		imapHost:    config.IMAPHost,
		imapPort:    config.IMAPPort,
		listenAddr:  config.ListenAddr,
		displayName: config.DisplayName,
		shortName:   config.ShortName,
	}

	return s, nil
}

// Start starts the autoconfiguration HTTP server
func (s *Server) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("server already running")
	}

	// Create HTTP server
	mux := http.NewServeMux()

	// Register handlers for both standard autoconfig paths
	mux.HandleFunc("/.well-known/autoconfig/mail/config-v1.1.xml", s.handleAutoconfig)
	mux.HandleFunc("/mail/config-v1.1.xml", s.handleAutoconfig)
	mux.HandleFunc("/autoconfig/mail/config-v1.1.xml", s.handleAutoconfig)

	// Add a root handler for debugging
	mux.HandleFunc("/", s.handleRoot)

	s.server = &http.Server{
		Addr:    s.listenAddr,
		Handler: s.logRequest(mux),
	}

	// Create listener
	listener, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}
	s.listener = listener

	// Start server in background
	go func() {
		log.Printf("Autoconfig server listening on %s", s.listenAddr)
		if err := s.server.Serve(s.listener); err != nil && err != http.ErrServerClosed {
			log.Printf("Autoconfig server error: %v", err)
		}
	}()

	s.running = true
	return nil
}

// Stop stops the autoconfiguration HTTP server
func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}

	if s.server != nil {
		if err := s.server.Close(); err != nil {
			return fmt.Errorf("failed to stop server: %w", err)
		}
	}

	if s.listener != nil {
		s.listener.Close()
	}

	s.running = false
	log.Printf("Autoconfig server stopped")
	return nil
}

// IsRunning returns true if the server is running
func (s *Server) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// GetListenAddr returns the address the server is listening on
func (s *Server) GetListenAddr() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.listenAddr
}

// handleAutoconfig handles requests for the autoconfiguration XML
func (s *Server) handleAutoconfig(w http.ResponseWriter, r *http.Request) {
	log.Printf("Autoconfig request from %s: %s", r.RemoteAddr, r.URL.Path)

	// Generate configuration XML
	config := s.generateConfig()

	// Marshal to XML
	xmlData, err := xml.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal XML: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set headers
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Write XML declaration and data
	w.Write([]byte(xml.Header))
	w.Write(xmlData)

	log.Printf("Autoconfig response sent successfully")
}

// handleRoot handles root requests for debugging
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
    <title>Yggmail Autoconfiguration Server</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px; }
        h1 { color: #333; }
        .info { background: #f0f0f0; padding: 15px; border-radius: 5px; margin: 20px 0; }
        code { background: #e0e0e0; padding: 2px 6px; border-radius: 3px; }
        ul { line-height: 1.8; }
    </style>
</head>
<body>
    <h1>Yggmail Autoconfiguration Server</h1>
    <div class="info">
        <p><strong>Status:</strong> Running</p>
        <p><strong>Mail Domain:</strong> %s</p>
        <p><strong>SMTP:</strong> %s:%s (plain text, no encryption)</p>
        <p><strong>IMAP:</strong> %s:%s (plain text, no encryption)</p>
    </div>

    <h2>Autoconfiguration URLs</h2>
    <ul>
        <li><a href="/.well-known/autoconfig/mail/config-v1.1.xml">/.well-known/autoconfig/mail/config-v1.1.xml</a></li>
        <li><a href="/mail/config-v1.1.xml">/mail/config-v1.1.xml</a></li>
        <li><a href="/autoconfig/mail/config-v1.1.xml">/autoconfig/mail/config-v1.1.xml</a></li>
    </ul>

    <h2>DeltaChat Setup</h2>
    <p>DeltaChat will automatically discover these settings when you configure your account.</p>
    <p>Make sure DeltaChat is configured to check <code>http://localhost:8080</code> for autoconfiguration.</p>

    <h2>Manual Configuration</h2>
    <p>If autoconfiguration doesn't work, use these settings:</p>
    <ul>
        <li><strong>IMAP Server:</strong> %s:%s</li>
        <li><strong>SMTP Server:</strong> %s:%s</li>
        <li><strong>Encryption:</strong> None (localhost only)</li>
        <li><strong>Authentication:</strong> Password (cleartext)</li>
        <li><strong>Username:</strong> Your full email address</li>
    </ul>
</body>
</html>`, s.mailDomain, s.smtpHost, s.smtpPort, s.imapHost, s.imapPort, s.imapHost, s.imapPort, s.smtpHost, s.smtpPort)
		return
	}

	http.NotFound(w, r)
}

// generateConfig creates the autoconfiguration XML structure
func (s *Server) generateConfig() ClientConfig {
	return ClientConfig{
		Version: "1.1",
		EmailProvider: EmailProvider{
			ID:               s.mailDomain,
			Domains:          []string{s.mailDomain},
			DisplayName:      s.displayName,
			DisplayShortName: s.shortName,
			IncomingServer: IncomingServer{
				Type:           "imap",
				Hostname:       s.imapHost,
				Port:           s.imapPort,
				SocketType:     "plain", // No SSL for localhost
				Authentication: "password-cleartext",
				Username:       "%EMAILADDRESS%", // Placeholder replaced by client
			},
			OutgoingServer: OutgoingServer{
				Type:           "smtp",
				Hostname:       s.smtpHost,
				Port:           s.smtpPort,
				SocketType:     "plain", // No SSL for localhost
				Authentication: "password-cleartext",
				Username:       "%EMAILADDRESS%", // Placeholder replaced by client
			},
		},
	}
}

// logRequest is a middleware that logs all HTTP requests
func (s *Server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("HTTP %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Log query parameters
		if len(r.URL.Query()) > 0 {
			log.Printf("Query params: %v", r.URL.Query())
		}

		// Log email address if provided
		if email := r.URL.Query().Get("emailaddress"); email != "" {
			log.Printf("Email address: %s", email)
		}

		next.ServeHTTP(w, r)
	})
}

// ParseSMTPAddress parses an SMTP address string into host and port
// Example: "127.0.0.1:1025" -> ("127.0.0.1", "1025")
func ParseSMTPAddress(addr string) (host, port string) {
	parts := strings.Split(addr, ":")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return addr, "25"
}

// ParseIMAPAddress parses an IMAP address string into host and port
// Example: "127.0.0.1:1143" -> ("127.0.0.1", "1143")
func ParseIMAPAddress(addr string) (host, port string) {
	parts := strings.Split(addr, ":")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return addr, "143"
}
