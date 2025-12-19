/**
 * useThemeManager - Hook for managing theme with system preference detection
 * Handles theme switching, system theme detection, and syncs with config
 */

import { useEffect } from 'react';
import { useConfigStore } from '../store/configStore';
import { useUIStore } from '../store/uiStore';
// @ts-ignore - Will be available after Wails generates bindings
import { GetSystemTheme } from '../wailsjs/go/main/App';

/**
 * Detects system theme preference using browser API
 * Falls back to 'light' if detection fails
 */
function detectBrowserTheme(): 'light' | 'dark' {
  if (typeof window === 'undefined') return 'light';

  try {
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    return prefersDark ? 'dark' : 'light';
  } catch {
    return 'light';
  }
}

/**
 * Detects system theme using backend (Windows registry, etc.)
 * Falls back to browser detection if backend fails
 */
async function detectSystemTheme(): Promise<'light' | 'dark'> {
  try {
    const theme = await GetSystemTheme();
    if (theme === 'dark' || theme === 'light') {
      return theme;
    }
  } catch (error) {
    console.warn('Failed to get system theme from backend:', error);
  }

  // Fallback to browser detection
  return detectBrowserTheme();
}

/**
 * Resolves the actual theme based on theme preference and system theme
 */
async function resolveTheme(themePreference: string): Promise<'light' | 'dark'> {
  if (themePreference === 'light') return 'light';
  if (themePreference === 'dark') return 'dark';

  // For 'system' or any other value, detect system theme
  return await detectSystemTheme();
}

/**
 * Hook that manages theme state and applies it to the DOM
 * Automatically syncs with config store and detects system theme changes
 *
 * @example
 * ```tsx
 * function App() {
 *   useThemeManager();
 *   return <YourApp />;
 * }
 * ```
 */
export function useThemeManager() {
  const config = useConfigStore((state) => state.config);
  const setResolvedTheme = useUIStore((state) => state.setResolvedTheme);
  const resolvedTheme = useUIStore((state) => state.resolvedTheme);

  // Update resolved theme when config changes
  useEffect(() => {
    if (!config) return;

    const updateTheme = async () => {
      const resolved = await resolveTheme(config.theme);
      setResolvedTheme(resolved);
    };

    updateTheme();
  }, [config?.theme, setResolvedTheme]);

  // Listen for system theme changes (when theme preference is 'system')
  useEffect(() => {
    if (!config || config.theme !== 'system') return;

    // Create media query listener for system theme changes
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');

    const handleChange = (e: MediaQueryListEvent) => {
      const newTheme = e.matches ? 'dark' : 'light';
      setResolvedTheme(newTheme);
    };

    // Modern browsers
    if (mediaQuery.addEventListener) {
      mediaQuery.addEventListener('change', handleChange);
      return () => mediaQuery.removeEventListener('change', handleChange);
    }
    // Legacy browsers
    else if (mediaQuery.addListener) {
      mediaQuery.addListener(handleChange);
      return () => mediaQuery.removeListener(handleChange);
    }
  }, [config?.theme, setResolvedTheme]);

  // Apply theme to DOM
  useEffect(() => {
    const html = document.documentElement;

    if (resolvedTheme === 'dark') {
      html.classList.add('dark');
    } else {
      html.classList.remove('dark');
    }
  }, [resolvedTheme]);

  return {
    resolvedTheme,
    themePreference: config?.theme || 'system',
  };
}
