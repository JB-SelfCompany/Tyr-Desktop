import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { useTranslation } from 'react-i18next';
import {
  Button,
  Input,
  GlassCard,
} from '../components';
import {
  SetOnboardingComplete,
  SetPassword,
  AddPeer,
  SetLanguage,
  SetTheme,
  GetSystemLanguage,
  GetDefaultPeers,
  FindAvailablePeers,
  GetCachedDiscoveredPeers,
  AddDiscoveredPeers,
  RestoreBackup,
  ShowOpenFileDialog,
  OnStartupComplete,
  GetConfig,
  RemovePeer,
} from '../../wailsjs/go/main/App';
import { LogPrint, EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';
import { toast } from '../components/ui/Toast';
import { core, models } from '../../wailsjs/go/models';

type DiscoveredPeer = core.DiscoveredPeer;
type RestoreOptionsDTO = models.RestoreOptionsDTO;

type OnboardingStep = 1 | 2 | 3;

/**
 * Onboarding Screen - First-time setup wizard
 */
interface OnboardingProps {
  onComplete: () => void | Promise<void>;
}

export function Onboarding({ onComplete }: OnboardingProps) {
  const { t, i18n } = useTranslation();
  const [step, setStep] = useState<OnboardingStep>(1);
  const [isProcessing, setIsProcessing] = useState(false);

  // Password form state
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [passwordError, setPasswordError] = useState('');

  // Peer configuration state
  const [useDefaultPeers, setUseDefaultPeers] = useState(true);
  const [selectedPeers, setSelectedPeers] = useState<DiscoveredPeer[]>([]);
  const [discoveredPeers, setDiscoveredPeers] = useState<DiscoveredPeer[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [searchProgress, setSearchProgress] = useState({ current: 0, total: 0, available_count: 0 });
  const [hasCachedPeers, setHasCachedPeers] = useState(false);

  // Peer discovery filters
  const [selectedProtocols, setSelectedProtocols] = useState<string[]>(['tls', 'quic']);
  const [maxRTT, setMaxRTT] = useState<number>(500);

  // Backup restore state
  const [showRestoreDialog, setShowRestoreDialog] = useState(false);
  const [backupFilePath, setBackupFilePath] = useState('');
  const [restorePassword, setRestorePassword] = useState('');

  const totalSteps = 3;

  // Auto-detect system language and theme on mount
  useEffect(() => {
    const initializeSettings = async () => {
      try {
        // Detect and set system language
        const systemLang = await GetSystemLanguage();
        LogPrint(`[Onboarding] Detected system language: ${systemLang}`);
        await SetLanguage(systemLang);
        await i18n.changeLanguage(systemLang);

        // Set dark theme (dark-only mode)
        await SetTheme('dark');
      } catch (error) {
        LogPrint(`[Onboarding] Error initializing settings: ${error}`);
      }
    };

    initializeSettings();
  }, [i18n]);

  // Load cached peers on mount
  useEffect(() => {
    const loadCachedPeers = async () => {
      try {
        const cached = await GetCachedDiscoveredPeers();
        if (cached && cached.length > 0) {
          LogPrint(`[Onboarding] Loaded ${cached.length} cached peers`);
          setDiscoveredPeers(cached);
          setHasCachedPeers(true);
        }
      } catch (error) {
        LogPrint(`[Onboarding] Error loading cached peers: ${error}`);
      }
    };

    loadCachedPeers();
  }, []);

  // Validate password
  const validatePassword = (): boolean => {
    if (password.length === 0) {
      setPasswordError(t('onboarding.error.passwordEmpty'));
      return false;
    }
    if (password.length < 6) {
      setPasswordError(t('onboarding.error.passwordShort'));
      return false;
    }
    if (password !== confirmPassword) {
      setPasswordError(t('onboarding.error.mismatch'));
      return false;
    }
    setPasswordError('');
    return true;
  };

  // Validate peer selection
  const validatePeers = (): boolean => {
    if (!useDefaultPeers && selectedPeers.length === 0) {
      toast.error(t('onboarding.error.noPeersSelected'));
      return false;
    }
    return true;
  };

  // Handle peer discovery
  const handleFindPeers = async () => {
    if (selectedProtocols.length === 0) {
      toast.error(t('peers.discovery.noProtocolsSelected'));
      return;
    }

    setIsSearching(true);
    setSearchProgress({ current: 0, total: 0, available_count: 0 });
    setUseDefaultPeers(false); // Disable default peers when searching

    try {
      const result = await FindAvailablePeers(
        selectedProtocols.join(','),
        '', // All regions
        maxRTT
      );

      if (result && result.peers) {
        LogPrint(`[Onboarding] Found ${result.available} peers out of ${result.total}`);
        setDiscoveredPeers(result.peers);
        setHasCachedPeers(false); // Not from cache
        toast.success(
          t('peers.discovery.foundPeers')
            .replace('{{count}}', String(result.available))
            .replace('{{total}}', String(result.total))
        );
      } else {
        toast.error(t('onboarding.discovery.noPeersFound'));
      }
    } catch (error) {
      LogPrint(`[Onboarding] Peer discovery error: ${error}`);
      toast.error(error instanceof Error ? error.message : t('peers.discovery.searchFailed'));
    } finally {
      setIsSearching(false);
      setSearchProgress({ current: 0, total: 0, available_count: 0 });
    }
  };

  // Subscribe to peer discovery progress
  useEffect(() => {
    const unsubscribe = EventsOn('peer:discovery:progress', (data: { current: number; total: number; available_count: number }) => {
      setSearchProgress(data);
    });

    return () => {
      EventsOff('peer:discovery:progress');
      if (unsubscribe) unsubscribe();
    };
  }, []);

  // Handle backup file selection
  const handleSelectBackupFile = async () => {
    try {
      const filePath = await ShowOpenFileDialog(t('onboarding.messages.selectBackupFile'));
      if (filePath) {
        setBackupFilePath(filePath);
      }
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : t('onboarding.messages.fileSelectionFailed')
      );
    }
  };

  // Handle restore from backup
  const handleRestoreBackup = async () => {
    if (!backupFilePath) {
      toast.error(t('onboarding.messages.fileRequiredMessage'));
      return;
    }

    if (!restorePassword) {
      toast.error(t('onboarding.messages.passwordRequiredMessage'));
      return;
    }

    setIsProcessing(true);
    try {
      const options: RestoreOptionsDTO = {
        backupPath: backupFilePath,
        password: restorePassword,
      };
      const result = await RestoreBackup(options);

      // Check if restore was successful (e.g., wrong password)
      if (!result.success) {
        throw new Error(result.message || t('onboarding.messages.restoreFailedMessage'));
      }

      LogPrint('[Onboarding] Backup restored successfully');

      // Initialize service manager after restore
      await OnStartupComplete();

      toast.success(t('onboarding.messages.setupCompleteMessage'));
      await onComplete();
    } catch (error) {
      LogPrint(`[Onboarding] Restore error: ${error}`);
      toast.error(
        error instanceof Error ? error.message : t('onboarding.messages.setupFailedMessage')
      );
    } finally {
      setIsProcessing(false);
    }
  };

  // Handle next step
  const handleNext = async () => {
    // Step 2: Validate password
    if (step === 2) {
      if (!validatePassword()) {
        return;
      }
    }

    // Step 3: Complete setup
    if (step === 3) {
      if (!validatePeers()) {
        return;
      }

      setIsProcessing(true);
      try {
        // Set password
        LogPrint('[Onboarding] Setting password...');
        await SetPassword(password);

        // CRITICAL: Remove all existing peers first (including default peers from config init)
        LogPrint('[Onboarding] Removing all existing peers before adding selected ones...');
        const currentConfig = await GetConfig();
        if (currentConfig.peers && currentConfig.peers.length > 0) {
          LogPrint(`[Onboarding] Found ${currentConfig.peers.length} existing peers to remove`);
          for (const peer of currentConfig.peers) {
            LogPrint(`[Onboarding] Removing existing peer: ${peer.address}`);
            await RemovePeer(peer.address);
          }
        }
        LogPrint('[Onboarding] All existing peers removed');

        // Add peers - STRICTLY one or the other, never both
        if (useDefaultPeers) {
          LogPrint('[Onboarding] Using default peers ONLY');
          const defaultPeers = await GetDefaultPeers();
          LogPrint(`[Onboarding] Adding ${defaultPeers.length} default peers`);
          for (const peer of defaultPeers) {
            LogPrint(`[Onboarding] Adding default peer: ${peer}`);
            await AddPeer(peer);
          }
        } else if (selectedPeers.length > 0) {
          // Only add discovered peers if we have selections AND default peers toggle is OFF
          LogPrint(`[Onboarding] Adding ${selectedPeers.length} discovered peers ONLY (no defaults)`);
          LogPrint(`[Onboarding] useDefaultPeers = ${useDefaultPeers}`);
          for (const peer of selectedPeers) {
            LogPrint(`[Onboarding] Adding discovered peer: ${peer.address}`);
          }
          await AddDiscoveredPeers(selectedPeers);
        } else {
          LogPrint('[Onboarding] No peers to add (useDefaultPeers=false, selectedPeers=0)');
        }

        // Mark onboarding complete
        LogPrint('[Onboarding] Marking onboarding complete...');
        await SetOnboardingComplete();

        // Initialize service manager
        LogPrint('[Onboarding] Initializing service...');
        await OnStartupComplete();

        toast.success(t('onboarding.messages.setupCompleteMessage'));
        await onComplete();
      } catch (error) {
        LogPrint(`[Onboarding] Setup error: ${error}`);
        toast.error(
          error instanceof Error ? error.message : t('onboarding.messages.setupFailedMessage')
        );
      } finally {
        setIsProcessing(false);
      }
      return;
    }

    // Move to next step
    setStep((step + 1) as OnboardingStep);
  };

  // Handle back
  const handleBack = () => {
    if (step > 1) {
      setStep((step - 1) as OnboardingStep);
    }
  };

  // Toggle peer selection
  const togglePeerSelection = (peer: DiscoveredPeer) => {
    const isSelected = selectedPeers.some(p => p.address === peer.address);
    if (isSelected) {
      setSelectedPeers(selectedPeers.filter(p => p.address !== peer.address));
    } else {
      setSelectedPeers([...selectedPeers, peer]);
    }
  };

  // Toggle protocol selection
  const toggleProtocol = (protocol: string) => {
    setSelectedProtocols(prev =>
      prev.includes(protocol)
        ? prev.filter(p => p !== protocol)
        : [...prev, protocol]
    );
  };

  const protocols = [
    { value: 'tcp', label: 'TCP' },
    { value: 'tls', label: 'TLS' },
    { value: 'quic', label: 'QUIC' },
    { value: 'ws', label: 'WebSocket' },
    { value: 'wss', label: 'WebSocket (TLS)' },
  ];

  // Render step content
  const renderStepContent = () => {
    // Step 1: Welcome
    if (step === 1) {
      return (
        <motion.div
          key="welcome"
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          exit={{ opacity: 0, x: -20 }}
          className="text-center space-y-4"
        >
          <div className="text-5xl sm:text-6xl mb-3">🚀</div>
          <h2 className="text-2xl sm:text-3xl font-semibold text-slate-100">
            {t('onboarding.welcomeTitle')}
          </h2>
          <p className="text-sm text-slate-400 max-w-xl mx-auto px-2">
            {t('onboarding.welcomeDescription')}
          </p>

          {/* Feature chips */}
          <div className="flex flex-wrap justify-center gap-2 pt-2">
            <div className="px-3 py-1.5 bg-emerald-500/20 border border-emerald-500/30 rounded-full">
              <span className="text-xs font-medium text-emerald-400">{t('onboarding.features.decentralized')}</span>
            </div>
            <div className="px-3 py-1.5 bg-blue-500/20 border border-blue-500/30 rounded-full">
              <span className="text-xs font-medium text-blue-400">{t('onboarding.features.encrypted')}</span>
            </div>
            <div className="px-3 py-1.5 bg-amber-500/20 border border-amber-500/30 rounded-full">
              <span className="text-xs font-medium text-amber-400">{t('onboarding.features.p2p')}</span>
            </div>
          </div>

          <div className="pt-4">
            <Button variant="primary" onClick={handleNext}>
              {t('onboarding.getStarted')}
            </Button>
          </div>
        </motion.div>
      );
    }

    // Step 2: Password Setup
    if (step === 2) {
      return (
        <motion.div
          key="password"
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          exit={{ opacity: 0, x: -20 }}
          className="space-y-3 w-full"
        >
          <div className="text-center space-y-1">
            <div className="text-3xl mb-2">🔐</div>
            <h2 className="text-xl sm:text-2xl font-semibold text-slate-100">
              {t('onboarding.password.title')}
            </h2>
            <p className="text-xs text-slate-400">
              {t('onboarding.password.description')}
            </p>
          </div>

          <div className="space-y-2.5">
            <Input
              label={t('onboarding.password.label')}
              type="password"
              placeholder={t('onboarding.password.placeholder')}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              error={passwordError}
            />
            <Input
              label={t('onboarding.password.confirmLabel')}
              type="password"
              placeholder={t('onboarding.password.confirmPlaceholder')}
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
            />
            <div className="bg-emerald-500/10 border border-emerald-500/20 rounded-lg p-2">
              <p className="text-xs text-slate-300">
                {t('onboarding.password.securityInfo')}
              </p>
            </div>
          </div>

          {/* Divider with OR */}
          <div className="flex items-center gap-3 py-2">
            <div className="flex-1 h-px bg-slate-700"></div>
            <span className="text-xs text-slate-500 font-medium">{t('onboarding.password.or')}</span>
            <div className="flex-1 h-px bg-slate-700"></div>
          </div>

          {/* Restore from backup button */}
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setShowRestoreDialog(true)}
            className="w-full"
          >
            {t('onboarding.password.restoreFromBackup')}
          </Button>

          {/* Restore dialog */}
          {showRestoreDialog && (
            <div className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
              <motion.div
                initial={{ opacity: 0, scale: 0.95 }}
                animate={{ opacity: 1, scale: 1 }}
                className="bg-slate-800 border border-slate-700 rounded-xl p-6 max-w-md w-full space-y-4"
              >
                <h3 className="text-xl font-semibold text-slate-100">{t('onboarding.backupFile.title')}</h3>

                <div className="space-y-3">
                  <div className="bg-slate-700 rounded-lg p-3">
                    {backupFilePath ? (
                      <p className="text-slate-200 font-mono text-xs break-all">{backupFilePath}</p>
                    ) : (
                      <p className="text-slate-400 text-center py-2 text-sm">
                        {t('onboarding.backupFile.noFileSelected')}
                      </p>
                    )}
                  </div>

                  <Button variant="primary" onClick={handleSelectBackupFile} className="w-full">
                    {t('onboarding.backupFile.browseFiles')}
                  </Button>

                  {backupFilePath && (
                    <>
                      <Input
                        label={t('onboarding.restorePassword.label')}
                        type="password"
                        placeholder={t('onboarding.restorePassword.placeholder')}
                        value={restorePassword}
                        onChange={(e) => setRestorePassword(e.target.value)}
                      />
                      <Button
                        variant="primary"
                        onClick={handleRestoreBackup}
                        disabled={isProcessing}
                        className="w-full"
                      >
                        {isProcessing ? t('onboarding.completion.restoring') : t('action.restore')}
                      </Button>
                    </>
                  )}
                </div>

                <Button variant="ghost" onClick={() => setShowRestoreDialog(false)} className="w-full">
                  {t('action.cancel')}
                </Button>
              </motion.div>
            </div>
          )}
        </motion.div>
      );
    }

    // Step 3: Peer Configuration
    if (step === 3) {
      const progressPercent = searchProgress.total > 0
        ? Math.round((searchProgress.current / searchProgress.total) * 100)
        : 0;

      return (
        <motion.div
          key="peers"
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          exit={{ opacity: 0, x: -20 }}
          className="space-y-2 w-full"
        >
          <div className="text-center flex-shrink-0">
            <div className="text-2xl">🌐</div>
            <h2 className="text-lg sm:text-xl font-semibold text-slate-100">
              {t('onboarding.peers.title')}
            </h2>
            <p className="text-xs text-slate-400">
              {t('onboarding.peers.description')}
            </p>
          </div>

          <div className="space-y-2 flex-shrink-0">
            {/* Protocol Selection & Max RTT */}
            <GlassCard padding="sm">
              <div className="space-y-2">
                <label className="block text-xs font-medium text-slate-200">
                  {t('peers.discovery.protocols')}
                </label>
                <div className="flex flex-wrap gap-1.5">
                  {protocols.map(proto => (
                    <button
                      key={proto.value}
                      onClick={() => toggleProtocol(proto.value)}
                      disabled={isSearching}
                      className={`px-2.5 py-1 text-xs rounded-lg font-medium transition-all ${
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

              {/* Max RTT */}
              <div className="pt-2">
                <label className="block text-xs font-medium text-slate-200 mb-1">
                  {t('peers.discovery.maxRTT')}
                </label>
                <input
                  type="number"
                  value={maxRTT}
                  onChange={(e) => setMaxRTT(parseInt(e.target.value) || 500)}
                  disabled={isSearching}
                  className="w-full px-3 py-1.5 text-sm bg-slate-700 border border-slate-600 rounded-lg text-slate-100 focus:border-emerald-500 focus:outline-none"
                />
              </div>
            </GlassCard>

            {/* Find Peers button & Use Default Peers toggle */}
            <div className="grid grid-cols-2 gap-2">
              <Button
                variant="primary"
                size="sm"
                onClick={handleFindPeers}
                disabled={isSearching || selectedProtocols.length === 0}
                className="w-full"
              >
                {isSearching ? t('peers.discovery.searching') : t('onboarding.peers.findPeers')}
              </Button>

              <GlassCard padding="sm">
                <div className="flex items-center justify-between h-full">
                  <span className="text-xs text-slate-200 mr-2">{t('onboarding.peers.useDefaultPeers')}</span>
                  <label className="relative inline-flex items-center cursor-pointer flex-shrink-0">
                    <input
                      type="checkbox"
                      checked={useDefaultPeers}
                      onChange={(e) => {
                        setUseDefaultPeers(e.target.checked);
                        if (e.target.checked) {
                          setSelectedPeers([]);
                        }
                      }}
                      className="sr-only peer"
                    />
                    <div className="w-9 h-5 bg-slate-600 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-emerald-500 rounded-full peer peer-checked:after:translate-x-full after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-500 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-emerald-500"></div>
                  </label>
                </div>
              </GlassCard>
            </div>

            {/* Cache indicator */}
            {hasCachedPeers && discoveredPeers.length > 0 && (
              <div className="flex items-center justify-center gap-2">
                <span className="text-xs text-slate-500">{t('onboarding.peers.fromCache')}</span>
              </div>
            )}

            {/* Search progress */}
            {isSearching && (
              <motion.div
                initial={{ opacity: 0, y: -10 }}
                animate={{ opacity: 1, y: 0 }}
              >
                <GlassCard padding="sm" variant="strong">
                  <div className="space-y-2">
                    <div className="flex items-center justify-between text-xs">
                      <span className="text-slate-300">
                        {searchProgress.current} / {searchProgress.total}
                      </span>
                      <span className="text-emerald-400 font-semibold">
                        {searchProgress.available_count} {t('peers.discovery.available')}
                      </span>
                    </div>
                    <div className="relative h-1.5 bg-slate-700 rounded-full overflow-hidden">
                      <motion.div
                        className="absolute inset-y-0 left-0 bg-emerald-500"
                        initial={{ width: 0 }}
                        animate={{ width: `${progressPercent}%` }}
                        transition={{ duration: 0.3 }}
                      />
                    </div>
                  </div>
                </GlassCard>
              </motion.div>
            )}

            {/* Instructions or peer list */}
            {!useDefaultPeers && discoveredPeers.length === 0 && !isSearching && (
              <GlassCard padding="sm">
                <p className="text-xs text-slate-400">
                  {t('onboarding.peers.instructions')}
                </p>
              </GlassCard>
            )}
          </div>

          {/* Discovered peers list - fixed max height to keep buttons visible */}
          {!useDefaultPeers && discoveredPeers.length > 0 && (
            <div className="flex-shrink-0">
              <GlassCard padding="sm" className="flex flex-col overflow-hidden">
                <div className="max-h-48 overflow-y-auto space-y-1.5 pr-1">
                  {[...discoveredPeers]
                    .sort((a, b) => a.rtt - b.rtt)
                    .map((peer) => {
                      const isSelected = selectedPeers.some(p => p.address === peer.address);
                      return (
                        <div
                          key={peer.address}
                          onClick={() => togglePeerSelection(peer)}
                          className={`p-2 rounded-lg cursor-pointer transition-all ${
                            isSelected
                              ? 'bg-emerald-500/20 border-2 border-emerald-500'
                              : 'bg-slate-700 border border-slate-600 hover:border-slate-500'
                          }`}
                        >
                          <div className="flex items-center justify-between">
                            <div className="flex-1 min-w-0">
                              <p className="text-slate-200 font-mono text-xs truncate">{peer.address}</p>
                              <p className="text-slate-500 text-[10px]">
                                {peer.protocol} • {peer.rtt}ms
                              </p>
                            </div>
                            {isSelected && (
                              <div className="ml-2 text-emerald-400 text-sm">✓</div>
                            )}
                          </div>
                        </div>
                      );
                    })}
                </div>
              </GlassCard>
            </div>
          )}
        </motion.div>
      );
    }

    return null;
  };

  return (
    <div className="h-screen bg-slate-900 flex flex-col">
      <div className="w-full max-w-6xl mx-auto py-6 px-4 sm:px-6 flex flex-col h-full">
        {/* Progress Bar */}
        {step > 1 && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            className="mb-4 flex-shrink-0"
          >
            <div className="bg-slate-800 border border-slate-700 rounded-lg p-3">
              <div className="flex justify-between items-center mb-1.5">
                <span className="text-xs text-slate-400">
                  {t('onboarding.progress.step')
                    .replace('{{current}}', String(step))
                    .replace('{{total}}', String(totalSteps))}
                </span>
                <span className="text-xs text-slate-400">
                  {Math.round((step / totalSteps) * 100)}%
                </span>
              </div>
              <div className="h-1.5 bg-slate-700 rounded-full overflow-hidden">
                <motion.div
                  initial={{ width: 0 }}
                  animate={{ width: `${(step / totalSteps) * 100}%` }}
                  transition={{ duration: 0.3, ease: 'easeOut' }}
                  className="h-full bg-emerald-500"
                />
              </div>
            </div>
          </motion.div>
        )}

        {/* Main Content */}
        <div className="flex-1 flex flex-col min-h-0">
          <GlassCard padding="lg" variant="strong" className="flex-1 flex flex-col min-h-0">
            <div className="flex-1 min-h-0 overflow-y-auto pr-2">
              <AnimatePresence mode="wait">{renderStepContent()}</AnimatePresence>
            </div>

            {/* Navigation Buttons */}
            {step > 1 && (
              <motion.div
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                className={`flex justify-between gap-3 border-t border-slate-700 flex-shrink-0 ${step === 3 ? 'mt-3 pt-3' : 'mt-6 pt-6'}`}
              >
                <Button variant="ghost" size="sm" onClick={handleBack} disabled={isProcessing}>
                  {t('onboarding.back')}
                </Button>
                <Button variant="primary" size="sm" onClick={handleNext} disabled={isProcessing}>
                  {step === 3 ? t('onboarding.finish') : t('onboarding.next')}
                </Button>
              </motion.div>
            )}
          </GlassCard>
        </div>
      </div>
    </div>
  );
}

export default Onboarding;
