import React from 'react';
import { GlassCard } from '../layout/GlassCard';
import { Badge } from '../ui/Badge';
import { Button } from '../ui/Button';
import { useI18n } from '../../hooks/useI18n';

export interface PeerInfo {
  address: string;
  enabled: boolean;
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

const formatBytes = (bytes?: number): string => {
  if (!bytes) return '0 B';
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(1024));
  return `${(bytes / Math.pow(1024, i)).toFixed(2)} ${sizes[i]}`;
};

const formatUptime = (seconds?: number): string => {
  if (!seconds) return '0s';
  const days = Math.floor(seconds / 86400);
  const hours = Math.floor((seconds % 86400) / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);

  if (days > 0) return `${days}d ${hours}h`;
  if (hours > 0) return `${hours}h ${minutes}m`;
  return `${minutes}m`;
};

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
      accentColor={peer.connected ? '#10b981' : undefined}
    >
      <div className="space-y-3">
        {/* Header: Address + Status */}
        <div className="flex items-start justify-between gap-3">
          <div className="flex-1 min-w-0">
            <p className="text-sm font-mono text-slate-200 truncate" title={peer.address}>
              {truncateAddress(peer.address, 30)}
            </p>
            {!peer.enabled && (
              <p className="text-xs text-slate-500 mt-1">
                {t('peers.status.disabled')}
              </p>
            )}
          </div>
          <Badge
            variant={peer.connected ? 'success' : peer.enabled ? 'warning' : 'default'}
            size="sm"
            animated={false}
          >
            <span className="text-sm leading-none">
              {peer.connected ? '●' : peer.enabled ? '◐' : '○'}
            </span>
            <span>
              {peer.connected
                ? t('peers.status.connected')
                : peer.enabled
                  ? t('peers.status.offline')
                  : t('peers.status.disabled')
              }
            </span>
          </Badge>
        </div>

        {/* Stats (only show if connected) */}
        {peer.connected && (
          variant === 'compact' ? (
            <div className="flex flex-wrap gap-x-4 gap-y-1 text-xs">
              {peer.latency !== undefined && (
                <span className="text-slate-300">
                  <span className="text-slate-500">Latency:</span> <span className="font-mono font-medium">{peer.latency}ms</span>
                </span>
              )}
              {peer.uptime !== undefined && (
                <span className="text-slate-300">
                  <span className="text-slate-500">Uptime:</span> <span className="font-mono font-medium">{formatUptime(peer.uptime)}</span>
                </span>
              )}
              {peer.rxBytes !== undefined && (
                <span className="text-slate-300">
                  <span className="text-slate-500">Down:</span> <span className="font-mono font-medium text-blue-400">{formatBytes(peer.rxBytes)}</span>
                </span>
              )}
              {peer.txBytes !== undefined && (
                <span className="text-slate-300">
                  <span className="text-slate-500">Up:</span> <span className="font-mono font-medium text-emerald-400">{formatBytes(peer.txBytes)}</span>
                </span>
              )}
            </div>
          ) : (
            <div className="grid grid-cols-2 gap-3 text-sm">
              {peer.latency !== undefined && (
                <div className="flex flex-col">
                  <span className="text-slate-500 text-xs">Latency</span>
                  <span className="text-slate-200 font-medium font-mono">{peer.latency}ms</span>
                </div>
              )}

              {peer.uptime !== undefined && (
                <div className="flex flex-col">
                  <span className="text-slate-500 text-xs">Uptime</span>
                  <span className="text-slate-200 font-medium font-mono">{formatUptime(peer.uptime)}</span>
                </div>
              )}

              {peer.rxBytes !== undefined && (
                <div className="flex flex-col">
                  <span className="text-slate-500 text-xs">Downloaded</span>
                  <span className="text-blue-400 font-medium font-mono">{formatBytes(peer.rxBytes)}</span>
                </div>
              )}

              {peer.txBytes !== undefined && (
                <div className="flex flex-col">
                  <span className="text-slate-500 text-xs">Uploaded</span>
                  <span className="text-emerald-400 font-medium font-mono">{formatBytes(peer.txBytes)}</span>
                </div>
              )}
            </div>
          )
        )}

        {/* Actions */}
        {showActions && (
          <div className="flex gap-2 pt-2 border-t border-slate-700">
            {onToggle && (
              <Button
                variant="ghost"
                size="sm"
                onClick={() => onToggle(peer.address)}
                fullWidth
              >
                {peer.enabled ? t('peers.card.disable') : t('peers.card.enable')}
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
