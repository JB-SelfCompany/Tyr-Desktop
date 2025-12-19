/**
 * UI Store - Manages UI state (theme, modals, notifications, loading)
 * Handles global UI state that doesn't belong to specific domains
 */

import { create } from 'zustand';

export type Theme = 'light' | 'dark' | 'system';
export type NotificationType = 'success' | 'error' | 'warning' | 'info';

export interface Notification {
  id: string;
  type: NotificationType;
  title: string;
  message: string;
  duration?: number; // ms, undefined = persistent
  timestamp: number;
}

export interface Modal {
  id: string;
  component: string; // Modal component name
  props?: Record<string, any>;
}

export interface UIState {
  // Theme
  theme: Theme;
  resolvedTheme: 'light' | 'dark'; // Actual theme after system preference resolution

  // Modals
  modals: Modal[];

  // Notifications (toast notifications)
  notifications: Notification[];

  // Global loading states
  isAppLoading: boolean;
  globalLoading: boolean;

  // Sidebar state (if applicable)
  sidebarCollapsed: boolean;

  // Actions - Theme
  setTheme: (theme: Theme) => void;
  setResolvedTheme: (theme: 'light' | 'dark') => void;

  // Actions - Modals
  openModal: (component: string, props?: Record<string, any>) => string;
  closeModal: (id: string) => void;
  closeAllModals: () => void;
  getTopModal: () => Modal | undefined;

  // Actions - Notifications
  showNotification: (
    type: NotificationType,
    title: string,
    message: string,
    duration?: number
  ) => string;
  dismissNotification: (id: string) => void;
  clearAllNotifications: () => void;

  // Actions - Loading
  setAppLoading: (loading: boolean) => void;
  setGlobalLoading: (loading: boolean) => void;

  // Actions - Sidebar
  toggleSidebar: () => void;
  setSidebarCollapsed: (collapsed: boolean) => void;
}

let notificationIdCounter = 0;
let modalIdCounter = 0;

export const useUIStore = create<UIState>((set, get) => ({
  // Initial state
  theme: 'system',
  resolvedTheme: 'dark', // Default to dark
  modals: [],
  notifications: [],
  isAppLoading: true,
  globalLoading: false,
  sidebarCollapsed: false,

  // Set theme preference
  setTheme: (theme: Theme) => {
    set({ theme });
  },

  // Set resolved theme (after system preference detection)
  setResolvedTheme: (theme: 'light' | 'dark') => {
    set({ resolvedTheme: theme });
  },

  // Open a modal
  openModal: (component: string, props?: Record<string, any>) => {
    const id = `modal-${++modalIdCounter}`;
    const modal: Modal = { id, component, props };
    set((state) => ({ modals: [...state.modals, modal] }));
    return id;
  },

  // Close a specific modal
  closeModal: (id: string) => {
    set((state) => ({
      modals: state.modals.filter((m) => m.id !== id),
    }));
  },

  // Close all modals
  closeAllModals: () => {
    set({ modals: [] });
  },

  // Get the top-most modal
  getTopModal: () => {
    const { modals } = get();
    return modals.length > 0 ? modals[modals.length - 1] : undefined;
  },

  // Show a notification
  showNotification: (
    type: NotificationType,
    title: string,
    message: string,
    duration = 5000
  ) => {
    const id = `notification-${++notificationIdCounter}`;
    const notification: Notification = {
      id,
      type,
      title,
      message,
      duration,
      timestamp: Date.now(),
    };

    set((state) => ({
      notifications: [...state.notifications, notification],
    }));

    // Auto-dismiss after duration
    if (duration) {
      setTimeout(() => {
        get().dismissNotification(id);
      }, duration);
    }

    return id;
  },

  // Dismiss a notification
  dismissNotification: (id: string) => {
    set((state) => ({
      notifications: state.notifications.filter((n) => n.id !== id),
    }));
  },

  // Clear all notifications
  clearAllNotifications: () => {
    set({ notifications: [] });
  },

  // Set app loading state (initial load)
  setAppLoading: (loading: boolean) => {
    set({ isAppLoading: loading });
  },

  // Set global loading state (for operations)
  setGlobalLoading: (loading: boolean) => {
    set({ globalLoading: loading });
  },

  // Toggle sidebar collapsed state
  toggleSidebar: () => {
    set((state) => ({ sidebarCollapsed: !state.sidebarCollapsed }));
  },

  // Set sidebar collapsed state
  setSidebarCollapsed: (collapsed: boolean) => {
    set({ sidebarCollapsed: collapsed });
  },
}));

// Helper functions for common notifications
export const showSuccess = (title: string, message: string, duration?: number) =>
  useUIStore.getState().showNotification('success', title, message, duration);

export const showError = (title: string, message: string, duration?: number) =>
  useUIStore.getState().showNotification('error', title, message, duration);

export const showWarning = (title: string, message: string, duration?: number) =>
  useUIStore.getState().showNotification('warning', title, message, duration);

export const showInfo = (title: string, message: string, duration?: number) =>
  useUIStore.getState().showNotification('info', title, message, duration);
