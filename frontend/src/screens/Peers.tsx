import { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import {
  Button,
  Input,
  Modal,
  GlassCard,
  HolographicBorder,
  PeerCard,
  LoadingSpinner,
  PeerDiscoveryModal,
} from '../components';
import { useServiceStatus } from '../hooks/useServiceStatus';
import { useConfig } from '../hooks/useConfig';
import { useI18n } from '../hooks/useI18n';
import {
  HotReloadPeers,
  GetConfig,
  SaveConfig,
} from '../../wailsjs/go/main/App';
import { toast } from '../components/ui/Toast';

/**
 * Peers Screen - Peer management
 *
 * Features:
 * - Peer list with cards (status, latency, bandwidth)
 * - Add new peer dialog
 * - Enable/Disable toggle (TODO: backend support)
 * - Delete confirmation
 * - Apply button (appears after changes, reloads peer configuration)
 */
export function Peers() {
  const { t } = useI18n();
  const [showAddModal, setShowAddModal] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [showDiscoveryModal, setShowDiscoveryModal] = useState(false);
  const [peerToDelete, setPeerToDelete] = useState<string | null>(null);
  const [newPeerAddress, setNewPeerAddress] = useState('');
  const [isProcessing, setIsProcessing] = useState(false);
  const [isHotReloading, setIsHotReloading] = useState(false);
  const [hasChanges, setHasChanges] = useState(false);

  // Local state for pending peer changes (not yet saved to config)
  const [localPeers, setLocalPeers] = useState<Array<{ address: string; enabled: boolean }>>([]);

  // Get peers from service status
  const { peers, fetchPeerStats } = useServiceStatus({
    refreshInterval: 5000,
    fetchOnMount: true,
  });

  // Get config for peer list
  const { config, loadConfig } = useConfig({ loadOnMount: true });

  // Initialize localPeers from config when it loads
  useEffect(() => {
    if (config?.peers && Array.isArray(config.peers)) {
      setLocalPeers(config.peers.map(p => ({ address: p.address, enabled: p.enabled })));
    } else {
      setLocalPeers([]);
    }
  }, [config]);

  // Auto-refresh peer stats
  useEffect(() => {
    const interval = setInterval(() => {
      fetchPeerStats();
    }, 5000);
    return () => clearInterval(interval);
  }, [fetchPeerStats]);

  // Handle add peer
  const handleAddPeer = async () => {
    if (!newPeerAddress.trim()) {
      toast.error(t('peers.messages.peerAddressEmpty'));
      return;
    }

    // Validate format (basic check)
    if (!newPeerAddress.includes('://')) {
      toast.error(t('peers.messages.invalidFormat'));
      return;
    }

    // Check if peer already exists
    if (localPeers.some(p => p.address === newPeerAddress.trim())) {
      toast.error(t('peers.messages.duplicatePeerMessage'));
      return;
    }

    setIsProcessing(true);
    try {
      // Add to local state only (not saving to config yet)
      setLocalPeers(prev => [...prev, { address: newPeerAddress.trim(), enabled: true }]);
      setNewPeerAddress('');
      setShowAddModal(false);
      setHasChanges(true);
      toast.success(t('peers.messages.peerAddedMessage'));
    } finally {
      setIsProcessing(false);
    }
  };

  // Handle delete peer
  const handleDeletePeer = async () => {
    if (!peerToDelete) return;

    setIsProcessing(true);
    try {
      // Remove from local state only (not saving to config yet)
      setLocalPeers(prev => prev.filter(p => p.address !== peerToDelete));
      setPeerToDelete(null);
      setShowDeleteModal(false);
      setHasChanges(true);
      toast.success(t('peers.messages.peerRemovedMessage'));
    } finally {
      setIsProcessing(false);
    }
  };

  // Handle apply changes (save to config and hot reload)
  const handleApplyChanges = async () => {
    setIsHotReloading(true);
    try {
      // First, save all changes to config
      // We need to completely replace the peer list in config
      const currentConfig = await GetConfig();

      // Update only the peers field
      currentConfig.peers = localPeers.map(p => ({
        address: p.address,
        enabled: p.enabled,
      }));

      await SaveConfig(currentConfig);

      // Reload config to sync
      await loadConfig();

      // Then hot-reload peers in the running service (no restart needed)
      await HotReloadPeers();
      await fetchPeerStats();

      setHasChanges(false);
      toast.success(t('peers.messages.changesAppliedMessage'));
    } catch (error) {
      toast.error(error instanceof Error ? error.message : t('peers.messages.applyFailedMessage'));
    } finally {
      setIsHotReloading(false);
    }
  };

  // Handle toggle peer
  const handleTogglePeer = async (address: string) => {
    // Find the peer in local state
    const localPeer = localPeers.find(p => p.address === address);
    if (!localPeer) return;

    setIsProcessing(true);
    try {
      // Toggle in local state only (not saving to config yet)
      setLocalPeers(prev =>
        prev.map(p =>
          p.address === address ? { ...p, enabled: !p.enabled } : p
        )
      );

      setHasChanges(true);
      const status = localPeer.enabled ? 'disabled' : 'enabled';
      toast.success(t(`peers.messages.peer${status === 'enabled' ? 'Enabled' : 'Disabled'}Message`));
    } finally {
      setIsProcessing(false);
    }
  };

  // Handle remove peer (show confirmation)
  const handleRemove = (address: string) => {
    setPeerToDelete(address);
    setShowDeleteModal(true);
  };

  // Merge local peers config with live peer stats
  const mergedPeers = localPeers.map(localPeer => {
    const livePeer = peers.find(p => p.address === localPeer.address);
    return {
      address: localPeer.address,
      enabled: localPeer.enabled,
      connected: livePeer?.connected || false,
      latency: livePeer?.latency || 0,
      rxBytes: livePeer?.rxBytes || 0,
      txBytes: livePeer?.txBytes || 0,
      uptime: livePeer?.uptime || 0,
    };
  });

  const connectedCount = mergedPeers.filter(p => p.connected).length;
  const totalCount = mergedPeers.length;

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
            {t('peers.peerManagement')}
          </h1>
          <p className="text-lg text-white/70 font-futuristic">
            {t('peers.peersConnected', { connected: connectedCount, total: totalCount })}
          </p>
        </div>
        <div className="flex gap-3">
          {hasChanges && (
            <Button
              variant="secondary"
              glow
              onClick={handleApplyChanges}
              disabled={isHotReloading}
            >
              {isHotReloading ? t('peers.applying') : t('peers.apply')}
            </Button>
          )}
          <Button
            variant="secondary"
            onClick={() => setShowDiscoveryModal(true)}
          >
            🔍 {t('peers.findPeersButton')}
          </Button>
          <Button
            variant="primary"
            glow
            onClick={() => setShowAddModal(true)}
          >
            {t('peers.addPeerButton')}
          </Button>
        </div>
      </motion.div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <motion.div
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.3, delay: 0.1 }}
        >
          <HolographicBorder borderWidth={2} animated>
            <GlassCard padding="md" variant="strong">
              <div className="text-center space-y-2">
                <p className="text-xs text-white/50 font-futuristic uppercase tracking-wide">
                  {t('peers.stats.connected')}
                </p>
                <p className="text-5xl font-bold text-neon-green">
                  {connectedCount}
                </p>
              </div>
            </GlassCard>
          </HolographicBorder>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.3, delay: 0.2 }}
        >
          <HolographicBorder borderWidth={2}>
            <GlassCard padding="md" variant="strong">
              <div className="text-center space-y-2">
                <p className="text-xs text-white/50 font-futuristic uppercase tracking-wide">
                  {t('peers.stats.disconnected')}
                </p>
                <p className="text-5xl font-bold text-neon-pink">
                  {totalCount - connectedCount}
                </p>
              </div>
            </GlassCard>
          </HolographicBorder>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.3, delay: 0.3 }}
        >
          <HolographicBorder borderWidth={2}>
            <GlassCard padding="md" variant="strong">
              <div className="text-center space-y-2">
                <p className="text-xs text-white/50 font-futuristic uppercase tracking-wide">
                  {t('peers.stats.totalPeers')}
                </p>
                <p className="text-5xl font-bold text-neon-cyan">
                  {totalCount}
                </p>
              </div>
            </GlassCard>
          </HolographicBorder>
        </motion.div>
      </div>

      {/* Peer List */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.4 }}
      >
        <GlassCard
          title={t('peers.allPeers')}
          subtitle={t('peers.allPeersSubtitle')}
          padding="lg"
        >
          {mergedPeers.length > 0 ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {mergedPeers.map((peer, index) => (
                <motion.div
                  key={peer.address}
                  initial={{ opacity: 0, scale: 0.9 }}
                  animate={{ opacity: 1, scale: 1 }}
                  transition={{ duration: 0.3, delay: 0.5 + index * 0.05 }}
                >
                  <PeerCard
                    peer={{
                      address: peer.address,
                      enabled: peer.enabled,
                      connected: peer.connected,
                      latency: peer.latency,
                      rxBytes: peer.rxBytes,
                      txBytes: peer.txBytes,
                      uptime: peer.uptime,
                    }}
                    showActions
                    onToggle={handleTogglePeer}
                    onRemove={handleRemove}
                  />
                </motion.div>
              ))}
            </div>
          ) : (
            <div className="text-center py-16 space-y-4">
              <div className="text-8xl mb-4">🌐</div>
              <p className="text-xl text-white/70">{t('peers.noPeersConfigured')}</p>
              <p className="text-sm text-white/50">
                {t('peers.addPeersPrompt')}
              </p>
              <Button
                variant="primary"
                glow
                onClick={() => setShowAddModal(true)}
                className="mt-4"
              >
                {t('peers.addFirstPeer')}
              </Button>
            </div>
          )}
        </GlassCard>
      </motion.div>

      {/* Info Box */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.6 }}
      >
        <HolographicBorder borderWidth={1}>
          <div className="bg-neon-cyan/10 border border-neon-cyan/30 rounded-lg p-6">
            <div className="flex items-start gap-4">
              <div className="text-4xl">💡</div>
              <div className="flex-1 space-y-2">
                <h3 className="text-white font-semibold text-lg">{t('peers.aboutPeers')}</h3>
                <p className="text-sm text-white/80 leading-relaxed">
                  {t('peers.aboutDescription')}
                </p>
                <div className="pt-2">
                  <p className="text-xs text-white/60">
                    {t('peers.defaultPeers')}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </HolographicBorder>
      </motion.div>

      {/* Add Peer Modal */}
      <Modal
        isOpen={showAddModal}
        onClose={() => setShowAddModal(false)}
        title={t('peers.modal.addPeer')}
        size="md"
      >
        <div className="space-y-4">
          <Input
            label={t('peers.modal.peerAddress')}
            placeholder={t('peers.modal.peerAddressPlaceholder')}
            value={newPeerAddress}
            onChange={(e) => setNewPeerAddress(e.target.value)}
            helperText={t('peers.modal.helperText')}
            autoFocus
          />
          <div className="bg-neon-cyan/10 border border-neon-cyan/30 rounded-lg p-4">
            <p className="text-sm text-white/90">
              <strong>{t('peers.modal.formatTitle')}</strong> {t('peers.modal.formatDescription')}
            </p>
            <ul className="list-disc list-inside text-sm text-white/70 mt-2 space-y-1">
              <li>{t('peers.modal.formatItem1')}</li>
              <li>{t('peers.modal.formatItem2')}</li>
              <li>{t('peers.modal.formatItem3')}</li>
            </ul>
          </div>
          <div className="flex gap-3 justify-end pt-4">
            <Button
              variant="ghost"
              onClick={() => setShowAddModal(false)}
              disabled={isProcessing}
            >
              {t('action.cancel')}
            </Button>
            <Button
              variant="primary"
              glow
              onClick={handleAddPeer}
              disabled={isProcessing}
            >
              {isProcessing ? t('peers.modal.adding') : t('peers.modal.addButton')}
            </Button>
          </div>
        </div>
      </Modal>

      {/* Delete Confirmation Modal */}
      <Modal
        isOpen={showDeleteModal}
        onClose={() => setShowDeleteModal(false)}
        title={t('peers.modal.removePeer')}
        size="sm"
      >
        <div className="space-y-4">
          <div className="bg-red-500/10 border border-red-500/30 rounded-lg p-4">
            <p className="text-white/90">
              {t('peers.modal.confirmRemove')}
            </p>
            {peerToDelete && (
              <p className="text-sm text-white/70 font-mono mt-2 break-all">
                {peerToDelete}
              </p>
            )}
          </div>
          <p className="text-sm text-white/70">
            {t('peers.modal.removeDescription')}
          </p>
          <div className="flex gap-3 justify-end pt-4">
            <Button
              variant="ghost"
              onClick={() => setShowDeleteModal(false)}
              disabled={isProcessing}
            >
              {t('action.cancel')}
            </Button>
            <Button
              variant="danger"
              onClick={handleDeletePeer}
              disabled={isProcessing}
            >
              {isProcessing ? t('peers.modal.removing') : t('peers.modal.removeButton')}
            </Button>
          </div>
        </div>
      </Modal>

      {/* Peer Discovery Modal */}
      <PeerDiscoveryModal
        isOpen={showDiscoveryModal}
        onClose={() => setShowDiscoveryModal(false)}
        onPeersAdded={(discoveredPeers) => {
          // Add discovered peers to local state without saving
          const newPeers = discoveredPeers.map(p => ({
            address: p.address,
            enabled: true,
          }));

          // Filter out duplicates
          const existingAddresses = new Set(localPeers.map(p => p.address));
          const uniqueNewPeers = newPeers.filter(p => !existingAddresses.has(p.address));

          if (uniqueNewPeers.length > 0) {
            setLocalPeers(prev => [...prev, ...uniqueNewPeers]);
            setHasChanges(true);
          }
        }}
      />
    </div>
  );
}

export default Peers;
