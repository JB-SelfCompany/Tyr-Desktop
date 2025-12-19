import React, { useState, useEffect, useRef, useMemo } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Badge } from '../ui/Badge';
import { Input } from '../ui/Input';
import { Button } from '../ui/Button';

export type LogLevel = 'INFO' | 'WARN' | 'ERROR' | 'DEBUG';

export interface LogEntry {
  timestamp: string;
  level: LogLevel;
  message: string;
}

interface LogViewerProps {
  logs: LogEntry[];
  maxHeight?: string;
  showFilters?: boolean;
  autoScroll?: boolean;
  onExport?: () => void;
}

const levelColors: Record<LogLevel, string> = {
  INFO: 'text-md-light-tertiary dark:text-md-dark-tertiary',
  WARN: 'text-md-light-onErrorContainer dark:text-md-dark-onErrorContainer',
  ERROR: 'text-md-light-error dark:text-md-dark-error',
  DEBUG: 'text-md-light-outline dark:text-md-dark-outline',
};

const levelBadgeVariant: Record<LogLevel, 'info' | 'warning' | 'error' | 'default'> = {
  INFO: 'info',
  WARN: 'warning',
  ERROR: 'error',
  DEBUG: 'default',
};

export const LogViewer: React.FC<LogViewerProps> = ({
  logs,
  maxHeight = '600px',
  showFilters = true,
  autoScroll = true,
  onExport,
}) => {
  const [searchTerm, setSearchTerm] = useState('');
  const [levelFilters, setLevelFilters] = useState<Set<LogLevel>>(
    new Set(['INFO', 'WARN', 'ERROR', 'DEBUG'])
  );
  const [isAtBottom, setIsAtBottom] = useState(true);
  const scrollContainerRef = useRef<HTMLDivElement>(null);
  const endRef = useRef<HTMLDivElement>(null);

  // Filter logs with useMemo for performance
  const filteredLogs = useMemo(() => {
    return logs.filter((log) => {
      const matchesSearch = log.message.toLowerCase().includes(searchTerm.toLowerCase());
      const matchesLevel = levelFilters.has(log.level);
      return matchesSearch && matchesLevel;
    });
  }, [logs, searchTerm, levelFilters]);

  // Toggle level filter
  const toggleLevelFilter = (level: LogLevel) => {
    setLevelFilters((prev) => {
      const newFilters = new Set(prev);
      if (newFilters.has(level)) {
        newFilters.delete(level);
      } else {
        newFilters.add(level);
      }
      return newFilters;
    });
  };

  // Auto-scroll logic - scroll to bottom when auto-scroll is enabled
  useEffect(() => {
    if (autoScroll && scrollContainerRef.current) {
      // Scroll the container to bottom instead of using scrollIntoView on endRef
      scrollContainerRef.current.scrollTop = scrollContainerRef.current.scrollHeight;
    }
  }, [logs, autoScroll]);

  // Track scroll position
  const handleScroll = () => {
    if (!scrollContainerRef.current) return;

    const { scrollTop, scrollHeight, clientHeight } = scrollContainerRef.current;
    const atBottom = scrollHeight - scrollTop - clientHeight < 50;
    setIsAtBottom(atBottom);
  };

  // Scroll to bottom button
  const scrollToBottom = () => {
    if (scrollContainerRef.current) {
      scrollContainerRef.current.scrollTo({
        top: scrollContainerRef.current.scrollHeight,
        behavior: 'smooth'
      });
      setIsAtBottom(true);
    }
  };

  return (
    <div className="flex flex-col gap-4">
      {/* Filters */}
      {showFilters && (
        <div className="flex flex-col sm:flex-row gap-3">
          {/* Search */}
          <div className="flex-1">
            <Input
              placeholder="Search logs..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
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

          {/* Level filters */}
          <div className="flex gap-2 flex-wrap">
            {(['INFO', 'WARN', 'ERROR', 'DEBUG'] as LogLevel[]).map((level) => (
              <Button
                key={level}
                variant={levelFilters.has(level) ? 'primary' : 'ghost'}
                size="sm"
                onClick={() => toggleLevelFilter(level)}
              >
                {level}
              </Button>
            ))}
          </div>

          {/* Export button */}
          {onExport && (
            <Button variant="secondary" size="sm" onClick={onExport}>
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"
                />
              </svg>
              Export
            </Button>
          )}
        </div>
      )}

      {/* Log container */}
      <div className="relative">
        <div
          ref={scrollContainerRef}
          className="bg-md-light-surfaceVariant/50 dark:bg-md-dark-surfaceVariant/50 backdrop-blur-lg border border-md-light-outline/30 dark:border-md-dark-outline/30 rounded-xl p-4 font-mono text-sm overflow-y-auto"
          style={{ maxHeight }}
          onScroll={handleScroll}
        >
          <div className="space-y-2">
            <AnimatePresence initial={false}>
              {filteredLogs.length === 0 ? (
                <motion.div
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  exit={{ opacity: 0 }}
                  className="text-center text-md-light-outline dark:text-md-dark-outline py-8"
                >
                  No logs to display
                </motion.div>
              ) : (
                filteredLogs.map((log, index) => (
                  <motion.div
                    key={`${log.timestamp}-${index}`}
                    initial={{ opacity: 0, x: -20 }}
                    animate={{ opacity: 1, x: 0 }}
                    exit={{ opacity: 0, x: 20 }}
                    transition={{ type: 'spring', stiffness: 300, damping: 30 }}
                    className="flex items-start gap-3 p-2 rounded-lg hover:bg-md-light-primaryContainer/20 dark:hover:bg-md-dark-primaryContainer/10 transition-colors"
                  >
                    {/* Timestamp */}
                    <span className="text-md-light-outline dark:text-md-dark-outline text-xs whitespace-nowrap flex-shrink-0">
                      {new Date(log.timestamp).toLocaleTimeString()}
                    </span>

                    {/* Level badge */}
                    <Badge variant={levelBadgeVariant[log.level]} size="sm">
                      {log.level}
                    </Badge>

                    {/* Message */}
                    <span className={`flex-1 ${levelColors[log.level]} break-words`}>
                      {log.message}
                    </span>
                  </motion.div>
                ))
              )}
            </AnimatePresence>
            <div ref={endRef} />
          </div>
        </div>

        {/* Scroll to bottom button */}
        {!isAtBottom && (
          <motion.button
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: 10 }}
            onClick={scrollToBottom}
            className="absolute bottom-4 right-4 p-2 rounded-full bg-md-light-primary dark:bg-md-dark-primary text-md-light-onPrimary dark:text-md-dark-onPrimary shadow-lg hover:bg-md-light-primaryContainer dark:hover:bg-md-dark-primaryContainer transition-colors"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M19 14l-7 7m0 0l-7-7m7 7V3"
              />
            </svg>
          </motion.button>
        )}
      </div>

      {/* Stats */}
      <div className="flex gap-4 text-sm text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant font-body">
        <span>Total: {logs.length}</span>
        <span>Filtered: {filteredLogs.length}</span>
      </div>
    </div>
  );
};
