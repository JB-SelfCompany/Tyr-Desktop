import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { motion } from 'framer-motion';
import {
  Button,
  GlassCard,
  HolographicBorder,
  StatusIndicator,
  PeerCard,
} from '../components';
import { useServiceStatus } from '../hooks/useServiceStatus';
import { useI18n } from '../hooks/useI18n';
import { CopyToClipboard, OpenDeltaChat } from '../../wailsjs/go/main/App';
import { showSuccess, showError } from '../store/uiStore';
import type { ServiceStatus } from '../components';

/**
 * Dashboard Screen - Main application screen
 *
 * Features:
 * - Service status with animated indicator
 * - Mail address with copy and DeltaChat button
 * - Server info (SMTP/IMAP addresses)
 * - Connected peers grid
 * - Start/Stop service control
 *
 * Layout: Bento Grid (3 columns) with Y2K Futurism design
 */
export function Dashboard() {
  const { t } = useI18n();
  const navigate = useNavigate();
  const [copiedField, setCopiedField] = useState<string | null>(null);

  // Service status hook with auto-refresh
  const {
    status,
    running,
    mailAddress,
    smtpAddress,
    imapAddress,
    databasePath,
    peers,
    startService,
    stopService,
    restartService,
    refreshAll,
    isStarting,
    isStopping,
  } = useServiceStatus({
    refreshInterval: 5000,
    fetchOnMount: true,
  });

  const serviceStatusValue: ServiceStatus = running ? 'Running' : 'Stopped';

  // Copy to clipboard with feedback
  const handleCopy = async (text: string, fieldName: string) => {
    try {
      await CopyToClipboard(text);
      setCopiedField(fieldName);
      showSuccess(t('dashboard.messages.copied'), t('dashboard.messages.copiedToClipboard').replace('{field}', fieldName));
      setTimeout(() => setCopiedField(null), 2000);
    } catch (error) {
      showError(t('dashboard.messages.copyFailed'), error instanceof Error ? error.message : t('dashboard.messages.copyFailedMessage'));
    }
  };

  // Open DeltaChat with dclogin:// URL for auto-configuration
  const handleOpenDeltaChat = async () => {
    if (!mailAddress) {
      showError(t('dashboard.messages.noMailAddress'), t('dashboard.messages.noMailAddressMessage'));
      return;
    }
    try {
      await OpenDeltaChat();
      showSuccess(t('dashboard.messages.deltachatOpened'), t('dashboard.messages.deltachatOpenedMessage'));
    } catch (error) {
      showError(t('dashboard.messages.openFailed'), error instanceof Error ? error.message : t('dashboard.messages.openFailedMessage'));
    }
  };

  // Service control actions
  const handleStartService = async () => {
    try {
      await startService();
      showSuccess(t('dashboard.messages.serviceStarted'), t('dashboard.messages.serviceStartedMessage'));
    } catch (error) {
      showError(t('dashboard.messages.startFailed'), error instanceof Error ? error.message : t('dashboard.messages.startFailedMessage'));
    }
  };

  const handleStopService = async () => {
    try {
      await stopService();
      showSuccess(t('dashboard.messages.serviceStopped'), t('dashboard.messages.serviceStoppedMessage'));
    } catch (error) {
      showError(t('dashboard.messages.stopFailed'), error instanceof Error ? error.message : t('dashboard.messages.stopFailedMessage'));
    }
  };

  const handleRestartService = async () => {
    try {
      await restartService();
      showSuccess(t('dashboard.messages.serviceRestarted'), t('dashboard.messages.serviceRestartedMessage'));
    } catch (error) {
      showError(t('dashboard.messages.restartFailed'), error instanceof Error ? error.message : t('dashboard.messages.restartFailedMessage'));
    }
  };

  return (
    <div className="space-y-4 md:space-y-6 pb-6">
      {/* Header */}
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, type: 'spring' }}
        className="text-center space-y-2"
      >
        <h1 className="text-5xl font-display font-black text-transparent bg-clip-text bg-iridescent bg-[length:200%_100%] animate-holographic-spin">
          {t('dashboard.title')}
        </h1>
        <p className="text-lg text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant font-futuristic">
          {t('dashboard.subtitle')}
        </p>
      </motion.div>

      {/* Main Dashboard - Vertical Stack */}
      <div className="space-y-4 md:space-y-6">
        {/* Service Status Card with Server Info */}
        <motion.div
          initial={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.5, delay: 0.1 }}
        >
          <GlassCard
            title={t('dashboard.serviceStatus')}
            accentColor="#006C4C"
            hoverable
            headerAction={
              <StatusIndicator status={serviceStatusValue} animated size="lg" />
            }
          >
            <div className="space-y-6">
              {/* Control Buttons */}
              <div className="flex flex-wrap gap-3">
                {running ? (
                  <>
                    <Button
                      variant="danger"
                      glow
                      onClick={handleStopService}
                      disabled={isStopping}
                      className="flex-1 min-w-[140px]"
                    >
                      {isStopping ? t('dashboard.stopping') : t('dashboard.stopService')}
                    </Button>
                    <Button
                      variant="secondary"
                      glow
                      onClick={handleRestartService}
                      disabled={isStopping || isStarting}
                      className="flex-1 min-w-[140px]"
                    >
                      {t('dashboard.restartService')}
                    </Button>
                  </>
                ) : (
                  <Button
                    variant="primary"
                    glow
                    onClick={handleStartService}
                    disabled={isStarting}
                    className="w-full"
                  >
                    {isStarting ? t('dashboard.starting') : t('dashboard.startService')}
                  </Button>
                )}
              </div>

              {/* Server Info Section */}
              <div className="pt-4 border-t border-md-light-outline/30 dark:border-md-dark-outline/30">
                <h3 className="text-sm font-bold text-md-light-onSurface dark:text-md-dark-onSurface font-futuristic uppercase tracking-wide mb-4">{t('dashboard.serverConfiguration')}</h3>

                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  {/* SMTP Server */}
                  <div>
                    <p className="text-xs text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant mb-1 font-futuristic uppercase tracking-wide">{t('dashboard.smtpServer')}</p>
                    <div className="bg-md-light-surfaceVariant dark:bg-md-dark-surfaceVariant rounded px-2 py-2 border border-md-light-outline/30 dark:border-md-dark-outline/30 overflow-x-auto">
                      <p className="text-xs font-mono text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant break-all">
                        {smtpAddress || t('dashboard.notConfigured')}
                      </p>
                    </div>
                    {smtpAddress && (
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleCopy(smtpAddress, t('dashboard.messages.smtpAddress'))}
                        className="mt-1 w-full"
                      >
                        {copiedField === t('dashboard.messages.smtpAddress') ? t('dashboard.copied') : t('dashboard.copy')}
                      </Button>
                    )}
                  </div>

                  {/* IMAP Server */}
                  <div>
                    <p className="text-xs text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant mb-1 font-futuristic uppercase tracking-wide">{t('dashboard.imapServer')}</p>
                    <div className="bg-md-light-surfaceVariant dark:bg-md-dark-surfaceVariant rounded px-2 py-2 border border-md-light-outline/30 dark:border-md-dark-outline/30 overflow-x-auto">
                      <p className="text-xs font-mono text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant break-all">
                        {imapAddress || t('dashboard.notConfigured')}
                      </p>
                    </div>
                    {imapAddress && (
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleCopy(imapAddress, t('dashboard.messages.imapAddress'))}
                        className="mt-1 w-full"
                      >
                        {copiedField === t('dashboard.messages.imapAddress') ? t('dashboard.copied') : t('dashboard.copy')}
                      </Button>
                    )}
                  </div>

                  {/* Database */}
                  <div>
                    <p className="text-xs text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant mb-1 font-futuristic uppercase tracking-wide">{t('dashboard.database')}</p>
                    <div className="bg-md-light-surfaceVariant dark:bg-md-dark-surfaceVariant rounded px-2 py-2 border border-md-light-outline/30 dark:border-md-dark-outline/30 overflow-x-auto">
                      <p className="text-xs font-mono text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant break-all">
                        {databasePath || t('dashboard.notConfigured')}
                      </p>
                    </div>
                    {databasePath && (
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleCopy(databasePath, t('dashboard.messages.databasePath'))}
                        className="mt-1 w-full"
                      >
                        {copiedField === t('dashboard.messages.databasePath') ? t('dashboard.copied') : t('dashboard.copy')}
                      </Button>
                    )}
                  </div>
                </div>
              </div>
            </div>
          </GlassCard>
        </motion.div>

        {/* Mail Address Card */}
        <motion.div
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.5, delay: 0.2 }}
        >
          <HolographicBorder animated borderWidth={2}>
            <GlassCard
              title={t('dashboard.mailAddress')}
              subtitle={t('dashboard.mailAddressSubtitle')}
              padding="lg"
            >
              <div className="space-y-4">
                {mailAddress ? (
                  <>
                    <div className="bg-md-light-primaryContainer dark:bg-md-dark-primaryContainer rounded-lg p-3 md:p-4 border border-md-light-primary/30 dark:border-md-dark-primary/30 overflow-x-auto">
                      <p className="text-xs md:text-sm font-mono text-md-light-onPrimaryContainer dark:text-md-dark-onPrimaryContainer break-all">
                        {mailAddress}
                      </p>
                    </div>
                    <div className="flex gap-3">
                      <Button
                        variant="primary"
                        glow
                        onClick={() => handleCopy(mailAddress, t('dashboard.messages.smtpAddress'))}
                        disabled={copiedField === t('dashboard.messages.smtpAddress')}
                        className="flex-1"
                      >
                        {copiedField === t('dashboard.messages.smtpAddress') ? t('dashboard.copied') : t('dashboard.copyAddressButton')}
                      </Button>
                      <Button
                        variant="secondary"
                        glow
                        onClick={handleOpenDeltaChat}
                        className="flex-1"
                      >
                        {t('dashboard.openInDeltaChat')}
                      </Button>
                    </div>
                  </>
                ) : (
                  <div className="text-center py-6 text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant">
                    <p>{t('dashboard.mailNotAvailable')}</p>
                    <p className="text-sm mt-2">{t('dashboard.startServicePrompt')}</p>
                  </div>
                )}
              </div>
            </GlassCard>
          </HolographicBorder>
        </motion.div>

        {/* Connected Peers Card */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.3 }}
        >
          <HolographicBorder animated borderWidth={2}>
            <GlassCard
              title={t('dashboard.connectedPeers')}
              padding="lg"
            >
              {/* Active Peers List - Only show connected peers */}
              {peers.filter(p => p.connected).length > 0 ? (
                <div className="space-y-3">
                  {peers
                    .filter(p => p.connected)
                    .map((peer, index) => (
                      <motion.div
                        key={peer.address}
                        initial={{ opacity: 0, scale: 0.9 }}
                        animate={{ opacity: 1, scale: 1 }}
                        transition={{ duration: 0.3, delay: 0.4 + index * 0.05 }}
                      >
                        <PeerCard
                          peer={{
                            address: peer.address,
                            connected: peer.connected,
                            latency: peer.latency,
                            rxBytes: peer.rxBytes,
                            txBytes: peer.txBytes,
                            uptime: peer.uptime,
                          }}
                          showActions={false}
                          variant="compact"
                        />
                      </motion.div>
                    ))}
                </div>
              ) : (
                <div className="text-center py-8 text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant space-y-2">
                  <div className="text-5xl mb-3">🌐</div>
                  <p className="text-base">{t('dashboard.noPeersConfigured')}</p>
                  <p className="text-xs">{t('dashboard.addPeersPrompt')}</p>
                  <Button
                    variant="primary"
                    glow
                    onClick={() => navigate('/peers')}
                    className="mt-3"
                  >
                    {t('dashboard.goToSettings')}
                  </Button>
                </div>
              )}
            </GlassCard>
          </HolographicBorder>
        </motion.div>
      </div>
    </div>
  );
}

export default Dashboard;
