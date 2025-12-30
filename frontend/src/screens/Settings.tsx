import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';
import {
  Button,
  Input,
  Modal,
  GlassCard,
  BentoGrid,
  BentoCard,
  HolographicBorder,
  Badge,
  LoadingSpinner,
} from '../components';
import { toast } from '../components/ui/Toast';
import { useConfig } from '../hooks/useConfig';
import { useI18n } from '../hooks/useI18n';
import {
  SetTheme,
  SetAutoStart,
  OpenURL,
  GetVersion,
  ChangePassword,
  RegenerateKeys,
  GetStorageStats,
  SetMaxMessageSizeMB,
  CreateBackup,
  RestoreBackup,
  ShowSaveFileDialog,
  ShowOpenFileDialog,
} from '../../wailsjs/go/main/App';

// Backup types
type BackupOptionsDTO = {
  backupPath?: string;
  includeDatabase: boolean;
  password: string;
};

type RestoreOptionsDTO = {
  backupPath: string;
  password: string;
};

type SettingsPage = 'hub' | 'general' | 'security' | 'storage' | 'backup' | 'about';

/**
 * Settings Screen - Application configuration
 *
 * Structure:
 * - Hub page with navigation cards
 * - Sub-pages: General, Network, Security, Backup, About
 *
 * Features:
 * - Language and theme switching
 * - Autostart configuration
 * - Peer management
 * - Password change
 * - Key regeneration (dangerous operation)
 * - Backup/restore shortcuts
 * - App info and links
 */
export function Settings() {
  const { t, changeLanguage, language } = useI18n();
  const [currentPage, setCurrentPage] = useState<SettingsPage>('hub');
  const [isProcessing, setIsProcessing] = useState(false);

  const { config, loadConfig } = useConfig({ loadOnMount: true });

  // Modal states
  const [showPasswordModal, setShowPasswordModal] = useState(false);
  const [showKeysModal, setShowKeysModal] = useState(false);

  // Form states
  const [newPassword, setNewPassword] = useState('');
  const [confirmNewPassword, setConfirmNewPassword] = useState('');
  const [currentPassword, setCurrentPassword] = useState('');

  // Version info - loaded from backend
  const [version, setVersion] = useState('loading...');

  // Storage states
  const [storageStats, setStorageStats] = useState<any>(null);
  const [maxMessageSize, setMaxMessageSize] = useState(10);

  // Backup states
  const navigate = useNavigate();
  const [createPassword, setCreatePassword] = useState('');
  const [confirmCreatePassword, setConfirmCreatePassword] = useState('');
  const [includeDatabase, setIncludeDatabase] = useState(true);
  const [isCreating, setIsCreating] = useState(false);
  const [createProgress, setCreateProgress] = useState(0);
  const [restoreFilePath, setRestoreFilePath] = useState('');
  const [restorePassword, setRestorePassword] = useState('');
  const [isRestoring, setIsRestoring] = useState(false);
  const [restoreProgress, setRestoreProgress] = useState(0);
  const [createPasswordError, setCreatePasswordError] = useState('');
  const [createProgressMessage, setCreateProgressMessage] = useState('');
  const [restoreProgressMessage, setRestoreProgressMessage] = useState('');

  // Load version on component mount
  useEffect(() => {
    GetVersion().then(setVersion).catch(() => setVersion('unknown'));
  }, []);

  // Load storage stats when storage page is opened
  useEffect(() => {
    if (currentPage === 'storage') {
      loadStorageStats();
    }
  }, [currentPage]);

  // Listen for backup progress events
  useEffect(() => {
    const backupProgressHandler = (data: { progress: number; message: string }) => {
      setCreateProgress(data.progress);
      setCreateProgressMessage(data.message);
    };

    const restoreProgressHandler = (data: { progress: number; message: string }) => {
      setRestoreProgress(data.progress);
      setRestoreProgressMessage(data.message);
    };

    const configRestoredHandler = async () => {
      await loadConfig();
      setIsRestoring(false);
      setRestoreProgress(0);
      setRestoreProgressMessage('');
      toast.success(t('backup.messages.backupRestoredMessage'));
      setCurrentPage('hub');
    };

    EventsOn('backup:progress', backupProgressHandler);
    EventsOn('restore:progress', restoreProgressHandler);
    EventsOn('config:restored', configRestoredHandler);

    return () => {
      EventsOff('backup:progress');
      EventsOff('restore:progress');
      EventsOff('config:restored');
    };
  }, [loadConfig, t]);

  const loadStorageStats = async () => {
    try {
      const stats = await GetStorageStats();
      setStorageStats(stats);
      setMaxMessageSize(stats.maxMessageSizeMB || 10);
    } catch (error) {
      console.error('Failed to load storage stats:', error);
      toast.error(t('storage.saveFailed'));
    }
  };

  // Handle language change
  const handleLanguageChange = async (lang: 'en' | 'ru') => {
    try {
      setIsProcessing(true);
      // changeLanguage will handle both i18n update and backend config save
      await changeLanguage(lang);
      toast.success(t('settings.messages.languageChanged'));
    } catch (error) {
      toast.error(error instanceof Error ? error.message : t('settings.messages.languageChangeFailed'));
    } finally {
      setIsProcessing(false);
    }
  };

  // Handle theme change
  const handleThemeChange = async (theme: 'light' | 'dark' | 'system') => {
    try {
      setIsProcessing(true);
      await SetTheme(theme);
      await loadConfig();
      toast.success(t('settings.messages.themeChanged'));
    } catch (error) {
      toast.error(error instanceof Error ? error.message : t('settings.messages.themeChangeFailed'));
    } finally {
      setIsProcessing(false);
    }
  };

  // Handle autostart change
  const handleAutoStartChange = async (enabled: boolean) => {
    try {
      setIsProcessing(true);
      await SetAutoStart(enabled);
      await loadConfig();
      toast.success(enabled ? t('settings.messages.autostartEnabled') : t('settings.messages.autostartDisabled'));
    } catch (error) {
      toast.error(error instanceof Error ? error.message : t('settings.messages.autostartChangeFailed'));
    } finally {
      setIsProcessing(false);
    }
  };

  // Handle max message size change
  const handleMaxMessageSizeChange = async () => {
    try {
      setIsProcessing(true);
      toast.loading(t('storage.savingSettings'));
      await SetMaxMessageSizeMB(maxMessageSize);
      toast.dismiss();
      toast.success(t('storage.settingsSaved'));
      await loadStorageStats();
    } catch (error) {
      toast.dismiss();
      toast.error(error instanceof Error ? error.message : t('storage.saveFailed'));
    } finally {
      setIsProcessing(false);
    }
  };

  // Handle password change
  const handleChangePassword = async () => {
    if (!currentPassword) {
      toast.error(t('settings.messages.passwordRequired'));
      return;
    }
    if (newPassword.length < 8) {
      toast.error(t('settings.messages.passwordShort'));
      return;
    }
    if (newPassword !== confirmNewPassword) {
      toast.error(t('settings.messages.passwordMismatch'));
      return;
    }
    setIsProcessing(true);
    try {
      await ChangePassword(currentPassword, newPassword);
      toast.success(t('settings.messages.passwordChanged'));
      setShowPasswordModal(false);
      setNewPassword('');
      setConfirmNewPassword('');
      setCurrentPassword('');
    } catch (error) {
      toast.error(error instanceof Error ? error.message : t('settings.messages.passwordChangeFailed'));
    } finally {
      setIsProcessing(false);
    }
  };

  // Handle key regeneration
  const handleRegenerateKeys = async () => {
    if (!currentPassword) {
      toast.error(t('settings.messages.passwordRequired'));
      return;
    }
    setIsProcessing(true);
    try {
      await RegenerateKeys(currentPassword);
      toast.success(t('settings.messages.keysRegeneratedMessage'));
      setShowKeysModal(false);
      setCurrentPassword('');
    } catch (error) {
      toast.error(error instanceof Error ? error.message : t('settings.messages.keyRegenerationFailed'));
    } finally {
      setIsProcessing(false);
    }
  };

  // Handle create backup
  const handleCreateBackup = async () => {
    if (createPassword.length < 8) {
      setCreatePasswordError(t('backup.messages.passwordShort'));
      return;
    }
    if (createPassword !== confirmCreatePassword) {
      setCreatePasswordError(t('backup.messages.passwordMismatch'));
      return;
    }
    setCreatePasswordError('');

    try {
      const now = new Date();
      const day = String(now.getDate()).padStart(2, '0');
      const month = String(now.getMonth() + 1).padStart(2, '0');
      const year = String(now.getFullYear()).slice(-2);
      const defaultFilename = `tbackup-${day}-${month}-${year}.tb`;

      const savePath = await ShowSaveFileDialog(t('backup.messages.saveBackup'), defaultFilename);
      if (!savePath) return;

      setIsCreating(true);
      setCreateProgress(0);
      setCreateProgressMessage('');

      try {
        const options: BackupOptionsDTO = {
          backupPath: savePath,
          includeDatabase,
          password: createPassword,
        };
        await CreateBackup(options);

        setTimeout(() => {
          toast.success(t('backup.messages.backupCreatedMessage', { path: savePath }));
          setCreatePassword('');
          setConfirmCreatePassword('');
          setIsCreating(false);
          setCreateProgress(0);
          setCreateProgressMessage('');
        }, 500);
      } catch (error) {
        throw error;
      }
    } catch (error) {
      setIsCreating(false);
      setCreateProgress(0);
      toast.error(error instanceof Error ? error.message : t('backup.messages.backupFailedMessage'));
    }
  };

  // Handle select backup file
  const handleSelectBackupFile = async () => {
    try {
      const filePath = await ShowOpenFileDialog(t('backup.messages.selectBackupFileTitle'));
      if (filePath) {
        setRestoreFilePath(filePath);
      }
    } catch (error) {
      toast.error(error instanceof Error ? error.message : t('backup.messages.fileSelectionFailedMessage'));
    }
  };

  // Handle restore backup
  const handleRestoreBackup = async () => {
    if (!restoreFilePath) {
      toast.error(t('backup.messages.selectBackupFile'));
      return;
    }
    if (!restorePassword) {
      toast.error(t('backup.messages.enterPasswordMessage'));
      return;
    }

    setIsRestoring(true);
    setRestoreProgress(0);
    setRestoreProgressMessage('');

    try {
      const options: RestoreOptionsDTO = {
        backupPath: restoreFilePath,
        password: restorePassword,
      };
      await RestoreBackup(options);

      setRestoreFilePath('');
      setRestorePassword('');
    } catch (error) {
      setIsRestoring(false);
      setRestoreProgress(0);
      setRestoreProgressMessage('');
      toast.error(error instanceof Error ? error.message : t('backup.messages.restoreFailedMessage'));
    }
  };

  // Render hub page with navigation cards
  const renderHub = () => (
    <motion.div
      key="hub"
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -20 }}
      className="space-y-6"
    >
      <div className="text-center space-y-2 mb-8">
        <h1 className="text-5xl font-display font-black text-transparent bg-clip-text bg-iridescent bg-[length:200%_100%] animate-holographic-spin">
          {t('settings.hub.title')}
        </h1>
        <p className="text-lg text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant font-futuristic">
          {t('settings.hub.subtitle')}
        </p>
      </div>

      <BentoGrid columns={3} gap="lg">
        <BentoCard span={1}>
          <HolographicBorder borderWidth={1}>
            <button
              onClick={() => setCurrentPage('general')}
              className="w-full p-8 bg-md-light-surfaceVariant dark:bg-md-dark-surfaceVariant hover:bg-md-light-primaryContainer/50 dark:hover:bg-md-dark-primaryContainer/30 rounded-xl transition-all text-left space-y-3"
            >
              <div className="text-5xl">⚙️</div>
              <h3 className="text-2xl font-bold text-md-light-onSurface dark:text-md-dark-onSurface">{t('settings.general')}</h3>
              <p className="text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant text-sm">
                {t('settings.hub.generalDescription')}
              </p>
            </button>
          </HolographicBorder>
        </BentoCard>

        <BentoCard span={1}>
          <HolographicBorder borderWidth={1}>
            <button
              onClick={() => setCurrentPage('security')}
              className="w-full p-8 bg-md-light-surfaceVariant dark:bg-md-dark-surfaceVariant hover:bg-md-light-primaryContainer/50 dark:hover:bg-md-dark-primaryContainer/30 rounded-xl transition-all text-left space-y-3"
            >
              <div className="text-5xl">🔐</div>
              <h3 className="text-2xl font-bold text-md-light-onSurface dark:text-md-dark-onSurface">{t('settings.security')}</h3>
              <p className="text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant text-sm">
                {t('settings.hub.securityDescription')}
              </p>
            </button>
          </HolographicBorder>
        </BentoCard>

        <BentoCard span={1}>
          <HolographicBorder borderWidth={1}>
            <button
              onClick={() => setCurrentPage('storage')}
              className="w-full p-8 bg-md-light-surfaceVariant dark:bg-md-dark-surfaceVariant hover:bg-md-light-primaryContainer/50 dark:hover:bg-md-dark-primaryContainer/30 rounded-xl transition-all text-left space-y-3"
            >
              <div className="text-5xl">💾</div>
              <h3 className="text-2xl font-bold text-md-light-onSurface dark:text-md-dark-onSurface">{t('settings.storage')}</h3>
              <p className="text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant text-sm">
                {t('settings.hub.storageDescription')}
              </p>
            </button>
          </HolographicBorder>
        </BentoCard>

        <BentoCard span={1}>
          <HolographicBorder borderWidth={1}>
            <button
              onClick={() => setCurrentPage('backup')}
              className="w-full p-8 bg-md-light-surfaceVariant dark:bg-md-dark-surfaceVariant hover:bg-md-light-primaryContainer/50 dark:hover:bg-md-dark-primaryContainer/30 rounded-xl transition-all text-left space-y-3"
            >
              <div className="text-5xl">📦</div>
              <h3 className="text-2xl font-bold text-md-light-onSurface dark:text-md-dark-onSurface">{t('settings.backup')}</h3>
              <p className="text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant text-sm">
                {t('settings.hub.backupDescription')}
              </p>
            </button>
          </HolographicBorder>
        </BentoCard>

        <BentoCard span={1}>
          <HolographicBorder borderWidth={1}>
            <button
              onClick={() => setCurrentPage('about')}
              className="w-full p-8 bg-md-light-surfaceVariant dark:bg-md-dark-surfaceVariant hover:bg-md-light-primaryContainer/50 dark:hover:bg-md-dark-primaryContainer/30 rounded-xl transition-all text-left space-y-3"
            >
              <div className="text-5xl">ℹ️</div>
              <h3 className="text-2xl font-bold text-md-light-onSurface dark:text-md-dark-onSurface">{t('settings.about')}</h3>
              <p className="text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant text-sm">
                {t('settings.hub.aboutDescription')}
              </p>
            </button>
          </HolographicBorder>
        </BentoCard>
      </BentoGrid>
    </motion.div>
  );

  // Render General page
  const renderGeneral = () => (
    <motion.div
      key="general"
      initial={{ opacity: 0, x: 50 }}
      animate={{ opacity: 1, x: 0 }}
      exit={{ opacity: 0, x: -50 }}
      className="space-y-6"
    >
      <GlassCard title={t('settings.generalSettings.title')} padding="lg" className="min-h-[500px]">
        <div className="space-y-6">
          {/* Language Selector */}
          <div>
            <label className="block text-sm font-medium text-md-light-onSurface dark:text-md-dark-onSurface mb-2">
              {t('label.language')}
            </label>
            <div className="flex gap-3">
              <Button
                variant={language === 'en' ? 'primary' : 'ghost'}
                glow={language === 'en'}
                onClick={() => handleLanguageChange('en')}
                disabled={isProcessing}
                className="flex-1"
              >
                🇬🇧 English
              </Button>
              <Button
                variant={language === 'ru' ? 'primary' : 'ghost'}
                glow={language === 'ru'}
                onClick={() => handleLanguageChange('ru')}
                disabled={isProcessing}
                className="flex-1"
              >
                🇷🇺 Русский
              </Button>
            </div>
          </div>

          {/* Theme Selector */}
          <div>
            <label className="block text-sm font-medium text-md-light-onSurface dark:text-md-dark-onSurface mb-2">
              {t('label.theme')}
            </label>
            <div className="flex gap-3">
              <Button
                variant={config?.theme === 'light' ? 'primary' : 'ghost'}
                glow={config?.theme === 'light'}
                onClick={() => handleThemeChange('light')}
                disabled={isProcessing}
                className="flex-1"
              >
                ☀️ {t('settings.theme.light')}
              </Button>
              <Button
                variant={config?.theme === 'dark' ? 'primary' : 'ghost'}
                glow={config?.theme === 'dark'}
                onClick={() => handleThemeChange('dark')}
                disabled={isProcessing}
                className="flex-1"
              >
                🌙 {t('settings.theme.dark')}
              </Button>
              <Button
                variant={config?.theme === 'system' ? 'primary' : 'ghost'}
                glow={config?.theme === 'system'}
                onClick={() => handleThemeChange('system')}
                disabled={isProcessing}
                className="flex-1"
              >
                💻 {t('settings.theme.system')}
              </Button>
            </div>
          </div>

          {/* Autostart Toggle */}
          <div className="flex items-center justify-between p-4 bg-md-light-surfaceVariant dark:bg-md-dark-surfaceVariant rounded-lg border border-md-light-outline/30 dark:border-md-dark-outline/30">
            <div>
              <p className="text-md-light-onSurface dark:text-md-dark-onSurface font-medium">{t('settings.autostart')}</p>
              <p className="text-sm text-md-light-outline dark:text-md-dark-outline mt-1">
                {t('settings.generalSettings.autostartDescription')}
              </p>
            </div>
            <button
              onClick={() => handleAutoStartChange(!config?.autoStart)}
              disabled={isProcessing}
              className={`relative w-14 h-8 rounded-full transition-colors ${
                config?.autoStart ? 'bg-md-light-primary dark:bg-md-dark-primary' : 'bg-md-light-outline/20 dark:bg-md-dark-outline/20'
              } ${isProcessing ? 'opacity-50 cursor-not-allowed' : ''}`}
            >
              <motion.div
                animate={{ x: config?.autoStart ? 24 : 2 }}
                transition={{ type: 'spring', stiffness: 500, damping: 30 }}
                className="absolute top-1 w-6 h-6 bg-white rounded-full shadow-lg"
              />
            </button>
          </div>
        </div>
      </GlassCard>
    </motion.div>
  );


  // Render Security page
  const renderSecurity = () => (
    <motion.div
      key="security"
      initial={{ opacity: 0, x: 50 }}
      animate={{ opacity: 1, x: 0 }}
      exit={{ opacity: 0, x: -50 }}
      className="space-y-6"
    >
      <GlassCard title={t('settings.securitySettings.title')} padding="lg" className="min-h-[500px]">
        <div className="space-y-6">
          {/* Change Password */}
          <div className="p-6 bg-md-light-surfaceVariant dark:bg-md-dark-surfaceVariant rounded-lg border border-md-light-outline/30 dark:border-md-dark-outline/30 space-y-3">
            <div className="flex items-center gap-3">
              <div className="text-3xl">🔐</div>
              <div className="flex-1">
                <h3 className="text-md-light-onSurface dark:text-md-dark-onSurface font-semibold">{t('settings.securitySettings.changePassword')}</h3>
                <p className="text-sm text-md-light-outline dark:text-md-dark-outline mt-1">
                  {t('settings.securitySettings.changePasswordDescription')}
                </p>
              </div>
            </div>
            <Button
              variant="secondary"
              glow
              onClick={() => setShowPasswordModal(true)}
              className="w-full"
            >
              {t('settings.securitySettings.changePassword')}
            </Button>
          </div>

          {/* Regenerate Keys */}
          <div className="p-6 bg-md-light-errorContainer/50 dark:bg-md-dark-errorContainer/30 rounded-lg border border-md-light-error/30 dark:border-md-dark-error/30 space-y-3">
            <div className="flex items-center gap-3">
              <div className="text-3xl">⚠️</div>
              <div className="flex-1">
                <h3 className="text-md-light-onSurface dark:text-md-dark-onSurface font-semibold">{t('settings.securitySettings.regenerateKeys')}</h3>
                <p className="text-sm text-md-light-onErrorContainer dark:text-md-dark-onErrorContainer mt-1">
                  {t('settings.securitySettings.regenerateKeysDescription')}
                </p>
              </div>
            </div>
            <Button
              variant="danger"
              onClick={() => setShowKeysModal(true)}
              className="w-full"
            >
              {t('settings.securitySettings.regenerateKeys')}
            </Button>
          </div>
        </div>
      </GlassCard>
    </motion.div>
  );

  // Render Storage page
  const renderStorage = () => (
    <motion.div
      key="storage"
      initial={{ opacity: 0, x: 50 }}
      animate={{ opacity: 1, x: 0 }}
      exit={{ opacity: 0, x: 50 }}
      className="space-y-6"
    >
      <div className="mb-6">
        <h2 className="text-3xl font-display font-black text-transparent bg-clip-text bg-iridescent bg-[length:200%_100%] animate-holographic-spin">
          {t('storage.title')}
        </h2>
        <p className="text-sm text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant font-futuristic">
          {t('storage.subtitle')}
        </p>
      </div>

      <HolographicBorder borderWidth={2}>
        <GlassCard padding="lg">
          {/* Limits Section */}
          <div className="space-y-4">
            <h3 className="text-xl font-bold text-md-light-onSurface dark:text-md-dark-onSurface">
              {t('storage.limits')}
            </h3>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-md-light-onSurface dark:text-md-dark-onSurface mb-2">
                  {t('storage.maxMessageSize')}
                </label>
                <p className="text-xs text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant mb-4">
                  {t('storage.maxMessageSizeDescription')}
                </p>
                <input
                  type="range"
                  min={10}
                  max={500}
                  step={10}
                  value={maxMessageSize}
                  onChange={(e) => setMaxMessageSize(Number(e.target.value))}
                  className="w-full h-2 bg-md-light-surfaceVariant dark:bg-md-dark-surfaceVariant rounded-lg appearance-none cursor-pointer accent-md-light-primary dark:accent-md-dark-primary"
                />
                <div className="flex justify-between text-xs text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant mt-2">
                  <span>10 {t('storage.mb')}</span>
                  <span className="text-lg font-bold text-md-light-primary dark:text-md-dark-primary">
                    {maxMessageSize} {t('storage.mb')}
                  </span>
                  <span>500 {t('storage.mb')}</span>
                </div>
              </div>
              <Button
                variant="primary"
                glow
                onClick={handleMaxMessageSizeChange}
                disabled={isProcessing || !storageStats}
                className="w-full"
              >
                {t('action.save')}
              </Button>
            </div>
          </div>
        </GlassCard>
      </HolographicBorder>
    </motion.div>
  );

  // Render Backup page
  const renderBackup = () => (
    <motion.div
      key="backup"
      initial={{ opacity: 0, x: 50 }}
      animate={{ opacity: 1, x: 0 }}
      exit={{ opacity: 0, x: -50 }}
      className="space-y-6 pb-6"
    >
      <div className="mb-6">
        <h2 className="text-3xl font-display font-black text-transparent bg-clip-text bg-iridescent bg-[length:200%_100%] animate-holographic-spin">
          {t('backup.title')}
        </h2>
        <p className="text-sm text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant font-futuristic">
          {t('backup.subtitle')}
        </p>
      </div>

      {/* Main Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Create Backup */}
        <HolographicBorder animated={isCreating} borderWidth={isCreating ? 3 : 2}>
          <GlassCard
            title={t('backup.createBackup')}
            subtitle={t('backup.createBackupSubtitle')}
            padding="lg"
          >
            {isCreating ? (
              <div className="space-y-6 py-8">
                <LoadingSpinner size="xl" variant="holographic" text={createProgressMessage || t('backup.creating')} />
                <div className="space-y-2">
                  <div className="flex justify-between text-sm text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant">
                    <span>{createProgressMessage || t('backup.progress')}</span>
                    <span>{createProgress}%</span>
                  </div>
                  <div className="h-2 bg-md-light-surfaceVariant dark:bg-md-dark-surfaceVariant rounded-full overflow-hidden">
                    <motion.div
                      animate={{ width: `${createProgress}%` }}
                      transition={{ duration: 0.3 }}
                      className="h-full bg-gradient-to-r from-md-light-primary via-md-light-secondary to-md-light-tertiary dark:from-md-dark-primary dark:via-md-dark-secondary dark:to-md-dark-tertiary"
                    />
                  </div>
                </div>
              </div>
            ) : (
              <div className="space-y-6">
                <div className="text-center mb-6">
                  <div className="text-7xl mb-4">💾</div>
                </div>

                <div className="space-y-4">
                  <Input
                    label={t('backup.backupPassword')}
                    type="password"
                    placeholder={t('backup.enterPassword')}
                    value={createPassword}
                    onChange={(e) => setCreatePassword(e.target.value)}
                    error={createPasswordError}
                  />
                  <Input
                    label={t('backup.confirmPassword')}
                    type="password"
                    placeholder={t('backup.reenterPassword')}
                    value={confirmCreatePassword}
                    onChange={(e) => setConfirmCreatePassword(e.target.value)}
                  />

                  {/* Include Database Checkbox */}
                  <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 p-4 bg-md-light-surfaceVariant dark:bg-md-dark-surfaceVariant rounded-lg border border-md-light-outline/30 dark:border-md-dark-outline/30">
                    <div className="flex-1">
                      <p className="text-md-light-onSurface dark:text-md-dark-onSurface font-medium">{t('backup.includeDatabaseLabel')}</p>
                      <p className="text-sm text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant mt-1">
                        {t('backup.includeDatabaseDescription')}
                      </p>
                    </div>
                    <button
                      onClick={() => setIncludeDatabase(!includeDatabase)}
                      className={`relative w-14 h-8 rounded-full transition-colors flex-shrink-0 ${
                        includeDatabase ? 'bg-md-light-primary dark:bg-md-dark-primary' : 'bg-md-light-outline/20 dark:bg-md-dark-outline/20'
                      }`}
                    >
                      <motion.div
                        animate={{ x: includeDatabase ? 24 : 2 }}
                        transition={{ type: 'spring', stiffness: 500, damping: 30 }}
                        className="absolute top-1 w-6 h-6 bg-white rounded-full shadow-lg"
                      />
                    </button>
                  </div>

                  <Button
                    variant="primary"
                    size="lg"
                    glow
                    onClick={handleCreateBackup}
                    className="w-full"
                  >
                    {t('backup.createButton')}
                  </Button>
                </div>

                <div className="bg-md-light-primaryContainer/30 dark:bg-md-dark-primaryContainer/30 border border-md-light-primary/30 dark:border-md-dark-primary/30 rounded-lg p-4 mt-6">
                  <p className="text-sm text-md-light-onSurface dark:text-md-dark-onSurface font-medium mb-2">
                    {t('backup.whatsIncluded')}
                  </p>
                  <ul className="list-disc list-inside text-sm text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant space-y-1">
                    <li>{t('backup.includedItem1')}</li>
                    <li>{t('backup.includedItem2')}</li>
                    <li>{t('backup.includedItem3')}</li>
                    {includeDatabase && <li>{t('backup.includedItem4')}</li>}
                  </ul>
                  <p className="text-xs text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant mt-3">
                    {t('backup.encryptionInfo')}
                  </p>
                </div>
              </div>
            )}
          </GlassCard>
        </HolographicBorder>

        {/* Restore Backup */}
        <HolographicBorder animated={isRestoring} borderWidth={isRestoring ? 3 : 2}>
          <GlassCard
            title={t('backup.restoreBackup')}
            subtitle={t('backup.restoreBackupSubtitle')}
            padding="lg"
          >
            {isRestoring ? (
              <div className="space-y-6 py-8">
                <LoadingSpinner size="xl" variant="holographic" text={restoreProgressMessage || t('backup.restoring')} />
                <div className="space-y-2">
                  <div className="flex justify-between text-sm text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant">
                    <span>{restoreProgressMessage || t('backup.progress')}</span>
                    <span>{restoreProgress}%</span>
                  </div>
                  <div className="h-2 bg-md-light-surfaceVariant dark:bg-md-dark-surfaceVariant rounded-full overflow-hidden">
                    <motion.div
                      animate={{ width: `${restoreProgress}%` }}
                      transition={{ duration: 0.3 }}
                      className="h-full bg-gradient-to-r from-md-light-tertiary via-md-light-secondary to-md-light-primary dark:from-md-dark-tertiary dark:via-md-dark-secondary dark:to-md-dark-primary"
                    />
                  </div>
                </div>
              </div>
            ) : (
              <div className="space-y-6">
                <div className="text-center mb-6">
                  <div className="text-7xl mb-4">📂</div>
                </div>

                <div className="space-y-4">
                  {/* File Selection */}
                  <div>
                    <label className="block text-sm font-medium text-md-light-onSurface dark:text-md-dark-onSurface mb-2">
                      {t('backup.backupFile')}
                    </label>
                    <div className="bg-md-light-surfaceVariant dark:bg-md-dark-surfaceVariant border border-md-light-outline/30 dark:border-md-dark-outline/30 rounded-lg p-4">
                      {restoreFilePath ? (
                        <div className="space-y-2">
                          <p className="text-xs text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant uppercase tracking-wide">{t('backup.selectedFile')}</p>
                          <p className="text-md-light-onSurface dark:text-md-dark-onSurface font-mono text-sm break-all">{restoreFilePath}</p>
                        </div>
                      ) : (
                        <p className="text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant text-center py-2">{t('backup.noFileSelected')}</p>
                      )}
                    </div>
                    <Button
                      variant="secondary"
                      glow
                      onClick={handleSelectBackupFile}
                      className="w-full mt-2"
                    >
                      {t('backup.browseFiles')}
                    </Button>
                  </div>

                  {/* Password Input */}
                  <Input
                    label={t('backup.backupPassword')}
                    type="password"
                    placeholder={t('backup.enterBackupPassword')}
                    value={restorePassword}
                    onChange={(e) => setRestorePassword(e.target.value)}
                  />

                  <Button
                    variant="primary"
                    size="lg"
                    glow
                    onClick={handleRestoreBackup}
                    disabled={!restoreFilePath || !restorePassword}
                    className="w-full"
                  >
                    {t('backup.restoreButton')}
                  </Button>
                </div>

                <div className="bg-md-light-errorContainer/20 dark:bg-md-dark-errorContainer/20 border border-md-light-error/30 dark:border-md-dark-error/30 rounded-lg p-4 mt-6">
                  <p className="text-sm text-md-light-onSurface dark:text-md-dark-onSurface font-medium mb-2">
                    {t('backup.warningTitle')}
                  </p>
                  <ul className="list-disc list-inside text-sm text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant space-y-1">
                    <li>{t('backup.warningItem1')}</li>
                    <li>{t('backup.warningItem2')}</li>
                    <li>{t('backup.warningItem3')}</li>
                  </ul>
                </div>
              </div>
            )}
          </GlassCard>
        </HolographicBorder>
      </div>

      {/* Info Section */}
      <HolographicBorder borderWidth={1}>
        <div className="bg-md-light-primaryContainer/20 dark:bg-md-dark-primaryContainer/20 border border-md-light-primary/30 dark:border-md-dark-primary/30 rounded-lg p-6">
          <div className="flex items-start gap-4">
            <div className="text-4xl">💡</div>
            <div className="flex-1 space-y-2">
              <h3 className="text-md-light-onSurface dark:text-md-dark-onSurface font-semibold text-lg">{t('backup.bestPractices')}</h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm text-md-light-onSurface dark:text-md-dark-onSurface">
                <div>
                  <p className="font-medium mb-2">{t('backup.doTitle')}</p>
                  <ul className="list-disc list-inside space-y-1 text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant">
                    <li>{t('backup.doItem1')}</li>
                    <li>{t('backup.doItem2')}</li>
                    <li>{t('backup.doItem3')}</li>
                    <li>{t('backup.doItem4')}</li>
                  </ul>
                </div>
                <div>
                  <p className="font-medium mb-2">{t('backup.dontTitle')}</p>
                  <ul className="list-disc list-inside space-y-1 text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant">
                    <li>{t('backup.dontItem1')}</li>
                    <li>{t('backup.dontItem2')}</li>
                    <li>{t('backup.dontItem3')}</li>
                    <li>{t('backup.dontItem4')}</li>
                  </ul>
                </div>
              </div>
              <div className="pt-3 mt-3 border-t border-md-light-outline/10 dark:border-md-dark-outline/10">
                <p className="text-xs text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant">
                  {t('backup.encryptionNotice')}
                </p>
              </div>
            </div>
          </div>
        </div>
      </HolographicBorder>
    </motion.div>
  );

  // Render About page
  const renderAbout = () => (
    <motion.div
      key="about"
      initial={{ opacity: 0, x: 50 }}
      animate={{ opacity: 1, x: 0 }}
      exit={{ opacity: 0, x: -50 }}
      className="space-y-6"
    >
      <GlassCard title={t('settings.aboutSettings.title')} padding="lg">
        <div className="space-y-6 text-center">
          <div className="flex justify-center mb-4">
            <img src="/appicon.png" alt="Tyr Desktop" className="w-24 h-24 sm:w-32 sm:h-32" />
          </div>
          <h2 className="text-2xl sm:text-4xl font-display font-black text-transparent bg-clip-text bg-iridescent bg-[length:200%_100%]">
            Tyr Desktop
          </h2>
          <Badge variant="info" size="lg">
            Version {version || '2.0.0'}
          </Badge>
          <p className="text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant text-sm sm:text-base px-4">
            {t('settings.aboutSettings.description')}
          </p>
          <div className="pt-6 space-y-3 flex flex-col items-stretch px-4">
            <Button
              variant="primary"
              glow
              onClick={() => OpenURL('https://github.com/JB-SelfCompany/Tyr-Desktop')}
              className="w-full"
            >
              {t('settings.aboutSettings.githubRepository')}
            </Button>
            <Button
              variant="ghost"
              onClick={() => OpenURL('https://yggdrasil-network.github.io')}
              className="w-full"
            >
              {t('settings.aboutSettings.yggdrasilNetwork')}
            </Button>
          </div>
          <div className="pt-6 border-t border-md-light-outline/30 dark:border-md-dark-outline/30">
            <p className="text-xs sm:text-sm text-md-light-outline dark:text-md-dark-outline px-4">
              {t('settings.aboutSettings.madeWith')}
            </p>
          </div>
        </div>
      </GlassCard>
    </motion.div>
  );

  return (
    <div className="space-y-6">
      {/* Back Button (only show on sub-pages) */}
      {currentPage !== 'hub' && (
        <Button variant="ghost" onClick={() => setCurrentPage('hub')}>
          {t('settings.backToSettings')}
        </Button>
      )}

      {/* Page Content */}
      <AnimatePresence mode="wait">
        {currentPage === 'hub' && renderHub()}
        {currentPage === 'general' && renderGeneral()}
        {currentPage === 'security' && renderSecurity()}
        {currentPage === 'storage' && renderStorage()}
        {currentPage === 'backup' && renderBackup()}
        {currentPage === 'about' && renderAbout()}
      </AnimatePresence>

      {/* Change Password Modal */}
      <Modal
        isOpen={showPasswordModal}
        onClose={() => setShowPasswordModal(false)}
        title={t('settings.securitySettings.changePassword')}
        size="md"
      >
        <div className="space-y-4">
          <Input
            label={t('settings.securitySettings.currentPassword')}
            type="password"
            placeholder={t('settings.securitySettings.enterCurrentPassword')}
            value={currentPassword}
            onChange={(e) => setCurrentPassword(e.target.value)}
          />
          <Input
            label={t('settings.securitySettings.newPassword')}
            type="password"
            placeholder={t('settings.securitySettings.enterNewPassword')}
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
          />
          <Input
            label={t('settings.securitySettings.confirmNewPassword')}
            type="password"
            placeholder={t('settings.securitySettings.reenterNewPassword')}
            value={confirmNewPassword}
            onChange={(e) => setConfirmNewPassword(e.target.value)}
          />
          <div className="flex gap-3 justify-end pt-4">
            <Button variant="ghost" onClick={() => setShowPasswordModal(false)}>
              {t('action.cancel')}
            </Button>
            <Button
              variant="primary"
              glow
              onClick={handleChangePassword}
              disabled={isProcessing}
            >
              {isProcessing ? t('settings.securitySettings.changing') : t('settings.securitySettings.changePassword')}
            </Button>
          </div>
        </div>
      </Modal>

      {/* Regenerate Keys Modal */}
      <Modal
        isOpen={showKeysModal}
        onClose={() => setShowKeysModal(false)}
        title={`⚠️ ${t('settings.securitySettings.regenerateKeys')}`}
        size="md"
      >
        <div className="space-y-4">
          <div className="bg-md-light-errorContainer/50 dark:bg-md-dark-errorContainer/30 border border-md-light-error/30 dark:border-md-dark-error/30 rounded-lg p-4">
            <p className="text-md-light-onSurface dark:text-md-dark-onSurface font-semibold mb-2">{t('settings.securitySettings.warningTitle')}</p>
            <p className="text-sm text-md-light-onSurface dark:text-md-dark-onSurface">
              {t('settings.securitySettings.warningDescription')}
            </p>
            <ul className="list-disc list-inside text-sm text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant mt-2 space-y-1">
              <li>{t('settings.securitySettings.warningItem1')}</li>
              <li>{t('settings.securitySettings.warningItem2')}</li>
              <li>{t('settings.securitySettings.warningItem3')}</li>
              <li>{t('settings.securitySettings.warningItem4')}</li>
            </ul>
          </div>
          <p className="text-md-light-onSurface dark:text-md-dark-onSurface">
            {t('settings.securitySettings.confirmWarning')}
          </p>
          <Input
            label={t('settings.securitySettings.currentPassword')}
            type="password"
            placeholder={t('settings.securitySettings.enterCurrentPassword')}
            value={currentPassword}
            onChange={(e) => setCurrentPassword(e.target.value)}
          />
          <div className="flex gap-3 justify-end pt-4">
            <Button variant="ghost" onClick={() => setShowKeysModal(false)}>
              {t('action.cancel')}
            </Button>
            <Button
              variant="danger"
              onClick={handleRegenerateKeys}
              disabled={isProcessing || !currentPassword}
            >
              {isProcessing ? t('settings.securitySettings.regenerating') : t('settings.securitySettings.yesRegenerateKeys')}
            </Button>
          </div>
        </div>
      </Modal>

    </div>
  );
}

export default Settings;
