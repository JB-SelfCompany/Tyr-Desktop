import React from 'react';
import { motion } from 'framer-motion';
import { Badge, BadgeVariant } from '../ui/Badge';
import { useI18n } from '../../hooks/useI18n';

export type ServiceStatus = 'Running' | 'Stopped' | 'Starting' | 'Stopping' | 'Error';

interface StatusIndicatorProps {
  status: ServiceStatus;
  animated?: boolean;
  size?: 'sm' | 'md' | 'lg';
  showText?: boolean;
  className?: string;
}

const statusConfig: Record<ServiceStatus, {
  variant: BadgeVariant;
  icon: string;
  textKey: string;
  color: string;
}> = {
  Running: {
    variant: 'success',
    icon: '●',
    textKey: 'dashboard.status.running',
    color: '#22c55e', // green-500
  },
  Stopped: {
    variant: 'default',
    icon: '○',
    textKey: 'dashboard.status.stopped',
    color: '#6b7280', // gray-500
  },
  Starting: {
    variant: 'warning',
    icon: '◐',
    textKey: 'dashboard.status.starting',
    color: '#f59e0b', // amber-500
  },
  Stopping: {
    variant: 'warning',
    icon: '◑',
    textKey: 'dashboard.status.stopping',
    color: '#f59e0b', // amber-500
  },
  Error: {
    variant: 'error',
    icon: '✕',
    textKey: 'dashboard.status.error',
    color: '#ef4444', // red-500
  },
};

export const StatusIndicator: React.FC<StatusIndicatorProps> = React.memo(({
  status,
  size = 'md',
  showText = true,
  className = '',
}) => {
  const { t } = useI18n();
  const config = statusConfig[status];
  const isTransitioning = status === 'Starting' || status === 'Stopping';

  return (
    <Badge
      variant={config.variant}
      size={size}
      animated={isTransitioning}
      className={className}
    >
      {/* Icon with optional rotation for transitioning states */}
      <motion.span
        className="text-sm leading-none"
        animate={isTransitioning ? { rotate: 360 } : {}}
        transition={
          isTransitioning
            ? {
                duration: 2,
                repeat: Infinity,
                ease: 'linear',
              }
            : {}
        }
        style={{ color: config.color }}
      >
        {config.icon}
      </motion.span>

      {/* Status text */}
      {showText && (
        <span className="font-medium">
          {t(config.textKey)}
        </span>
      )}
    </Badge>
  );
});
