/**
 * useConfig - Hook for managing application configuration
 * Provides convenient access to config store with auto-loading
 */

import { useEffect } from 'react';
import { useConfigStore, type Theme, type Language } from '../store/configStore';

export interface UseConfigOptions {
  /**
   * Whether to load config automatically on mount
   * @default true
   */
  loadOnMount?: boolean;
}

/**
 * Hook that provides access to configuration with auto-loading
 *
 * @param options - Configuration options
 * @returns Config store state and actions
 *
 * @example
 * ```tsx
 * const { config, setTheme, setLanguage } = useConfig();
 *
 * if (config) {
 *   console.log('Current theme:', config.theme);
 *   console.log('Current language:', config.language);
 * }
 * ```
 */
export function useConfig(options: UseConfigOptions = {}) {
  const { loadOnMount = true } = options;

  const config = useConfigStore((state) => state.config);
  const loadConfig = useConfigStore((state) => state.loadConfig);
  const isLoading = useConfigStore((state) => state.isLoading);

  useEffect(() => {
    if (loadOnMount && !config) {
      loadConfig();
    }
  }, [loadOnMount, config, loadConfig]);

  const store = useConfigStore.getState();
  return {
    ...store,
    config,
    isLoading,
  };
}

/**
 * Hook that returns only peer-related actions
 * Useful when you only need to manage peers
 *
 * @example
 * ```tsx
 * const { peers, addPeer, removePeer } = usePeers();
 * ```
 */
export function usePeers() {
  const config = useConfigStore((state) => state.config);
  const addPeer = useConfigStore((state) => state.addPeer);
  const removePeer = useConfigStore((state) => state.removePeer);
  const enablePeer = useConfigStore((state) => state.enablePeer);
  const disablePeer = useConfigStore((state) => state.disablePeer);
  const getDefaultPeers = useConfigStore((state) => state.getDefaultPeers);
  const getPeer = useConfigStore((state) => state.getPeer);
  const getEnabledPeers = useConfigStore((state) => state.getEnabledPeers);
  const isSaving = useConfigStore((state) => state.isSaving);

  return {
    peers: config?.peers || [],
    addPeer,
    removePeer,
    enablePeer,
    disablePeer,
    getDefaultPeers,
    getPeer,
    getEnabledPeers,
    isSaving,
  };
}

/**
 * Hook that returns only UI preference actions
 * Useful when you only need theme/language controls
 *
 * @example
 * ```tsx
 * const { theme, language, setTheme, setLanguage } = useUIPreferences();
 * ```
 */
export function useUIPreferences() {
  const config = useConfigStore((state) => state.config);
  const setTheme = useConfigStore((state) => state.setTheme);
  const setLanguage = useConfigStore((state) => state.setLanguage);
  const setAutoStart = useConfigStore((state) => state.setAutoStart);
  const isSaving = useConfigStore((state) => state.isSaving);

  return {
    theme: (config?.theme || 'system') as Theme,
    language: (config?.language || 'en') as Language,
    autoStart: config?.autoStart || false,
    setTheme,
    setLanguage,
    setAutoStart,
    isSaving,
  };
}

/**
 * Hook that returns current theme
 * Useful for theme-aware components
 *
 * @example
 * ```tsx
 * const theme = useTheme();
 * ```
 */
export function useTheme() {
  const config = useConfigStore((state) => state.config);
  return (config?.theme || 'system') as Theme;
}

/**
 * Hook that returns current language
 * Useful for language-aware components
 *
 * @example
 * ```tsx
 * const language = useLanguage();
 * ```
 */
export function useLanguage() {
  const config = useConfigStore((state) => state.config);
  return (config?.language || 'en') as Language;
}

/**
 * Hook that returns autostart setting
 *
 * @example
 * ```tsx
 * const { autoStart, setAutoStart } = useAutoStart();
 * ```
 */
export function useAutoStart() {
  const config = useConfigStore((state) => state.config);
  const setAutoStart = useConfigStore((state) => state.setAutoStart);

  return {
    autoStart: config?.autoStart || false,
    setAutoStart,
  };
}
