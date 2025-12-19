/**
 * useServiceStatus - Hook for auto-refreshing service status
 * Automatically fetches service status and peer statistics at regular intervals
 */

import { useEffect, useRef } from 'react';
import { useServiceStore } from '../store/serviceStore';

export interface UseServiceStatusOptions {
  /**
   * Auto-refresh interval in milliseconds
   * Set to 0 to disable auto-refresh
   * @default 5000 (5 seconds)
   */
  refreshInterval?: number;

  /**
   * Whether to fetch immediately on mount
   * @default true
   */
  fetchOnMount?: boolean;

  /**
   * Whether to fetch peer stats along with service status
   * @default true
   */
  includePeerStats?: boolean;
}

/**
 * Hook that automatically fetches and updates service status
 *
 * @param options - Configuration options
 * @returns Service store state and actions
 *
 * @example
 * ```tsx
 * const { status, running, startService } = useServiceStatus({
 *   refreshInterval: 3000, // Refresh every 3 seconds
 * });
 * ```
 */
export function useServiceStatus(options: UseServiceStatusOptions = {}) {
  const {
    refreshInterval = 5000,
    fetchOnMount = true,
    includePeerStats = true,
  } = options;

  const fetchStatus = useServiceStore((state) => state.fetchStatus);
  const fetchPeerStats = useServiceStore((state) => state.fetchPeerStats);
  const running = useServiceStore((state) => state.running);

  const intervalRef = useRef<number | null>(null);

  useEffect(() => {
    // Fetch on mount
    if (fetchOnMount) {
      fetchStatus();
      if (includePeerStats && running) {
        fetchPeerStats();
      }
    }

    // Set up auto-refresh interval
    if (refreshInterval > 0) {
      intervalRef.current = window.setInterval(() => {
        fetchStatus();
        if (includePeerStats && useServiceStore.getState().running) {
          fetchPeerStats();
        }
      }, refreshInterval);
    }

    // Cleanup
    return () => {
      if (intervalRef.current !== null) {
        clearInterval(intervalRef.current);
        intervalRef.current = null;
      }
    };
  }, [refreshInterval, fetchOnMount, includePeerStats]);

  // Return the entire service store state and actions
  return useServiceStore();
}

/**
 * Hook that returns only the service running status
 * Useful when you only need to check if service is running
 *
 * @example
 * ```tsx
 * const isRunning = useServiceRunning();
 * ```
 */
export function useServiceRunning() {
  return useServiceStore((state) => state.running);
}

/**
 * Hook that returns peer statistics
 *
 * @example
 * ```tsx
 * const peers = usePeerStats();
 * ```
 */
export function usePeerStats() {
  return useServiceStore((state) => state.peers);
}

/**
 * Hook that returns service actions only
 * Useful when you don't need the state, only actions
 *
 * @example
 * ```tsx
 * const { startService, stopService, restartService } = useServiceActions();
 * ```
 */
export function useServiceActions() {
  return useServiceStore((state) => ({
    startService: state.startService,
    stopService: state.stopService,
    restartService: state.restartService,
    initializeService: state.initializeService,
    hotReloadPeers: state.hotReloadPeers,
    refreshAll: state.refreshAll,
  }));
}
