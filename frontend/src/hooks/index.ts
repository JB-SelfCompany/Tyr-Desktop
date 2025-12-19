/**
 * Hooks Index - Centralized exports for all custom hooks
 * Provides convenient single import point for all hooks
 */

export {
  useServiceStatus,
  useServiceRunning,
  usePeerStats,
  useServiceActions,
} from './useServiceStatus';
export type { UseServiceStatusOptions } from './useServiceStatus';

export {
  useEventStream,
  useEventStreamWithHandlers,
  useLogEvents,
  useMailEvents,
  useConnectionEvents,
  useServiceStatusEvents,
  usePeerStatsEvents,
  EventNames,
} from './useEventStream';

export {
  useConfig,
  usePeers,
  useUIPreferences,
  useTheme,
  useLanguage,
  useAutoStart,
} from './useConfig';
export type { UseConfigOptions } from './useConfig';

export { useI18n, useTranslate, useCurrentLanguage } from './useI18n';

export { useThemeManager } from './useThemeManager';
