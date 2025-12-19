/**
 * Config Store - Manages application configuration
 * Handles peers, UI preferences (theme, language), and autostart settings
 */

import { create } from 'zustand';
import {
  GetConfig,
  SaveConfig,
  AddPeer,
  RemovePeer,
  EnablePeer,
  DisablePeer,
  SetLanguage,
  SetTheme,
  SetAutoStart,
  SetPassword,
  GetDefaultPeers,
} from '../wailsjs/go/main/App';
import type { ConfigDTO, PeerConfigDTO } from '../wailsjs/go/main/models';

export type Theme = 'light' | 'dark' | 'system';
export type Language = 'en' | 'ru';

export interface ConfigState {
  // Configuration
  config: ConfigDTO | null;

  // Loading states
  isLoading: boolean;
  isSaving: boolean;

  // Actions
  loadConfig: () => Promise<void>;
  saveConfig: (config: ConfigDTO) => Promise<void>;

  // Peer management
  addPeer: (address: string) => Promise<void>;
  removePeer: (address: string) => Promise<void>;
  enablePeer: (address: string) => Promise<void>;
  disablePeer: (address: string) => Promise<void>;
  getDefaultPeers: () => Promise<string[]>;

  // UI preferences
  setLanguage: (language: Language) => Promise<void>;
  setTheme: (theme: Theme) => Promise<void>;
  setAutoStart: (enabled: boolean) => Promise<void>;

  // Security
  setPassword: (password: string) => Promise<void>;

  // Helpers
  getPeer: (address: string) => PeerConfigDTO | undefined;
  getEnabledPeers: () => PeerConfigDTO[];
}

export const useConfigStore = create<ConfigState>((set, get) => ({
  // Initial state
  config: null,
  isLoading: false,
  isSaving: false,

  // Load configuration from backend
  loadConfig: async () => {
    try {
      set({ isLoading: true });
      const config = await GetConfig();
      set({ config, isLoading: false });
    } catch (error) {
      console.error('Failed to load config:', error);
      set({ isLoading: false });
      throw error;
    }
  },

  // Save entire configuration
  saveConfig: async (config: ConfigDTO) => {
    try {
      set({ isSaving: true });
      await SaveConfig(config);
      set({ config, isSaving: false });
    } catch (error) {
      console.error('Failed to save config:', error);
      set({ isSaving: false });
      throw error;
    }
  },

  // Add a new peer
  addPeer: async (address: string) => {
    try {
      set({ isSaving: true });
      await AddPeer(address);
      await get().loadConfig(); // Reload config to get updated peer list
      set({ isSaving: false });
    } catch (error) {
      console.error('Failed to add peer:', error);
      set({ isSaving: false });
      throw error;
    }
  },

  // Remove a peer
  removePeer: async (address: string) => {
    try {
      set({ isSaving: true });
      await RemovePeer(address);
      await get().loadConfig();
      set({ isSaving: false });
    } catch (error) {
      console.error('Failed to remove peer:', error);
      set({ isSaving: false });
      throw error;
    }
  },

  // Enable a peer
  enablePeer: async (address: string) => {
    try {
      set({ isSaving: true });
      await EnablePeer(address);
      await get().loadConfig();
      set({ isSaving: false });
    } catch (error) {
      console.error('Failed to enable peer:', error);
      set({ isSaving: false });
      throw error;
    }
  },

  // Disable a peer
  disablePeer: async (address: string) => {
    try {
      set({ isSaving: true });
      await DisablePeer(address);
      await get().loadConfig();
      set({ isSaving: false });
    } catch (error) {
      console.error('Failed to disable peer:', error);
      set({ isSaving: false });
      throw error;
    }
  },

  // Get default recommended peers
  getDefaultPeers: async () => {
    try {
      const peers = await GetDefaultPeers();
      return peers;
    } catch (error) {
      console.error('Failed to get default peers:', error);
      return [];
    }
  },

  // Set UI language
  setLanguage: async (language: Language) => {
    try {
      set({ isSaving: true });
      await SetLanguage(language);
      await get().loadConfig();
      set({ isSaving: false });
    } catch (error) {
      console.error('Failed to set language:', error);
      set({ isSaving: false });
      throw error;
    }
  },

  // Set UI theme
  setTheme: async (theme: Theme) => {
    try {
      set({ isSaving: true });
      await SetTheme(theme);
      await get().loadConfig();
      set({ isSaving: false });
    } catch (error) {
      console.error('Failed to set theme:', error);
      set({ isSaving: false });
      throw error;
    }
  },

  // Set autostart preference
  setAutoStart: async (enabled: boolean) => {
    try {
      set({ isSaving: true });
      await SetAutoStart(enabled);
      await get().loadConfig();
      set({ isSaving: false });
    } catch (error) {
      console.error('Failed to set autostart:', error);
      set({ isSaving: false });
      throw error;
    }
  },

  // Set password (stored in OS keyring)
  setPassword: async (password: string) => {
    try {
      set({ isSaving: true });
      await SetPassword(password);
      set({ isSaving: false });
    } catch (error) {
      console.error('Failed to set password:', error);
      set({ isSaving: false });
      throw error;
    }
  },

  // Get a specific peer by address
  getPeer: (address: string) => {
    const { config } = get();
    if (!config) return undefined;
    return config.peers.find((p) => p.address === address);
  },

  // Get only enabled peers
  getEnabledPeers: () => {
    const { config } = get();
    if (!config) return [];
    return config.peers.filter((p) => p.enabled);
  },
}));
