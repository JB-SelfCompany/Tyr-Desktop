package peerdiscovery

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/JB-SelfCompany/Tyr-Desktop/internal/core"
)

// FindAvailablePeers discovers available Yggdrasil peers with the given filters
// Emits "peer:discovery:progress" events during discovery
// Returns the final result with list of discovered peers
//
// Parameters:
//   - ctx: Wails runtime context for emitting events to frontend
//   - discoveryCtx: Context for peer discovery operations (can be cancelled)
//   - cfg: Configuration to cache discovered peers
//   - protocols: comma-separated list (e.g., "tcp,tls,quic"), empty for all
//   - region: filter by region, empty for all
//   - maxRTTMs: maximum RTT in milliseconds, 0 for no limit (default: 5000)
//
// Returns: PeerDiscoveryResult with list of available peers
func FindAvailablePeers(
	ctx context.Context,
	discoveryCtx context.Context,
	cfg *core.Config,
	protocols string,
	region string,
	maxRTTMs int,
) (*core.PeerDiscoveryResult, error) {
	log.Printf("[FindAvailablePeers] Starting peer discovery: protocols=%s, region=%s, maxRTT=%dms",
		protocols, region, maxRTTMs)

	// Create peer discovery manager
	pdm := core.NewPeerDiscoveryManager()

	// Create timeout context (60 seconds)
	timeoutCtx, cancel := context.WithTimeout(discoveryCtx, 60*time.Second)
	defer cancel()

	// Progress callback - emit events to frontend
	progressCallback := func(progress core.PeerDiscoveryProgress) {
		if ctx != nil {
			runtime.EventsEmit(ctx, "peer:discovery:progress", progress)
		}
	}

	// Start discovery
	result, err := pdm.FindAvailablePeers(timeoutCtx, protocols, region, maxRTTMs, progressCallback)
	if err != nil {
		log.Printf("[FindAvailablePeers] Error: %v", err)
		return nil, fmt.Errorf("failed to discover peers: %w", err)
	}

	log.Printf("[FindAvailablePeers] Discovery complete: %d available out of %d checked (elapsed: %s)",
		result.Available, result.Total, result.Elapsed)

	// Cache the discovered peers
	if cfg != nil && len(result.Peers) > 0 {
		if err := cfg.CacheDiscoveredPeers(result.Peers); err != nil {
			log.Printf("[FindAvailablePeers] Warning: failed to cache peers: %v", err)
		}
		// Save config to persist cache
		if err := cfg.Save(); err != nil {
			log.Printf("[FindAvailablePeers] Warning: failed to save config: %v", err)
		}
	}

	return result, nil
}

// GetCachedDiscoveredPeers returns cached discovered peers if available
// Returns nil if cache is expired (> 24 hours) or empty
func GetCachedDiscoveredPeers(cfg *core.Config) []core.DiscoveredPeer {
	if cfg == nil {
		log.Println("[GetCachedDiscoveredPeers] Config not initialized")
		return nil
	}

	cached := cfg.GetCachedDiscoveredPeers()
	if cached == nil {
		log.Println("[GetCachedDiscoveredPeers] No cached peers or cache expired")
		return nil
	}

	log.Printf("[GetCachedDiscoveredPeers] Returning %d cached peers", len(cached))
	return cached
}

// ClearCachedDiscoveredPeers clears the cached discovered peers
func ClearCachedDiscoveredPeers(cfg *core.Config) error {
	if cfg == nil {
		return fmt.Errorf("config not initialized")
	}

	log.Println("[ClearCachedDiscoveredPeers] Clearing cache")
	cfg.ClearCachedDiscoveredPeers()

	// Save config to persist
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// GetAvailableRegions returns a list of all available peer regions
func GetAvailableRegions(discoveryCtx context.Context) ([]string, error) {
	log.Println("[GetAvailableRegions] Fetching available regions")

	// Create peer discovery manager
	pdm := core.NewPeerDiscoveryManager()

	// Create timeout context (30 seconds)
	ctx, cancel := context.WithTimeout(discoveryCtx, 30*time.Second)
	defer cancel()

	// Get regions
	regions, err := pdm.GetAvailableRegions(ctx)
	if err != nil {
		log.Printf("[GetAvailableRegions] Error: %v", err)
		return nil, fmt.Errorf("failed to get regions: %w", err)
	}

	log.Printf("[GetAvailableRegions] Found %d regions", len(regions))
	return regions, nil
}

// CheckCustomPeers checks a list of user-provided peer URIs for availability
// Returns list of peers with their status (available/unavailable) and RTT
//
// Parameters:
//   - discoveryCtx: Context for peer discovery operations (can be cancelled)
//   - peerURIs: array of peer URIs (e.g., ["tls://host:port", "quic://host:port"])
//
// Returns: array of DiscoveredPeer with availability status and RTT
func CheckCustomPeers(discoveryCtx context.Context, peerURIs []string) ([]core.DiscoveredPeer, error) {
	if len(peerURIs) == 0 {
		return nil, fmt.Errorf("no peer URIs provided")
	}

	log.Printf("[CheckCustomPeers] Checking %d custom peers", len(peerURIs))

	// Create peer discovery manager
	pdm := core.NewPeerDiscoveryManager()

	// Create timeout context (60 seconds)
	ctx, cancel := context.WithTimeout(discoveryCtx, 60*time.Second)
	defer cancel()

	// Get batching parameters for concurrency
	_, concurrency, _ := core.GetBatchingParams()

	// Check peers
	results, err := pdm.CheckCustomPeers(ctx, peerURIs, concurrency)
	if err != nil {
		log.Printf("[CheckCustomPeers] Error: %v", err)
		return nil, fmt.Errorf("failed to check peers: %w", err)
	}

	// Count available peers
	available := 0
	for _, peer := range results {
		if peer.Available {
			available++
		}
	}

	log.Printf("[CheckCustomPeers] Results: %d/%d peers available", available, len(results))
	return results, nil
}

// AddDiscoveredPeer adds a discovered peer to the configuration
// This is a convenience method that converts DiscoveredPeer to PeerConfig
func AddDiscoveredPeer(cfg *core.Config, peer core.DiscoveredPeer) error {
	if cfg == nil {
		return fmt.Errorf("config not initialized")
	}

	log.Printf("[AddDiscoveredPeer] Adding peer: %s (RTT: %dms)", peer.Address, peer.RTT)

	// Add peer to config
	if err := cfg.AddPeer(peer.Address); err != nil {
		return fmt.Errorf("failed to add peer: %w", err)
	}

	// Save config
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	log.Printf("[AddDiscoveredPeer] Peer added successfully: %s", peer.Address)
	return nil
}

// AddDiscoveredPeers adds multiple discovered peers to the configuration
// This is a batch operation that adds all peers at once
func AddDiscoveredPeers(cfg *core.Config, peers []core.DiscoveredPeer) error {
	if cfg == nil {
		return fmt.Errorf("config not initialized")
	}

	if len(peers) == 0 {
		return fmt.Errorf("no peers provided")
	}

	log.Printf("[AddDiscoveredPeers] Adding %d peers", len(peers))

	// Add all peers
	added := 0
	for _, peer := range peers {
		if err := cfg.AddPeer(peer.Address); err != nil {
			log.Printf("[AddDiscoveredPeers] Warning: failed to add peer %s: %v", peer.Address, err)
			continue
		}
		added++
	}

	// Save config
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	log.Printf("[AddDiscoveredPeers] Successfully added %d/%d peers", added, len(peers))
	return nil
}

// GetPeerDiscoverySystemInfo returns system information for debugging
// Useful for troubleshooting peer discovery issues
func GetPeerDiscoverySystemInfo() map[string]interface{} {
	batchSize, concurrency, pauseMs := core.GetBatchingParams()

	info := map[string]interface{}{
		"system_info": core.GetSystemInfo(),
		"batching": map[string]interface{}{
			"batch_size":  batchSize,
			"concurrency": concurrency,
			"pause_ms":    pauseMs,
		},
		"cache_ttl_hours": core.CacheTTLHours,
	}

	log.Printf("[GetPeerDiscoverySystemInfo] %+v", info)
	return info
}
