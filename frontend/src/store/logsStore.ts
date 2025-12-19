/**
 * Logs Store - Manages real-time logs with filtering
 * Handles log events from backend with level and search filtering
 */

import { create } from 'zustand';
import type { LogEventDTO } from '../wailsjs/go/main/models';

export type LogLevel = 'INFO' | 'WARN' | 'ERROR' | 'DEBUG' | 'ALL';

export interface LogsState {
  // All logs (unfiltered)
  logs: LogEventDTO[];

  // Filtering
  levelFilter: Set<LogLevel>;
  searchFilter: string;

  // UI state
  maxLogs: number;
  autoScroll: boolean;
  isPaused: boolean;

  // Actions
  addLog: (log: LogEventDTO) => void;
  clearLogs: () => void;
  setLevelFilter: (levels: LogLevel[]) => void;
  toggleLevel: (level: LogLevel) => void;
  setSearchFilter: (search: string) => void;
  setAutoScroll: (enabled: boolean) => void;
  togglePause: () => void;
  setMaxLogs: (max: number) => void;
  exportLogs: () => string;

  // Computed
  getFilteredLogs: () => LogEventDTO[];
}

const DEFAULT_MAX_LOGS = 1000;

export const useLogsStore = create<LogsState>((set, get) => ({
  // Initial state
  logs: [],
  levelFilter: new Set<LogLevel>(['INFO', 'WARN', 'ERROR', 'DEBUG']),
  searchFilter: '',
  maxLogs: DEFAULT_MAX_LOGS,
  autoScroll: true,
  isPaused: false,

  // Add a new log entry (from event stream)
  addLog: (log: LogEventDTO) => {
    const { isPaused, maxLogs } = get();

    // Skip if paused
    if (isPaused) return;

    set((state) => {
      const newLogs = [...state.logs, log];

      // Keep only the last maxLogs entries
      if (newLogs.length > maxLogs) {
        return { logs: newLogs.slice(-maxLogs) };
      }

      return { logs: newLogs };
    });
  },

  // Clear all logs
  clearLogs: () => {
    set({ logs: [] });
  },

  // Set level filter (replaces all)
  setLevelFilter: (levels: LogLevel[]) => {
    set({ levelFilter: new Set(levels) });
  },

  // Toggle a specific level
  toggleLevel: (level: LogLevel) => {
    set((state) => {
      const newFilter = new Set(state.levelFilter);
      if (newFilter.has(level)) {
        newFilter.delete(level);
      } else {
        newFilter.add(level);
      }
      return { levelFilter: newFilter };
    });
  },

  // Set search filter
  setSearchFilter: (search: string) => {
    set({ searchFilter: search.toLowerCase() });
  },

  // Set auto-scroll behavior
  setAutoScroll: (enabled: boolean) => {
    set({ autoScroll: enabled });
  },

  // Toggle pause state
  togglePause: () => {
    set((state) => ({ isPaused: !state.isPaused }));
  },

  // Set maximum number of logs to keep
  setMaxLogs: (max: number) => {
    set({ maxLogs: max });
    const { logs } = get();
    if (logs.length > max) {
      set({ logs: logs.slice(-max) });
    }
  },

  // Export logs as text
  exportLogs: () => {
    const logs = get().getFilteredLogs();
    return logs
      .map((log) => {
        const timestamp = new Date(log.timestamp).toISOString();
        return `[${timestamp}] [${log.level}] [${log.tag}] ${log.message}`;
      })
      .join('\n');
  },

  // Get filtered logs based on level and search
  getFilteredLogs: () => {
    const { logs, levelFilter, searchFilter } = get();

    return logs.filter((log) => {
      // Filter by level
      if (!levelFilter.has(log.level as LogLevel)) {
        return false;
      }

      // Filter by search term
      if (searchFilter) {
        const searchableText = `${log.level} ${log.tag} ${log.message}`.toLowerCase();
        if (!searchableText.includes(searchFilter)) {
          return false;
        }
      }

      return true;
    });
  },
}));
