import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
  Modal,
  Button,
  Input,
  LoadingSpinner,
  Badge,
  HolographicBorder,
  GlassCard,
} from '..';
import { useI18n } from '../../hooks/useI18n';
import {
  FindAvailablePeers,
  GetCachedDiscoveredPeers,
  ClearCachedDiscoveredPeers,
  GetAvailableRegions,
  CancelPeerDiscovery,
} from '../../../wailsjs/go/main/App';
import { core } from '../../../wailsjs/go/models';
import { EventsOn, EventsOff } from '../../../wailsjs/runtime/runtime';
import { showSuccess, showError } from '../../store/uiStore';

interface PeerDiscoveryModalProps {
  isOpen: boolean;
  onClose: () => void;
  onPeersAdded?: (peers: core.DiscoveredPeer[]) => void;
}

interface DiscoveryProgress {
  current: number;
  total: number;
  available_count: number;
}

export function PeerDiscoveryModal({ isOpen, onClose, onPeersAdded }: PeerDiscoveryModalProps) {
  const { t } = useI18n();
  const [isSearching, setIsSearching] = useState(false);
  const [discoveredPeers, setDiscoveredPeers] = useState<core.DiscoveredPeer[]>([]);
  const [selectedPeers, setSelectedPeers] = useState<Set<string>>(new Set());
  const [progress, setProgress] = useState<DiscoveryProgress>({ current: 0, total: 0, available_count: 0 });

  // Filters
  const [selectedProtocols, setSelectedProtocols] = useState<string[]>(['tls', 'quic']);
  const [selectedRegion, setSelectedRegion] = useState<string>('');
  const [maxRTT, setMaxRTT] = useState<number>(5000);
  const [availableRegions, setAvailableRegions] = useState<string[]>([]);

  // View state
  const [showCached, setShowCached] = useState(false);

  const protocols = [
    { value: 'tcp', label: 'TCP' },
    { value: 'tls', label: 'TLS' },
    { value: 'quic', label: 'QUIC' },
    { value: 'ws', label: 'WebSocket' },
    { value: 'wss', label: 'WebSocket (TLS)' },
  ];

  // Load cached peers on mount
  useEffect(() => {
    if (isOpen) {
      loadCachedPeers();
      loadAvailableRegions();
    }
  }, [isOpen]);

  // Subscribe to progress events
  useEffect(() => {
    const unsubscribe = EventsOn('peer:discovery:progress', (data: DiscoveryProgress) => {
      setProgress(data);
    });

    return () => {
      EventsOff('peer:discovery:progress');
      if (unsubscribe) unsubscribe();
    };
  }, []);

  const loadCachedPeers = async () => {
    try {
      const cached = await GetCachedDiscoveredPeers();
      if (cached && cached.length > 0) {
        setDiscoveredPeers(cached);
        setShowCached(true);
      }
    } catch (error) {
      console.error('Failed to load cached peers:', error);
    }
  };

  const loadAvailableRegions = async () => {
    try {
      const regions = await GetAvailableRegions();
      setAvailableRegions(regions || []);
    } catch (error) {
      console.error('Failed to load regions:', error);
    }
  };

  const handleSearch = async () => {
    setIsSearching(true);
    setShowCached(false);
    setDiscoveredPeers([]);
    setSelectedPeers(new Set());
    setProgress({ current: 0, total: 0, available_count: 0 });

    try {
      const result = await FindAvailablePeers(
        selectedProtocols.join(','),
        selectedRegion,
        maxRTT
      );

      if (result && result.peers) {
        setDiscoveredPeers(result.peers);
        showSuccess(
          t('peers.discovery.searchComplete'),
          t('peers.discovery.foundPeers', { count: result.available, total: result.total })
        );
      }
    } catch (error) {
      showError(
        t('peers.discovery.searchFailed'),
        error instanceof Error ? error.message : t('peers.discovery.searchFailedMessage')
      );
    } finally {
      setIsSearching(false);
    }
  };

  const handleStopSearch = async () => {
    try {
      await CancelPeerDiscovery();
      setIsSearching(false);
      showSuccess(t('peers.discovery.searchStopped'), t('peers.discovery.searchStoppedMessage'));
    } catch (error) {
      showError(t('action.error'), error instanceof Error ? error.message : 'Failed to stop search');
    }
  };

  const handleClearCache = async () => {
    try {
      await ClearCachedDiscoveredPeers();
      setDiscoveredPeers([]);
      setShowCached(false);
      showSuccess(t('peers.discovery.cacheCleared'), t('peers.discovery.cacheClearedMessage'));
    } catch (error) {
      showError(t('action.error'), error instanceof Error ? error.message : 'Failed to clear cache');
    }
  };

  const handleToggleSelect = (address: string) => {
    const newSelected = new Set(selectedPeers);
    if (newSelected.has(address)) {
      newSelected.delete(address);
    } else {
      newSelected.add(address);
    }
    setSelectedPeers(newSelected);
  };

  const handleSelectAll = () => {
    if (selectedPeers.size === discoveredPeers.length) {
      setSelectedPeers(new Set());
    } else {
      setSelectedPeers(new Set(discoveredPeers.map(p => p.address)));
    }
  };

  const handleAddSelected = () => {
    if (selectedPeers.size === 0) {
      showError(t('peers.discovery.noPeersSelected'), t('peers.discovery.noPeersSelectedMessage'));
      return;
    }

    const peersToAdd = discoveredPeers.filter(p => selectedPeers.has(p.address));

    showSuccess(
      t('peers.discovery.peersAdded'),
      t('peers.discovery.peersAddedMessage', { count: selectedPeers.size })
    );

    onPeersAdded?.(peersToAdd);
    onClose();
  };

  const toggleProtocol = (protocol: string) => {
    setSelectedProtocols(prev =>
      prev.includes(protocol)
        ? prev.filter(p => p !== protocol)
        : [...prev, protocol]
    );
  };

  const progressPercent = progress.total > 0
    ? Math.round((progress.current / progress.total) * 100)
    : 0;

  const sortedPeers = [...discoveredPeers].sort((a, b) => a.rtt - b.rtt);

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={t('peers.discovery.title')}
      size="xl"
    >
      <div className="space-y-4 max-h-[70vh] overflow-y-auto pr-2 scrollbar-thin">
        {/* Filters */}
        <GlassCard padding="md" variant="subtle">
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-slate-200 mb-2">
                {t('peers.discovery.protocols')}
              </label>
              <div className="flex flex-wrap gap-2">
                {protocols.map(proto => (
                  <button
                    key={proto.value}
                    onClick={() => toggleProtocol(proto.value)}
                    className={`px-4 py-2 rounded-lg font-medium transition-all ${
                      selectedProtocols.includes(proto.value)
                        ? 'bg-emerald-500/20 border-2 border-emerald-500 text-emerald-400'
                        : 'bg-slate-700 border border-slate-600 text-slate-300 hover:bg-slate-600'
                    }`}
                  >
                    {proto.label}
                  </button>
                ))}
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-slate-200 mb-2">
                  {t('peers.discovery.region')}
                </label>
                <select
                  value={selectedRegion}
                  onChange={(e) => setSelectedRegion(e.target.value)}
                  className="w-full px-4 py-2 bg-slate-800 border border-slate-600 rounded-xl text-slate-100 focus:border-emerald-500 focus:ring-2 focus:ring-emerald-500/50 focus:outline-none [&>option]:bg-slate-800 [&>option]:text-slate-100"
                >
                  <option value="">{t('peers.discovery.allRegions')}</option>
                  {availableRegions.map(region => (
                    <option key={region} value={region}>{region}</option>
                  ))}
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-200 mb-2">
                  {t('peers.discovery.maxRTT')}
                </label>
                <input
                  type="number"
                  value={maxRTT}
                  onChange={(e) => setMaxRTT(parseInt(e.target.value) || 5000)}
                  className="w-full px-4 py-2 bg-slate-800 border border-slate-600 rounded-xl text-slate-100 focus:border-emerald-500 focus:ring-2 focus:ring-emerald-500/50 focus:outline-none"
                />
                <p className="text-xs text-slate-400 mt-1">{t('peers.discovery.maxRTTHelper')}</p>
              </div>
            </div>

            <div className="flex gap-3">
              {isSearching ? (
                <Button
                  variant="danger"
                  onClick={handleStopSearch}
                  className="flex-1"
                >
                  {t('peers.discovery.stopButton')}
                </Button>
              ) : (
                <Button
                  variant="primary"
                  onClick={handleSearch}
                  disabled={selectedProtocols.length === 0}
                  className="flex-1"
                >
                  {t('peers.discovery.searchButton')}
                </Button>
              )}
              {showCached && !isSearching && (
                <Button
                  variant="ghost"
                  onClick={handleClearCache}
                >
                  {t('peers.discovery.clearCache')}
                </Button>
              )}
            </div>
          </div>
        </GlassCard>

        {/* Progress */}
        {isSearching && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
          >
            <GlassCard padding="md" variant="strong">
              <div className="space-y-3">
                <div className="flex items-center justify-between text-sm">
                  <span className="text-slate-200">
                    {t('peers.discovery.checking')} {progress.current} / {progress.total}
                  </span>
                  <span className="text-emerald-400 font-bold">
                    {progress.available_count} {t('peers.discovery.available')}
                  </span>
                </div>
                <div className="relative h-2 bg-slate-700 rounded-full overflow-hidden">
                  <motion.div
                    className="absolute inset-y-0 left-0 bg-gradient-to-r from-emerald-500 to-emerald-400"
                    initial={{ width: 0 }}
                    animate={{ width: `${progressPercent}%` }}
                    transition={{ duration: 0.3 }}
                  />
                </div>
                <p className="text-xs text-slate-400 text-center">
                  {progressPercent}% {t('peers.discovery.complete')}
                </p>
              </div>
            </GlassCard>
          </motion.div>
        )}

        {/* Results Header */}
        {discoveredPeers.length > 0 && (
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <h3 className="text-lg font-semibold text-slate-100">
                {showCached ? t('peers.discovery.cachedResults') : t('peers.discovery.results')}
              </h3>
              <Badge variant="success">{discoveredPeers.length}</Badge>
            </div>
            <div className="flex gap-2">
              <Button
                variant="ghost"
                size="sm"
                onClick={handleSelectAll}
              >
                {selectedPeers.size === discoveredPeers.length
                  ? t('action.deselectAll')
                  : t('action.selectAll')}
              </Button>
              <Button
                variant="primary"
                size="sm"
                onClick={handleAddSelected}
                disabled={selectedPeers.size === 0}
              >
                {t('peers.discovery.addSelected')} ({selectedPeers.size})
              </Button>
            </div>
          </div>
        )}

        {/* Discovered Peers List */}
        {discoveredPeers.length > 0 && (
          <div className="max-h-64 overflow-y-auto space-y-2 pr-2 scrollbar-thin">
            <AnimatePresence>
              {sortedPeers.map((peer, index) => (
                <motion.div
                  key={peer.address}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: 20 }}
                  transition={{ duration: 0.2, delay: index * 0.03 }}
                >
                  <HolographicBorder
                    borderWidth={selectedPeers.has(peer.address) ? 2 : 1}
                    animated={selectedPeers.has(peer.address)}
                  >
                    <button
                      onClick={() => handleToggleSelect(peer.address)}
                      className={`w-full text-left p-4 rounded-lg transition-all ${
                        selectedPeers.has(peer.address)
                          ? 'bg-emerald-500/10'
                          : 'bg-slate-700/50 hover:bg-slate-700'
                      }`}
                    >
                      <div className="flex items-center justify-between gap-4">
                        <div className="flex-1 min-w-0">
                          <div className="flex items-center gap-2 mb-1">
                            <Badge variant={peer.protocol as any}>{peer.protocol.toUpperCase()}</Badge>
                            {peer.region && (
                              <span className="text-xs text-slate-400">📍 {peer.region}</span>
                            )}
                          </div>
                          <p className="text-sm font-mono text-slate-200 truncate">
                            {peer.address}
                          </p>
                        </div>
                        <div className="flex items-center gap-4">
                          <div className="text-right">
                            <p className="text-xs text-slate-500">{t('peers.discovery.rtt')}</p>
                            <p className={`text-sm font-bold ${
                              peer.rtt < 100 ? 'text-emerald-400' :
                              peer.rtt < 300 ? 'text-yellow-400' :
                              'text-red-400'
                            }`}>
                              {peer.rtt}ms
                            </p>
                          </div>
                          <div className="w-6 h-6 rounded-full border-2 flex items-center justify-center">
                            {selectedPeers.has(peer.address) && (
                              <motion.div
                                initial={{ scale: 0 }}
                                animate={{ scale: 1 }}
                                className="w-3 h-3 rounded-full bg-emerald-500"
                              />
                            )}
                          </div>
                        </div>
                      </div>
                    </button>
                  </HolographicBorder>
                </motion.div>
              ))}
            </AnimatePresence>
          </div>
        )}

        {/* Empty State */}
        {!isSearching && discoveredPeers.length === 0 && (
          <div className="text-center py-12 space-y-4">
            <div className="text-6xl mb-4">🔍</div>
            <p className="text-lg text-slate-100/70">
              {t('peers.discovery.noPeersFound')}
            </p>
            <p className="text-sm text-slate-500">
              {t('peers.discovery.noPeersFoundMessage')}
            </p>
          </div>
        )}

        {/* Info */}
        <div className="bg-emerald-500/10 border border-emerald-500/30 rounded-lg p-4">
          <p className="text-xs text-slate-100/80">
            💡 {t('peers.discovery.infoMessage')}
          </p>
        </div>
      </div>
    </Modal>
  );
}

export default PeerDiscoveryModal;
