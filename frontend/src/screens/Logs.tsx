import { useState } from 'react';
import { motion } from 'framer-motion';
import { useNavigate } from 'react-router-dom';
import {
  Button,
  Input,
  GlassCard,
  HolographicBorder,
  Badge,
  LogViewer,
} from '../components';
import type { LogEntry } from '../components';
import { useLogsStore } from '../store/logsStore';
import { useI18n } from '../hooks/useI18n';

/**
 * Logs Screen - Real-time log viewer
 *
 * Features:
 * - Real-time log stream (EventsOn 'service:log')
 * - Level filter checkboxes (INFO, WARN, ERROR, DEBUG)
 * - Search input (filter by message)
 * - Virtualized list (react-window) for performance
 * - Auto-scroll logic (scroll to bottom only if already at bottom)
 * - Export logs button (download as .txt)
 */
export function Logs() {
  const { t } = useI18n();
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState('');
  const [levelFilters, setLevelFilters] = useState<Set<string>>(
    new Set(['INFO', 'WARN', 'ERROR', 'DEBUG'])
  );
  const [autoScroll, setAutoScroll] = useState(true);

  // Get logs from store
  const logs = useLogsStore((state) => state.logs);
  const clearLogs = useLogsStore((state) => state.clearLogs);
  const exportLogs = useLogsStore((state) => state.exportLogs);

  // Filter logs
  const filteredLogs = logs
    .filter((log) => levelFilters.has(log.level))
    .filter((log) => {
      if (!searchQuery.trim()) return true;
      const query = searchQuery.toLowerCase();
      return (
        log.message.toLowerCase().includes(query) ||
        log.tag.toLowerCase().includes(query)
      );
    });

  // Convert to LogEntry format
  const logEntries: LogEntry[] = filteredLogs.map((log) => ({
    timestamp: log.timestamp,
    level: log.level as 'INFO' | 'WARN' | 'ERROR' | 'DEBUG',
    message: `[${log.tag}] ${log.message}`,
  }));

  // Toggle level filter
  const toggleLevelFilter = (level: string) => {
    const newFilters = new Set(levelFilters);
    if (newFilters.has(level)) {
      newFilters.delete(level);
    } else {
      newFilters.add(level);
    }
    setLevelFilters(newFilters);
  };

  // Handle export logs
  const handleExportLogs = () => {
    const logsText = exportLogs();
    const blob = new Blob([logsText], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `tyr-logs-${new Date().toISOString()}.txt`;
    a.click();
    URL.revokeObjectURL(url);
  };

  // Handle clear logs
  const handleClearLogs = () => {
    if (window.confirm(t('logs.clearConfirm'))) {
      clearLogs();
    }
  };

  // Level counts
  const levelCounts = {
    INFO: logs.filter((log) => log.level === 'INFO').length,
    WARN: logs.filter((log) => log.level === 'WARN').length,
    ERROR: logs.filter((log) => log.level === 'ERROR').length,
    DEBUG: logs.filter((log) => log.level === 'DEBUG').length,
  };

  return (
    <div className="space-y-6 pb-6">
      {/* Header */}
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, type: 'spring' }}
        className="flex items-center justify-between"
      >
        <div className="space-y-2">
          <h1 className="text-5xl font-display font-black text-transparent bg-clip-text bg-iridescent bg-[length:200%_100%] animate-holographic-spin">
            {t('logs.serviceLogs')}
          </h1>
          <p className="text-lg text-white/70 font-futuristic">
            {t('logs.logEntries', { filtered: filteredLogs.length, total: logs.length })}
          </p>
        </div>
        <div className="flex gap-3">
          <Button variant="secondary" glow onClick={handleExportLogs} disabled={logs.length === 0}>
            {t('logs.exportLogs')}
          </Button>
          <Button variant="danger" onClick={handleClearLogs} disabled={logs.length === 0}>
            {t('logs.clearLogsButton')}
          </Button>
        </div>
      </motion.div>

      {/* Stats Cards */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <motion.div
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.3, delay: 0.1 }}
        >
          <HolographicBorder borderWidth={levelFilters.has('INFO') ? 2 : 1}>
            <button
              onClick={() => toggleLevelFilter('INFO')}
              className={`w-full p-4 rounded-xl transition-all ${
                levelFilters.has('INFO')
                  ? 'bg-space-blue/50 scale-105'
                  : 'bg-space-blue/20 hover:bg-space-blue/30'
              }`}
            >
              <div className="text-center space-y-2">
                <Badge variant="info" size="sm">
                  INFO
                </Badge>
                <p className="text-3xl font-bold text-white">{levelCounts.INFO}</p>
              </div>
            </button>
          </HolographicBorder>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.3, delay: 0.2 }}
        >
          <HolographicBorder borderWidth={levelFilters.has('WARN') ? 2 : 1}>
            <button
              onClick={() => toggleLevelFilter('WARN')}
              className={`w-full p-4 rounded-xl transition-all ${
                levelFilters.has('WARN')
                  ? 'bg-space-blue/50 scale-105'
                  : 'bg-space-blue/20 hover:bg-space-blue/30'
              }`}
            >
              <div className="text-center space-y-2">
                <Badge variant="warning" size="sm">
                  WARN
                </Badge>
                <p className="text-3xl font-bold text-white">{levelCounts.WARN}</p>
              </div>
            </button>
          </HolographicBorder>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.3, delay: 0.3 }}
        >
          <HolographicBorder borderWidth={levelFilters.has('ERROR') ? 2 : 1}>
            <button
              onClick={() => toggleLevelFilter('ERROR')}
              className={`w-full p-4 rounded-xl transition-all ${
                levelFilters.has('ERROR')
                  ? 'bg-space-blue/50 scale-105'
                  : 'bg-space-blue/20 hover:bg-space-blue/30'
              }`}
            >
              <div className="text-center space-y-2">
                <Badge variant="error" size="sm">
                  ERROR
                </Badge>
                <p className="text-3xl font-bold text-white">{levelCounts.ERROR}</p>
              </div>
            </button>
          </HolographicBorder>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.3, delay: 0.4 }}
        >
          <HolographicBorder borderWidth={levelFilters.has('DEBUG') ? 2 : 1}>
            <button
              onClick={() => toggleLevelFilter('DEBUG')}
              className={`w-full p-4 rounded-xl transition-all ${
                levelFilters.has('DEBUG')
                  ? 'bg-space-blue/50 scale-105'
                  : 'bg-space-blue/20 hover:bg-space-blue/30'
              }`}
            >
              <div className="text-center space-y-2">
                <Badge variant="default" size="sm">
                  DEBUG
                </Badge>
                <p className="text-3xl font-bold text-white">{levelCounts.DEBUG}</p>
              </div>
            </button>
          </HolographicBorder>
        </motion.div>
      </div>

      {/* Filters */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.5 }}
      >
        <GlassCard padding="md">
          <div className="flex flex-col md:flex-row gap-4 items-center">
            <div className="flex-1 w-full">
              <Input
                placeholder={t('logs.searchPlaceholder')}
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                leftIcon={
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                    />
                  </svg>
                }
              />
            </div>
            <div className="flex items-center gap-3">
              <span className="text-sm text-white/70 font-futuristic">{t('logs.autoScrollLabel')}</span>
              <button
                onClick={() => setAutoScroll(!autoScroll)}
                className={`relative w-14 h-8 rounded-full transition-colors ${
                  autoScroll ? 'bg-neon-green' : 'bg-white/20'
                }`}
              >
                <motion.div
                  animate={{ x: autoScroll ? 24 : 2 }}
                  transition={{ type: 'spring', stiffness: 500, damping: 30 }}
                  className="absolute top-1 w-6 h-6 bg-white rounded-full shadow-lg"
                />
              </button>
            </div>
          </div>
        </GlassCard>
      </motion.div>

      {/* Log Viewer */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.6 }}
      >
        <GlassCard title={t('logs.logStream')} subtitle={t('logs.logStreamSubtitle')} padding="lg">
          {logEntries.length > 0 ? (
            <LogViewer
              logs={logEntries}
              maxHeight="600px"
              showFilters={false}
              autoScroll={autoScroll}
              onExport={handleExportLogs}
            />
          ) : (
            <div className="text-center py-16 space-y-4">
              <div className="text-8xl mb-4">📋</div>
              <p className="text-xl text-white/70">
                {logs.length === 0 ? t('logs.noLogsAvailable') : t('logs.noLogsMatchFilters')}
              </p>
              <p className="text-sm text-white/50">
                {logs.length === 0
                  ? t('logs.logsWillAppear')
                  : t('logs.adjustFilters')}
              </p>
              {logs.length === 0 && (
                <div className="pt-4">
                  <Button variant="primary" glow onClick={() => navigate('/')}>
                    {t('logs.goToDashboard')}
                  </Button>
                </div>
              )}
            </div>
          )}
        </GlassCard>
      </motion.div>

      {/* Info Box */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.7 }}
      >
        <HolographicBorder borderWidth={1}>
          <div className="bg-neon-cyan/10 border border-neon-cyan/30 rounded-lg p-6">
            <div className="flex items-start gap-4">
              <div className="text-4xl">💡</div>
              <div className="flex-1 space-y-2">
                <h3 className="text-white font-semibold text-lg">{t('logs.aboutLogs')}</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm text-white/80">
                  <div>
                    <p className="font-medium mb-2">{t('logs.logLevels')}</p>
                    <ul className="space-y-1 text-white/70">
                      <li>
                        <Badge variant="info" size="sm" className="mr-2">
                          INFO
                        </Badge>
                        {t('logs.info')}
                      </li>
                      <li>
                        <Badge variant="warning" size="sm" className="mr-2">
                          WARN
                        </Badge>
                        {t('logs.warn')}
                      </li>
                      <li>
                        <Badge variant="error" size="sm" className="mr-2">
                          ERROR
                        </Badge>
                        {t('logs.error')}
                      </li>
                      <li>
                        <Badge variant="default" size="sm" className="mr-2">
                          DEBUG
                        </Badge>
                        {t('logs.debug')}
                      </li>
                    </ul>
                  </div>
                  <div>
                    <p className="font-medium mb-2">{t('logs.features')}</p>
                    <ul className="list-disc list-inside space-y-1 text-white/70">
                      <li>{t('logs.feature1')}</li>
                      <li>{t('logs.feature2')}</li>
                      <li>{t('logs.feature3')}</li>
                      <li>{t('logs.feature4')}</li>
                      <li>{t('logs.feature5')}</li>
                    </ul>
                  </div>
                </div>
                <div className="pt-3 mt-3 border-t border-white/10">
                  <p className="text-xs text-white/60">
                    {t('logs.notice')}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </HolographicBorder>
      </motion.div>
    </div>
  );
}

export default Logs;
