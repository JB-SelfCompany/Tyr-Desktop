/**
 * Service Store - Manages service lifecycle and peer statistics
 * Uses Zustand for state management with TypeScript support
 */

import { create } from 'zustand';
import {
  GetServiceStatus,
  StartService,
  StopService,
  RestartService,
  InitializeService,
  GetPeerStats,
  HotReloadPeers,
  GetMailAddress,
  IsServiceRunning,
} from '../wailsjs/go/main/App';
import type { ServiceStatusDTO, PeerInfoDTO } from '../wailsjs/go/main/models';

export interface ServiceState {
  // Service status
  status: string;
  running: boolean;
  mailAddress: string;
  smtpAddress: string;
  imapAddress: string;
  databasePath: string;
  errorMessage: string;

  // Peer statistics
  peers: PeerInfoDTO[];

  // Loading states
  isLoading: boolean;
  isStarting: boolean;
  isStopping: boolean;
  isRestarting: boolean;
  isHotReloading: boolean;

  // Last update timestamp
  lastUpdate: number;

  // Actions
  fetchStatus: () => Promise<void>;
  fetchPeerStats: () => Promise<void>;
  startService: () => Promise<void>;
  stopService: () => Promise<void>;
  restartService: () => Promise<void>;
  initializeService: () => Promise<void>;
  hotReloadPeers: () => Promise<void>;
  refreshAll: () => Promise<void>;

  // State setters (for event updates)
  setStatus: (status: ServiceStatusDTO) => void;
  setPeers: (peers: PeerInfoDTO[]) => void;
  setMailAddress: (address: string) => void;
}

export const useServiceStore = create<ServiceState>((set, get) => ({
  // Initial state
  status: 'Stopped',
  running: false,
  mailAddress: '',
  smtpAddress: '',
  imapAddress: '',
  databasePath: '',
  errorMessage: '',
  peers: [],
  isLoading: false,
  isStarting: false,
  isStopping: false,
  isRestarting: false,
  isHotReloading: false,
  lastUpdate: 0,

  // Fetch service status from backend
  fetchStatus: async () => {
    try {
      set({ isLoading: true });
      const status = await GetServiceStatus();
      set({
        status: status.status,
        running: status.running,
        mailAddress: status.mailAddress,
        smtpAddress: status.smtpAddress,
        imapAddress: status.imapAddress,
        databasePath: status.databasePath,
        errorMessage: status.errorMessage || '',
        lastUpdate: Date.now(),
        isLoading: false,
      });
    } catch (error) {
      console.error('Failed to fetch service status:', error);
      set({
        status: 'Error',
        running: false,
        errorMessage: error instanceof Error ? error.message : 'Unknown error',
        isLoading: false,
      });
    }
  },

  // Fetch peer statistics from backend
  fetchPeerStats: async () => {
    try {
      const peers = await GetPeerStats();
      set({ peers: peers || [], lastUpdate: Date.now() });
    } catch (error) {
      console.error('Failed to fetch peer stats:', error);
      set({ peers: [] });
    }
  },

  // Initialize the service (must be called before starting)
  initializeService: async () => {
    try {
      set({ isLoading: true });
      await InitializeService();
      await get().fetchStatus();
      set({ isLoading: false });
    } catch (error) {
      console.error('Failed to initialize service:', error);
      set({
        status: 'Error',
        running: false,
        errorMessage: error instanceof Error ? error.message : 'Failed to initialize service',
        isLoading: false,
      });
      throw error;
    }
  },

  // Start the service
  startService: async () => {
    try {
      set({ isStarting: true, errorMessage: '' });
      await StartService();
      // Wait a bit for service to start
      await new Promise((resolve) => setTimeout(resolve, 500));
      await get().fetchStatus();
      await get().fetchPeerStats();
      set({ isStarting: false });
    } catch (error) {
      console.error('Failed to start service:', error);
      set({
        isStarting: false,
        errorMessage: error instanceof Error ? error.message : 'Failed to start service',
      });
      throw error;
    }
  },

  // Stop the service
  stopService: async () => {
    try {
      set({ isStopping: true, errorMessage: '' });
      await StopService();
      await new Promise((resolve) => setTimeout(resolve, 500));
      await get().fetchStatus();
      set({ isStopping: false, peers: [] });
    } catch (error) {
      console.error('Failed to stop service:', error);
      set({
        isStopping: false,
        errorMessage: error instanceof Error ? error.message : 'Failed to stop service',
      });
      throw error;
    }
  },

  // Restart the service
  restartService: async () => {
    try {
      set({ isRestarting: true, errorMessage: '' });
      await RestartService();
      await new Promise((resolve) => setTimeout(resolve, 1000));
      await get().fetchStatus();
      await get().fetchPeerStats();
      set({ isRestarting: false });
    } catch (error) {
      console.error('Failed to restart service:', error);
      set({
        isRestarting: false,
        errorMessage: error instanceof Error ? error.message : 'Failed to restart service',
      });
      throw error;
    }
  },

  // Hot-reload peers without restarting service
  hotReloadPeers: async () => {
    try {
      set({ isHotReloading: true, errorMessage: '' });
      await HotReloadPeers();
      await new Promise((resolve) => setTimeout(resolve, 500));
      await get().fetchPeerStats();
      set({ isHotReloading: false });
    } catch (error) {
      console.error('Failed to hot-reload peers:', error);
      set({
        isHotReloading: false,
        errorMessage: error instanceof Error ? error.message : 'Failed to hot-reload peers',
      });
      throw error;
    }
  },

  // Refresh all service data
  refreshAll: async () => {
    await get().fetchStatus();
    const { running } = get();
    if (running) {
      await get().fetchPeerStats();
    }
  },

  // Set status from event (real-time update)
  setStatus: (status: ServiceStatusDTO) => {
    set({
      status: status.status,
      running: status.running,
      mailAddress: status.mailAddress,
      smtpAddress: status.smtpAddress,
      imapAddress: status.imapAddress,
      databasePath: status.databasePath,
      errorMessage: status.errorMessage || '',
      lastUpdate: Date.now(),
    });
  },

  // Set peers from event (real-time update)
  setPeers: (peers: PeerInfoDTO[]) => {
    set({ peers: peers || [], lastUpdate: Date.now() });
  },

  // Set mail address from event
  setMailAddress: (address: string) => {
    set({ mailAddress: address, lastUpdate: Date.now() });
  },
}));
