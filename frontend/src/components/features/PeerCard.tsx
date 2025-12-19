import React from 'react';
import { motion } from 'framer-motion';
import { GlassCard } from '../layout/GlassCard';
import { Badge } from '../ui/Badge';
import { Button } from '../ui/Button';
import { useI18n } from '../../hooks/useI18n';

export interface PeerInfo {
  address: string;
  connected: boolean;
  latency?: number;
  rxBytes?: number;
  txBytes?: number;
  uptime?: number;
}

interface PeerCardProps {
  peer: PeerInfo;
  onToggle?: (address: string) => void;
  onRemove?: (address: string) => void;
  showActions?: boolean;
  variant?: 'default' | 'compact';
}

// Format bytes to human-readable format
const formatBytes = (bytes?: number): string => {
  if (!bytes) return '0 B';
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(1024));
  return `${(bytes / Math.pow(1024, i)).toFixed(2)} ${sizes[i]}`;
};

// Format uptime to human-readable format
const formatUptime = (seconds?: number): string => {
  if (!seconds) return '0s';
  const days = Math.floor(seconds / 86400);
  const hours = Math.floor((seconds % 86400) / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);

  if (days > 0) return `${days}d ${hours}h`;
  if (hours > 0) return `${hours}h ${minutes}m`;
  return `${minutes}m`;
};

// Truncate address
const truncateAddress = (address: string, length: number = 20): string => {
  if (address.length <= length) return address;
  const start = Math.floor(length / 2);
  const end = Math.ceil(length / 2);
  return `${address.substring(0, start)}...${address.substring(address.length - end)}`;
};

export const PeerCard: React.FC<PeerCardProps> = React.memo(({
  peer,
  onToggle,
  onRemove,
  showActions = true,
  variant = 'default',
}) => {
  const { t } = useI18n();

  return (
    <GlassCard
      variant="default"
      padding="md"
      hoverable
      accentColor={peer.connected ? '#006C4C' : undefined}
    >
      <div className="space-y-3">
        {/* Header: Address + Status */}
        <div className="flex items-start justify-between gap-3">
          <div className="flex-1 min-w-0">
            <p className="text-sm font-mono text-md-light-onSurface dark:text-md-dark-onSurface truncate" title={peer.address}>
              {truncateAddress(peer.address, 30)}
            </p>
          </div>
          <Badge
            variant={peer.connected ? 'success' : 'default'}
            size="sm"
            animated={false}
          >
            <span className="text-base leading-none">
              {peer.connected ? '●' : '○'}
            </span>
            <span>{peer.connected ? t('peers.status.connected') : t('peers.status.offline')}</span>
          </Badge>
        </div>

        {/* Stats (only show if connected) */}
        {peer.connected && (
          variant === 'compact' ? (
            <div className="flex flex-wrap gap-x-4 gap-y-1 text-xs">
              {peer.latency !== undefined && (
                <span className="text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant">
                  <span className="text-md-light-outline dark:text-md-dark-outline">Latency:</span> <span className="font-mono font-semibold">{peer.latency}ms</span>
                </span>
              )}
              {peer.uptime !== undefined && (
                <span className="text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant">
                  <span className="text-md-light-outline dark:text-md-dark-outline">Uptime:</span> <span className="font-mono font-semibold">{formatUptime(peer.uptime)}</span>
                </span>
              )}
              {peer.rxBytes !== undefined && (
                <span className="text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant">
                  <span className="text-md-light-outline dark:text-md-dark-outline">Down:</span> <span className="font-mono font-semibold text-md-light-tertiary dark:text-md-dark-tertiary">{formatBytes(peer.rxBytes)}</span>
                </span>
              )}
              {peer.txBytes !== undefined && (
                <span className="text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant">
                  <span className="text-md-light-outline dark:text-md-dark-outline">Up:</span> <span className="font-mono font-semibold text-md-light-secondary dark:text-md-dark-secondary">{formatBytes(peer.txBytes)}</span>
                </span>
              )}
            </div>
          ) : (
            <div className="grid grid-cols-2 gap-3 text-sm">
              {/* Latency */}
              {peer.latency !== undefined && (
                <div className="flex flex-col">
                  <span className="text-md-light-outline dark:text-md-dark-outline text-xs font-body">Latency</span>
                  <span className="text-md-light-onSurface dark:text-md-dark-onSurface font-semibold font-mono">{peer.latency}ms</span>
                </div>
              )}

              {/* Uptime */}
              {peer.uptime !== undefined && (
                <div className="flex flex-col">
                  <span className="text-md-light-outline dark:text-md-dark-outline text-xs font-body">Uptime</span>
                  <span className="text-md-light-onSurface dark:text-md-dark-onSurface font-semibold font-mono">{formatUptime(peer.uptime)}</span>
                </div>
              )}

              {/* RX (Download) */}
              {peer.rxBytes !== undefined && (
                <div className="flex flex-col">
                  <span className="text-md-light-outline dark:text-md-dark-outline text-xs font-body">Downloaded</span>
                  <span className="text-md-light-tertiary dark:text-md-dark-tertiary font-semibold font-mono">{formatBytes(peer.rxBytes)}</span>
                </div>
              )}

              {/* TX (Upload) */}
              {peer.txBytes !== undefined && (
                <div className="flex flex-col">
                  <span className="text-md-light-outline dark:text-md-dark-outline text-xs font-body">Uploaded</span>
                  <span className="text-md-light-secondary dark:text-md-dark-secondary font-semibold font-mono">{formatBytes(peer.txBytes)}</span>
                </div>
              )}
            </div>
          )
        )}

        {/* Actions */}
        {showActions && (
          <div className="flex gap-2 pt-2 border-t border-md-light-outline/30 dark:border-md-dark-outline/30">
            {onToggle && (
              <Button
                variant="ghost"
                size="sm"
                onClick={() => onToggle(peer.address)}
                fullWidth
              >
                {peer.connected ? t('peers.card.disable') : t('peers.card.enable')}
              </Button>
            )}
            {onRemove && (
              <Button
                variant="danger"
                size="sm"
                onClick={() => onRemove(peer.address)}
              >
                <svg
                  className="w-4 h-4"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                  />
                </svg>
              </Button>
            )}
          </div>
        )}
      </div>
    </GlassCard>
  );
});
