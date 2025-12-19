package core

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/JB-SelfCompany/yggpeers"
)

// DiscoveredPeer represents a discovered Yggdrasil peer from the public peer list
type DiscoveredPeer struct {
	// Address is the full peer URI (e.g., "tls://host:port")
	Address string `json:"address" toml:"address"`

	// Protocol is the connection protocol (tcp, tls, quic, ws, wss)
	Protocol string `json:"protocol" toml:"protocol"`

	// Region is the geographical region of the peer
	Region string `json:"region" toml:"region"`

	// RTT is the round-trip time in milliseconds
	RTT int64 `json:"rtt" toml:"rtt"`

	// Available indicates if the peer is currently reachable
	Available bool `json:"available" toml:"available"`

	// ResponseMS is the response time from publicnodes.json
	ResponseMS int `json:"response_ms" toml:"response_ms"`

	// LastSeen is the Unix timestamp when peer was last seen
	LastSeen int64 `json:"last_seen" toml:"last_seen"`
}

// ToPeerConfig converts a DiscoveredPeer to a PeerConfig
func (dp *DiscoveredPeer) ToPeerConfig() PeerConfig {
	return PeerConfig{
		Address: dp.Address,
		Enabled: true,
	}
}

// PeerDiscoveryManager handles peer discovery operations
type PeerDiscoveryManager struct {
	manager *yggpeers.Manager
}

// NewPeerDiscoveryManager creates a new peer discovery manager
func NewPeerDiscoveryManager() *PeerDiscoveryManager {
	manager := yggpeers.NewManager(
		yggpeers.WithCacheTTL(24*time.Hour),
		yggpeers.WithTimeout(5*time.Second),
	)

	return &PeerDiscoveryManager{
		manager: manager,
	}
}

// GetBatchingParams returns optimal batching parameters based on platform
// Desktop systems typically have better network connections than mobile
func GetBatchingParams() (batchSize, concurrency, pauseMs int) {
	// Desktop defaults - more aggressive than mobile
	// These values are optimized for typical desktop connections (50-1000 Mbps)
	return 40, 20, 150
}

// PeerDiscoveryProgress represents progress information during peer discovery
type PeerDiscoveryProgress struct {
	Current        int `json:"current"`
	Total          int `json:"total"`
	AvailableCount int `json:"available_count"`
}

// PeerDiscoveryResult represents the final result of peer discovery
type PeerDiscoveryResult struct {
	Peers     []DiscoveredPeer `json:"peers"`
	Total     int              `json:"total"`
	Available int              `json:"available"`
	Elapsed   string           `json:"elapsed"`
}

// FindAvailablePeers finds available peers with the given filters
// protocols: comma-separated list (e.g., "tcp,tls,quic"), empty for all
// region: filter by region, empty for all
// maxRTTMs: maximum RTT in milliseconds, 0 for no limit
// progressCallback: called periodically with progress updates
func (pdm *PeerDiscoveryManager) FindAvailablePeers(
	ctx context.Context,
	protocols string,
	region string,
	maxRTTMs int,
	progressCallback func(PeerDiscoveryProgress),
) (*PeerDiscoveryResult, error) {
	startTime := time.Now()

	// Parse protocols
	var protoList []yggpeers.Protocol
	if protocols != "" {
		for _, p := range splitAndTrim(protocols, ",") {
			protoList = append(protoList, yggpeers.Protocol(p))
		}
	}

	// Parse region
	var regionList []string
	if region != "" {
		regionList = []string{region}
	}

	// Build filter options
	filter := &yggpeers.FilterOptions{
		Protocols: protoList,
		Regions:   regionList,
		OnlyUp:    true,
	}

	if maxRTTMs > 0 {
		filter.MaxRTT = time.Duration(maxRTTMs) * time.Millisecond
	} else {
		filter.MaxRTT = 5 * time.Second
	}

	// Get all peers
	allPeers, err := pdm.manager.GetPeers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch peers: %w", err)
	}

	// Pre-filter by protocol and region
	filtered := pdm.manager.FilterPeers(allPeers, filter, yggpeers.SortByRTT)
	total := len(filtered)

	if progressCallback != nil {
		progressCallback(PeerDiscoveryProgress{
			Current:        0,
			Total:          total,
			AvailableCount: 0,
		})
	}

	// Get batching parameters
	batchSize, concurrency, pauseMs := GetBatchingParams()

	// Check peers in batches
	availablePeers := make([]DiscoveredPeer, 0)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		batch := filtered[i:end]

		// Check batch
		err := pdm.manager.CheckPeers(ctx, batch, concurrency)
		if err != nil {
			log.Printf("Error checking batch: %v", err)
		}

		// Collect available peers
		for _, peer := range batch {
			if peer.Available && matchesMaxRTT(peer, filter.MaxRTT) {
				dp := DiscoveredPeer{
					Address:    peer.Address,
					Protocol:   string(peer.Protocol),
					Region:     peer.Region,
					RTT:        peer.RTT.Milliseconds(),
					Available:  peer.Available,
					ResponseMS: peer.ResponseMS,
					LastSeen:   peer.LastSeen,
				}
				availablePeers = append(availablePeers, dp)
			}
		}

		// Send progress update
		if progressCallback != nil {
			progressCallback(PeerDiscoveryProgress{
				Current:        end,
				Total:          total,
				AvailableCount: len(availablePeers),
			})
		}

		// Pause between batches (optional, for rate limiting)
		if end < total && pauseMs > 0 {
			time.Sleep(time.Duration(pauseMs) * time.Millisecond)
		}
	}

	elapsed := time.Since(startTime)

	return &PeerDiscoveryResult{
		Peers:     availablePeers,
		Total:     total,
		Available: len(availablePeers),
		Elapsed:   elapsed.String(),
	}, nil
}

// GetAvailableRegions returns a list of all available regions
func (pdm *PeerDiscoveryManager) GetAvailableRegions(ctx context.Context) ([]string, error) {
	peers, err := pdm.manager.GetPeers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get peers: %w", err)
	}

	// Extract unique regions
	regionMap := make(map[string]bool)
	for _, peer := range peers {
		if peer.Region != "" {
			regionMap[peer.Region] = true
		}
	}

	regions := make([]string, 0, len(regionMap))
	for region := range regionMap {
		regions = append(regions, region)
	}

	return regions, nil
}

// CheckCustomPeers checks a list of user-provided peer URIs
func (pdm *PeerDiscoveryManager) CheckCustomPeers(
	ctx context.Context,
	peerURIs []string,
	concurrency int,
) ([]DiscoveredPeer, error) {
	// Convert URIs to peers
	peers := make([]*yggpeers.Peer, 0, len(peerURIs))
	for _, uri := range peerURIs {
		peer, err := parsePeerURI(uri)
		if err != nil {
			log.Printf("Warning: invalid peer URI %s: %v", uri, err)
			continue
		}
		peers = append(peers, peer)
	}

	if len(peers) == 0 {
		return nil, fmt.Errorf("no valid peer URIs provided")
	}

	// Check all peers
	if concurrency <= 0 {
		_, concurrency, _ = GetBatchingParams()
	}

	err := pdm.manager.CheckPeers(ctx, peers, concurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to check peers: %w", err)
	}

	// Convert results
	results := make([]DiscoveredPeer, 0, len(peers))
	for _, peer := range peers {
		dp := DiscoveredPeer{
			Address:    peer.Address,
			Protocol:   string(peer.Protocol),
			Region:     peer.Region,
			RTT:        peer.RTT.Milliseconds(),
			Available:  peer.Available,
			ResponseMS: peer.ResponseMS,
			LastSeen:   peer.LastSeen,
		}
		results = append(results, dp)
	}

	return results, nil
}

// Helper functions

func matchesMaxRTT(peer *yggpeers.Peer, maxRTT time.Duration) bool {
	if maxRTT == 0 {
		return true
	}
	return peer.RTT <= maxRTT
}

func parsePeerURI(uri string) (*yggpeers.Peer, error) {
	// Parse protocol from URI (e.g., "tls://host:port")
	parts := splitAndTrim(uri, "://")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid URI format: %s", uri)
	}

	protocol := yggpeers.Protocol(parts[0])

	// Extract host and port (simple parsing, yggpeers will validate)
	return &yggpeers.Peer{
		Address:  uri,
		Protocol: protocol,
	}, nil
}

func splitAndTrim(s, sep string) []string {
	parts := make([]string, 0)
	for _, part := range split(s, sep) {
		trimmed := trim(part)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

func split(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	// Simple split implementation
	result := make([]string, 0)
	start := 0
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

func trim(s string) string {
	// Trim leading and trailing whitespace
	start := 0
	end := len(s)

	for start < end && isSpace(s[start]) {
		start++
	}

	for end > start && isSpace(s[end-1]) {
		end--
	}

	return s[start:end]
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

// GetSystemInfo returns system information for optimal batching parameters
func GetSystemInfo() string {
	return fmt.Sprintf("OS: %s, Arch: %s, CPUs: %d",
		runtime.GOOS, runtime.GOARCH, runtime.NumCPU())
}
