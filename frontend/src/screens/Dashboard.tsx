import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { motion } from 'framer-motion';
import {
  Button,
  GlassCard,
  StatusIndicator,
  PeerCard,
} from '../components';
import { useServiceStatus } from '../hooks/useServiceStatus';
import { useI18n } from '../hooks/useI18n';
import { CopyToClipboard, OpenDeltaChat, GetStorageStats } from '../../wailsjs/go/main/App';
import { toast } from '../components/ui/Toast';
import type { ServiceStatus } from '../components';

/**
 * Dashboard Screen - Main application screen
 */
export function Dashboard() {
  const { t } = useI18n();
  const navigate = useNavigate();
  const [copiedField, setCopiedField] = useState<string | null>(null);
  const [storageStats, setStorageStats] = useState<any>(null);
  const [showDeltaChat, setShowDeltaChat] = useState(false);
  const [showEmailClients, setShowEmailClients] = useState(false);

  useEffect(() => {
    const loadStorageStats = async () => {
      try {
        const stats = await GetStorageStats();
        setStorageStats(stats);
      } catch (error) {
        console.error('Failed to load storage stats:', error);
      }
    };
    loadStorageStats();
  }, []);

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
    isRestarting,
  } = useServiceStatus({
    refreshInterval: 5000,
    fetchOnMount: true,
  });

  const serviceStatusValue: ServiceStatus = running ? 'Running' : 'Stopped';

  const handleCopy = async (text: string, fieldName: string) => {
    try {
      await CopyToClipboard(text);
      setCopiedField(fieldName);
      toast.success(t('dashboard.messages.copiedToClipboard'));
      setTimeout(() => setCopiedField(null), 2000);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : t('dashboard.messages.copyFailedMessage'));
    }
  };

  const handleOpenDeltaChat = async () => {
    if (!mailAddress) {
      toast.error(t('dashboard.messages.noMailAddressMessage'));
      return;
    }
    try {
      await OpenDeltaChat();
      toast.success(t('dashboard.messages.deltachatOpenedMessage'));
    } catch (error) {
      toast.error(error instanceof Error ? error.message : t('dashboard.messages.openFailedMessage'));
    }
  };

  const handleStartService = async () => {
    try {
      await startService();
      toast.success(t('dashboard.messages.serviceStartedMessage'));
    } catch (error) {
      toast.error(error instanceof Error ? error.message : t('dashboard.messages.startFailedMessage'));
    }
  };

  const handleStopService = async () => {
    try {
      await stopService();
      toast.success(t('dashboard.messages.serviceStoppedMessage'));
    } catch (error) {
      toast.error(error instanceof Error ? error.message : t('dashboard.messages.stopFailedMessage'));
    }
  };

  const handleRestartService = async () => {
    try {
      await restartService();
      toast.success(t('dashboard.messages.serviceRestartedMessage'));
    } catch (error) {
      toast.error(error instanceof Error ? error.message : t('dashboard.messages.restartFailedMessage'));
    }
  };

  return (
    <div className="space-y-6 pb-6">
      {/* Header */}
      <motion.div
        initial={{ opacity: 0, y: -10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.2 }}
      >
        <h1 className="text-2xl font-semibold text-slate-100">
          {t('dashboard.title')}
        </h1>
        <p className="text-sm text-slate-400 mt-1">
          {t('dashboard.subtitle')}
        </p>
      </motion.div>

      {/* Main Dashboard Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Service Status Card */}
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.2, delay: 0.05 }}
          className="lg:col-span-2"
        >
          <GlassCard
            title={t('dashboard.serviceStatus')}
            accentColor="#10b981"
            headerAction={
              <StatusIndicator status={serviceStatusValue} size="md" />
            }
          >
            <div className="space-y-5">
              {/* Control Buttons */}
              <div className="flex flex-wrap gap-3">
                {running ? (
                  <>
                    <Button
                      variant="danger"
                      onClick={handleStopService}
                      disabled={isStopping || isRestarting}
                      className="flex-1 min-w-[140px]"
                    >
                      {isStopping ? t('dashboard.stopping') : t('dashboard.stopService')}
                    </Button>
                    <Button
                      variant="secondary"
                      onClick={handleRestartService}
                      disabled={isStopping || isStarting || isRestarting}
                      className="flex-1 min-w-[140px]"
                    >
                      {isRestarting ? t('dashboard.restarting') : t('dashboard.restartService')}
                    </Button>
                  </>
                ) : (
                  <Button
                    variant="primary"
                    onClick={handleStartService}
                    disabled={isStarting || isRestarting}
                    className="w-full"
                  >
                    {isStarting ? t('dashboard.starting') : t('dashboard.startService')}
                  </Button>
                )}
              </div>

              {/* Server Info Section */}
              <div className="pt-4 border-t border-slate-700">
                <h3 className="text-xs font-medium text-slate-400 uppercase tracking-wide mb-3">
                  {t('dashboard.serverConfiguration')}
                </h3>

                <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
                  {/* SMTP Server */}
                  <div>
                    <p className="text-xs text-slate-400 mb-1">{t('dashboard.smtpServer')}</p>
                    <div className="bg-slate-700 rounded-lg px-3 py-2">
                      <p className="text-xs font-mono text-slate-300 break-all">
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
                    <p className="text-xs text-slate-400 mb-1">{t('dashboard.imapServer')}</p>
                    <div className="bg-slate-700 rounded-lg px-3 py-2">
                      <p className="text-xs font-mono text-slate-300 break-all">
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
                    <p className="text-xs text-slate-400 mb-1">{t('dashboard.database')}</p>
                    <div className="bg-slate-700 rounded-lg px-3 py-2">
                      <p className="text-xs font-mono text-slate-300 break-all">
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
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.2, delay: 0.1 }}
        >
          <GlassCard
            title={t('dashboard.mailAddress')}
            subtitle={t('dashboard.mailAddressSubtitle')}
          >
            <div className="space-y-4">
              {mailAddress ? (
                <>
                  <div className="bg-emerald-900/20 border border-emerald-500/30 rounded-lg p-3">
                    <p className="text-sm font-mono text-emerald-400 break-all">
                      {mailAddress}
                    </p>
                  </div>
                  <div className="flex gap-3">
                    <Button
                      variant="primary"
                      onClick={() => handleCopy(mailAddress, t('dashboard.messages.mailAddress'))}
                      disabled={copiedField === t('dashboard.messages.mailAddress')}
                      className="flex-1"
                    >
                      {copiedField === t('dashboard.messages.mailAddress') ? t('dashboard.copied') : t('dashboard.copyAddressButton')}
                    </Button>
                    <Button
                      variant="secondary"
                      onClick={handleOpenDeltaChat}
                      className="flex-1"
                    >
                      {t('dashboard.openInDeltaChat')}
                    </Button>
                  </div>
                </>
              ) : (
                <div className="text-center py-6 text-slate-400">
                  <p>{t('dashboard.mailNotAvailable')}</p>
                  <p className="text-sm mt-2">{t('dashboard.startServicePrompt')}</p>
                </div>
              )}
            </div>
          </GlassCard>
        </motion.div>

        {/* Storage Card */}
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.2, delay: 0.15 }}
        >
          <GlassCard title={t('storage.title')}>
            <div className="grid grid-cols-2 gap-3">
              <div className="bg-slate-700 rounded-lg px-3 py-2">
                <p className="text-xs text-slate-400 mb-1">{t('storage.maxMessageSize')}</p>
                <p className="text-sm font-semibold text-emerald-400">
                  {storageStats?.maxMessageSizeMB || 10} {t('storage.mb')}
                </p>
              </div>
              <div className="bg-slate-700 rounded-lg px-3 py-2">
                <p className="text-xs text-slate-400 mb-1">{t('storage.databaseSize')}</p>
                <p className="text-sm font-semibold text-slate-200">
                  {storageStats?.databaseSizeMB?.toFixed(2) || '0.00'} {t('storage.mb')}
                </p>
              </div>
              <div className="bg-slate-700 rounded-lg px-3 py-2">
                <p className="text-xs text-slate-400 mb-1">{t('storage.filesSize')}</p>
                <p className="text-sm font-semibold text-slate-200">
                  {storageStats?.filesSizeMB?.toFixed(2) || '0.00'} {t('storage.mb')}
                </p>
              </div>
              <div className="bg-emerald-900/30 border border-emerald-500/20 rounded-lg px-3 py-2">
                <p className="text-xs text-emerald-400/80 mb-1">{t('storage.totalSize')}</p>
                <p className="text-sm font-semibold text-emerald-400">
                  {storageStats?.totalSizeMB?.toFixed(2) || '0.00'} {t('storage.mb')}
                </p>
              </div>
            </div>
          </GlassCard>
        </motion.div>

        {/* Connected Peers Card */}
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.2, delay: 0.2 }}
          className="lg:col-span-2"
        >
          <GlassCard title={t('dashboard.connectedPeers')}>
            {peers.filter(p => p.connected).length > 0 ? (
              <div className="space-y-2">
                {peers
                  .filter(p => p.connected)
                  .map((peer, index) => (
                    <motion.div
                      key={peer.address}
                      initial={{ opacity: 0 }}
                      animate={{ opacity: 1 }}
                      transition={{ duration: 0.2, delay: index * 0.03 }}
                    >
                      <PeerCard
                        peer={{
                          address: peer.address,
                          enabled: true,
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
              <div className="text-center py-8 text-slate-400">
                <div className="text-4xl mb-3">🌐</div>
                <p>{t('dashboard.noPeersConfigured')}</p>
                <p className="text-xs mt-1">{t('dashboard.addPeersPrompt')}</p>
                <Button
                  variant="primary"
                  onClick={() => navigate('/peers')}
                  className="mt-4"
                >
                  {t('dashboard.goToSettings')}
                </Button>
              </div>
            )}
          </GlassCard>
        </motion.div>

        {/* DeltaChat Setup Card */}
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.2, delay: 0.25 }}
        >
          <GlassCard
            title={t('deltachat.title')}
            subtitle={t('deltachat.subtitle')}
          >
            <div className="space-y-3">
              <Button
                variant="secondary"
                onClick={() => setShowDeltaChat(!showDeltaChat)}
                className="w-full"
              >
                {showDeltaChat ? `▼ ${t('deltachat.hideInstructions')}` : `▶ ${t('deltachat.showInstructions')}`}
              </Button>

              <motion.div
                initial={false}
                animate={{ height: showDeltaChat ? 'auto' : 0, opacity: showDeltaChat ? 1 : 0 }}
                transition={{ duration: 0.2 }}
                style={{ overflow: 'hidden' }}
              >
                <div className="space-y-3 pt-1">
                  <div className="bg-slate-700 rounded-lg p-4">
                    <h4 className="text-sm font-semibold text-emerald-400 mb-2">
                      {t('deltachat.automatic.title')}
                    </h4>
                    <ol className="space-y-1 text-sm text-slate-300 list-decimal list-inside">
                      <li>{t('deltachat.automatic.step1')}</li>
                      <li>{t('deltachat.automatic.step2')}</li>
                      <li>{t('deltachat.automatic.step3')}</li>
                      <li>{t('deltachat.automatic.step4')}</li>
                    </ol>
                  </div>

                  <div className="bg-slate-700 rounded-lg p-4">
                    <h4 className="text-sm font-semibold text-emerald-400 mb-2">
                      {t('deltachat.manual.title')}
                    </h4>
                    <ol className="space-y-1 text-sm text-slate-300 list-decimal list-inside">
                      <li>{t('deltachat.manual.step1')}</li>
                      <li>{t('deltachat.manual.step2')}</li>
                      <li>{t('deltachat.manual.step3')}</li>
                      <li>{t('deltachat.manual.step4')}</li>
                      <li>{t('deltachat.manual.step5')}</li>
                    </ol>
                  </div>
                </div>
              </motion.div>
            </div>
          </GlassCard>
        </motion.div>

        {/* Email Clients Setup Card */}
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.2, delay: 0.3 }}
        >
          <GlassCard
            title={t('emailClients.title')}
            subtitle={t('emailClients.subtitle')}
          >
            <div className="space-y-3">
              <Button
                variant="secondary"
                onClick={() => setShowEmailClients(!showEmailClients)}
                className="w-full"
              >
                {showEmailClients ? `▼ ${t('emailClients.hideInstructions')}` : `▶ ${t('emailClients.showInstructions')}`}
              </Button>

              <motion.div
                initial={false}
                animate={{ height: showEmailClients ? 'auto' : 0, opacity: showEmailClients ? 1 : 0 }}
                transition={{ duration: 0.2 }}
                style={{ overflow: 'hidden' }}
              >
                <div className="space-y-3 pt-1">
                  {/* Server Configuration */}
                  <div className="bg-emerald-900/20 border border-emerald-500/20 rounded-lg p-4">
                    <h4 className="text-sm font-semibold text-emerald-400 mb-2">
                      {t('emailClients.serverConfig.title')}
                    </h4>
                    <div className="grid grid-cols-2 gap-2 text-sm">
                      <div>
                        <p className="text-emerald-400/70 text-xs">{t('emailClients.serverConfig.imap')}</p>
                        <p className="font-mono text-emerald-300">127.0.0.1:1143</p>
                      </div>
                      <div>
                        <p className="text-emerald-400/70 text-xs">{t('emailClients.serverConfig.smtp')}</p>
                        <p className="font-mono text-emerald-300">127.0.0.1:1025</p>
                      </div>
                      <div>
                        <p className="text-emerald-400/70 text-xs">{t('emailClients.serverConfig.encryption')}</p>
                        <p className="font-mono text-emerald-300">{t('emailClients.serverConfig.noEncryption')}</p>
                      </div>
                      <div>
                        <p className="text-emerald-400/70 text-xs">{t('emailClients.serverConfig.password')}</p>
                        <p className="font-mono text-emerald-300">{t('emailClients.serverConfig.tyrPassword')}</p>
                      </div>
                    </div>
                  </div>

                  {/* Thunderbird */}
                  <div className="bg-slate-700 rounded-lg p-4">
                    <h4 className="text-sm font-semibold text-emerald-400 mb-2">
                      📧 {t('emailClients.thunderbird.title')}
                    </h4>
                    <ol className="space-y-1 text-sm text-slate-300 list-decimal list-inside">
                      <li>{t('emailClients.thunderbird.step1')}</li>
                      <li>{t('emailClients.thunderbird.step2')}</li>
                      <li>{t('emailClients.thunderbird.step3')}</li>
                      <li>{t('emailClients.thunderbird.step4')}</li>
                      <li>{t('emailClients.thunderbird.step5')}</li>
                      <li>{t('emailClients.thunderbird.step6')}</li>
                    </ol>
                  </div>

                  {/* Mailspring */}
                  <div className="bg-slate-700 rounded-lg p-4">
                    <h4 className="text-sm font-semibold text-emerald-400 mb-2">
                      📧 {t('emailClients.mailspring.title')}
                    </h4>
                    <ol className="space-y-1 text-sm text-slate-300 list-decimal list-inside">
                      <li>{t('emailClients.mailspring.step1')}</li>
                      <li>{t('emailClients.mailspring.step2')}</li>
                      <li>{t('emailClients.mailspring.step3')}</li>
                      <li>{t('emailClients.mailspring.step4')}</li>
                      <li>{t('emailClients.mailspring.step5')}</li>
                    </ol>
                  </div>

                  {/* Apple Mail */}
                  <div className="bg-slate-700 rounded-lg p-4">
                    <h4 className="text-sm font-semibold text-emerald-400 mb-2">
                      📧 {t('emailClients.appleMail.title')}
                    </h4>
                    <ol className="space-y-1 text-sm text-slate-300 list-decimal list-inside">
                      <li>{t('emailClients.appleMail.step1')}</li>
                      <li>{t('emailClients.appleMail.step2')}</li>
                      <li>{t('emailClients.appleMail.step3')}</li>
                      <li>{t('emailClients.appleMail.step4')}</li>
                      <li>{t('emailClients.appleMail.step5')}</li>
                    </ol>
                  </div>
                </div>
              </motion.div>
            </div>
          </GlassCard>
        </motion.div>
      </div>
    </div>
  );
}

export default Dashboard;
