import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';
import {
  Button,
  Input,
  Modal,
  GlassCard,
  Badge,
  LoadingSpinner,
} from '../components';
import { toast } from '../components/ui/Toast';
import { useConfig } from '../hooks/useConfig';
import { useI18n } from '../hooks/useI18n';
import {
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

export function Settings() {
  const { t, changeLanguage, language } = useI18n();
  const [currentPage, setCurrentPage] = useState<SettingsPage>('hub');
  const [isProcessing, setIsProcessing] = useState(false);

  const { config, loadConfig } = useConfig({ loadOnMount: true });

  const [showPasswordModal, setShowPasswordModal] = useState(false);
  const [showKeysModal, setShowKeysModal] = useState(false);

  const [newPassword, setNewPassword] = useState('');
  const [confirmNewPassword, setConfirmNewPassword] = useState('');
  const [currentPassword, setCurrentPassword] = useState('');

  const [version, setVersion] = useState('loading...');
  const [storageStats, setStorageStats] = useState<any>(null);
  const [maxMessageSize, setMaxMessageSize] = useState(10);

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

  useEffect(() => {
    GetVersion().then(setVersion).catch(() => setVersion('unknown'));
  }, []);

  useEffect(() => {
    if (currentPage === 'storage') {
      loadStorageStats();
    }
  }, [currentPage]);

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

  const handleLanguageChange = async (lang: 'en' | 'ru') => {
    try {
      setIsProcessing(true);
      await changeLanguage(lang);
      toast.success(t('settings.messages.languageChanged'));
    } catch (error) {
      toast.error(error instanceof Error ? error.message : t('settings.messages.languageChangeFailed'));
    } finally {
      setIsProcessing(false);
    }
  };

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
      setIsCreating(false);
      setCreateProgress(0);
      toast.error(error instanceof Error ? error.message : t('backup.messages.backupFailedMessage'));
    }
  };

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

  // Hub page
  const renderHub = () => (
    <motion.div
      key="hub"
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0 }}
      className="space-y-6"
    >
      <div className="mb-6">
        <h1 className="text-2xl font-semibold text-slate-100">{t('settings.hub.title')}</h1>
        <p className="text-sm text-slate-400 mt-1">{t('settings.hub.subtitle')}</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {/* General */}
        <button
          onClick={() => setCurrentPage('general')}
          className="p-6 bg-slate-800 hover:bg-slate-700 border border-slate-700 hover:border-slate-600 rounded-xl transition-all text-left"
        >
          <div className="text-4xl mb-3">⚙️</div>
          <h3 className="text-lg font-semibold text-slate-100">{t('settings.general')}</h3>
          <p className="text-sm text-slate-400 mt-1">{t('settings.hub.generalDescription')}</p>
        </button>

        {/* Security */}
        <button
          onClick={() => setCurrentPage('security')}
          className="p-6 bg-slate-800 hover:bg-slate-700 border border-slate-700 hover:border-slate-600 rounded-xl transition-all text-left"
        >
          <div className="text-4xl mb-3">🔐</div>
          <h3 className="text-lg font-semibold text-slate-100">{t('settings.security')}</h3>
          <p className="text-sm text-slate-400 mt-1">{t('settings.hub.securityDescription')}</p>
        </button>

        {/* Storage */}
        <button
          onClick={() => setCurrentPage('storage')}
          className="p-6 bg-slate-800 hover:bg-slate-700 border border-slate-700 hover:border-slate-600 rounded-xl transition-all text-left"
        >
          <div className="text-4xl mb-3">💾</div>
          <h3 className="text-lg font-semibold text-slate-100">{t('settings.storage')}</h3>
          <p className="text-sm text-slate-400 mt-1">{t('settings.hub.storageDescription')}</p>
        </button>

        {/* Backup */}
        <button
          onClick={() => setCurrentPage('backup')}
          className="p-6 bg-slate-800 hover:bg-slate-700 border border-slate-700 hover:border-slate-600 rounded-xl transition-all text-left"
        >
          <div className="text-4xl mb-3">📦</div>
          <h3 className="text-lg font-semibold text-slate-100">{t('settings.backup')}</h3>
          <p className="text-sm text-slate-400 mt-1">{t('settings.hub.backupDescription')}</p>
        </button>

        {/* About */}
        <button
          onClick={() => setCurrentPage('about')}
          className="p-6 bg-slate-800 hover:bg-slate-700 border border-slate-700 hover:border-slate-600 rounded-xl transition-all text-left"
        >
          <div className="text-4xl mb-3">ℹ️</div>
          <h3 className="text-lg font-semibold text-slate-100">{t('settings.about')}</h3>
          <p className="text-sm text-slate-400 mt-1">{t('settings.hub.aboutDescription')}</p>
        </button>
      </div>
    </motion.div>
  );

  // General page
  const renderGeneral = () => (
    <motion.div
      key="general"
      initial={{ opacity: 0, x: 20 }}
      animate={{ opacity: 1, x: 0 }}
      exit={{ opacity: 0 }}
      className="space-y-6"
    >
      <GlassCard title={t('settings.generalSettings.title')} padding="lg">
        <div className="space-y-6">
          {/* Language Selector */}
          <div>
            <label className="block text-sm font-medium text-slate-200 mb-2">
              {t('label.language')}
            </label>
            <div className="flex gap-3">
              <Button
                variant={language === 'en' ? 'primary' : 'ghost'}
                onClick={() => handleLanguageChange('en')}
                disabled={isProcessing}
                className="flex-1"
              >
                🇬🇧 English
              </Button>
              <Button
                variant={language === 'ru' ? 'primary' : 'ghost'}
                onClick={() => handleLanguageChange('ru')}
                disabled={isProcessing}
                className="flex-1"
              >
                🇷🇺 Русский
              </Button>
            </div>
          </div>

          {/* Autostart Toggle */}
          <div className="flex items-center justify-between p-4 bg-slate-700 rounded-lg">
            <div>
              <p className="text-slate-200 font-medium">{t('settings.autostart')}</p>
              <p className="text-sm text-slate-400 mt-1">
                {t('settings.generalSettings.autostartDescription')}
              </p>
            </div>
            <button
              onClick={() => handleAutoStartChange(!config?.autoStart)}
              disabled={isProcessing}
              className={`relative w-12 h-6 rounded-full transition-colors ${
                config?.autoStart ? 'bg-emerald-500' : 'bg-slate-600'
              } ${isProcessing ? 'opacity-50 cursor-not-allowed' : ''}`}
            >
              <motion.div
                animate={{ x: config?.autoStart ? 24 : 2 }}
                transition={{ type: 'spring', stiffness: 500, damping: 30 }}
                className="absolute top-1 w-4 h-4 bg-white rounded-full shadow"
              />
            </button>
          </div>
        </div>
      </GlassCard>
    </motion.div>
  );

  // Security page
  const renderSecurity = () => (
    <motion.div
      key="security"
      initial={{ opacity: 0, x: 20 }}
      animate={{ opacity: 1, x: 0 }}
      exit={{ opacity: 0 }}
      className="space-y-6"
    >
      <GlassCard title={t('settings.securitySettings.title')} padding="lg">
        <div className="space-y-6">
          {/* Change Password */}
          <div className="p-4 bg-slate-700 rounded-lg space-y-3">
            <div className="flex items-center gap-3">
              <div className="text-2xl">🔐</div>
              <div className="flex-1">
                <h3 className="text-slate-100 font-medium">{t('settings.securitySettings.changePassword')}</h3>
                <p className="text-sm text-slate-400 mt-1">
                  {t('settings.securitySettings.changePasswordDescription')}
                </p>
              </div>
            </div>
            <Button
              variant="secondary"
              onClick={() => setShowPasswordModal(true)}
              className="w-full"
            >
              {t('settings.securitySettings.changePassword')}
            </Button>
          </div>

          {/* Regenerate Keys */}
          <div className="p-4 bg-red-500/10 border border-red-500/30 rounded-lg space-y-3">
            <div className="flex items-center gap-3">
              <div className="text-2xl">⚠️</div>
              <div className="flex-1">
                <h3 className="text-slate-100 font-medium">{t('settings.securitySettings.regenerateKeys')}</h3>
                <p className="text-sm text-red-400 mt-1">
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

  // Storage page
  const renderStorage = () => (
    <motion.div
      key="storage"
      initial={{ opacity: 0, x: 20 }}
      animate={{ opacity: 1, x: 0 }}
      exit={{ opacity: 0 }}
      className="space-y-6"
    >
      <div className="mb-4">
        <h2 className="text-xl font-semibold text-slate-100">{t('storage.title')}</h2>
        <p className="text-sm text-slate-400 mt-1">{t('storage.subtitle')}</p>
      </div>

      <GlassCard padding="lg">
        <div className="space-y-4">
          <h3 className="text-lg font-semibold text-slate-100">{t('storage.limits')}</h3>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-slate-200 mb-2">
                {t('storage.maxMessageSize')}
              </label>
              <p className="text-xs text-slate-400 mb-4">
                {t('storage.maxMessageSizeDescription')}
              </p>
              <input
                type="range"
                min={10}
                max={500}
                step={10}
                value={maxMessageSize}
                onChange={(e) => setMaxMessageSize(Number(e.target.value))}
                className="w-full h-2 bg-slate-700 rounded-lg appearance-none cursor-pointer accent-emerald-500"
              />
              <div className="flex justify-between text-xs text-slate-400 mt-2">
                <span>10 {t('storage.mb')}</span>
                <span className="text-lg font-semibold text-emerald-400">
                  {maxMessageSize} {t('storage.mb')}
                </span>
                <span>500 {t('storage.mb')}</span>
              </div>
            </div>
            <Button
              variant="primary"
              onClick={handleMaxMessageSizeChange}
              disabled={isProcessing || !storageStats}
              className="w-full"
            >
              {t('action.save')}
            </Button>
          </div>
        </div>
      </GlassCard>
    </motion.div>
  );

  // Backup page
  const renderBackup = () => (
    <motion.div
      key="backup"
      initial={{ opacity: 0, x: 20 }}
      animate={{ opacity: 1, x: 0 }}
      exit={{ opacity: 0 }}
      className="space-y-6"
    >
      <div className="mb-4">
        <h2 className="text-xl font-semibold text-slate-100">{t('backup.title')}</h2>
        <p className="text-sm text-slate-400 mt-1">{t('backup.subtitle')}</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Create Backup */}
        <GlassCard title={t('backup.createBackup')} subtitle={t('backup.createBackupSubtitle')} padding="lg">
          {isCreating ? (
            <div className="space-y-4 py-6">
              <LoadingSpinner size="lg" text={createProgressMessage || t('backup.creating')} />
              <div className="space-y-2">
                <div className="flex justify-between text-sm text-slate-400">
                  <span>{createProgressMessage || t('backup.progress')}</span>
                  <span>{createProgress}%</span>
                </div>
                <div className="h-2 bg-slate-700 rounded-full overflow-hidden">
                  <motion.div
                    animate={{ width: `${createProgress}%` }}
                    className="h-full bg-emerald-500"
                  />
                </div>
              </div>
            </div>
          ) : (
            <div className="space-y-4">
              <div className="text-center mb-4">
                <div className="text-5xl mb-2">💾</div>
              </div>

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

              {/* Include Database Toggle */}
              <div className="flex items-center justify-between p-3 bg-slate-700 rounded-lg">
                <div className="flex-1">
                  <p className="text-slate-200 text-sm font-medium">{t('backup.includeDatabaseLabel')}</p>
                  <p className="text-xs text-slate-400 mt-1">{t('backup.includeDatabaseDescription')}</p>
                </div>
                <button
                  onClick={() => setIncludeDatabase(!includeDatabase)}
                  className={`relative w-12 h-6 rounded-full transition-colors ${
                    includeDatabase ? 'bg-emerald-500' : 'bg-slate-600'
                  }`}
                >
                  <motion.div
                    animate={{ x: includeDatabase ? 24 : 2 }}
                    transition={{ type: 'spring', stiffness: 500, damping: 30 }}
                    className="absolute top-1 w-4 h-4 bg-white rounded-full shadow"
                  />
                </button>
              </div>

              <Button variant="primary" size="lg" onClick={handleCreateBackup} className="w-full">
                {t('backup.createButton')}
              </Button>

              <div className="bg-emerald-500/10 border border-emerald-500/20 rounded-lg p-3 mt-4">
                <p className="text-sm text-slate-200 font-medium mb-2">{t('backup.whatsIncluded')}</p>
                <ul className="list-disc list-inside text-xs text-slate-400 space-y-1">
                  <li>{t('backup.includedItem1')}</li>
                  <li>{t('backup.includedItem2')}</li>
                  <li>{t('backup.includedItem3')}</li>
                  {includeDatabase && <li>{t('backup.includedItem4')}</li>}
                </ul>
              </div>
            </div>
          )}
        </GlassCard>

        {/* Restore Backup */}
        <GlassCard title={t('backup.restoreBackup')} subtitle={t('backup.restoreBackupSubtitle')} padding="lg">
          {isRestoring ? (
            <div className="space-y-4 py-6">
              <LoadingSpinner size="lg" text={restoreProgressMessage || t('backup.restoring')} />
              <div className="space-y-2">
                <div className="flex justify-between text-sm text-slate-400">
                  <span>{restoreProgressMessage || t('backup.progress')}</span>
                  <span>{restoreProgress}%</span>
                </div>
                <div className="h-2 bg-slate-700 rounded-full overflow-hidden">
                  <motion.div
                    animate={{ width: `${restoreProgress}%` }}
                    className="h-full bg-emerald-500"
                  />
                </div>
              </div>
            </div>
          ) : (
            <div className="space-y-4">
              <div className="text-center mb-4">
                <div className="text-5xl mb-2">📂</div>
              </div>

              {/* File Selection */}
              <div>
                <label className="block text-sm font-medium text-slate-200 mb-2">{t('backup.backupFile')}</label>
                <div className="bg-slate-700 rounded-lg p-3">
                  {restoreFilePath ? (
                    <div className="space-y-1">
                      <p className="text-xs text-slate-400 uppercase">{t('backup.selectedFile')}</p>
                      <p className="text-slate-200 font-mono text-sm break-all">{restoreFilePath}</p>
                    </div>
                  ) : (
                    <p className="text-slate-400 text-center py-2">{t('backup.noFileSelected')}</p>
                  )}
                </div>
                <Button variant="secondary" onClick={handleSelectBackupFile} className="w-full mt-2">
                  {t('backup.browseFiles')}
                </Button>
              </div>

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
                onClick={handleRestoreBackup}
                disabled={!restoreFilePath || !restorePassword}
                className="w-full"
              >
                {t('backup.restoreButton')}
              </Button>

              <div className="bg-red-500/10 border border-red-500/20 rounded-lg p-3 mt-4">
                <p className="text-sm text-slate-200 font-medium mb-2">{t('backup.warningTitle')}</p>
                <ul className="list-disc list-inside text-xs text-slate-400 space-y-1">
                  <li>{t('backup.warningItem1')}</li>
                  <li>{t('backup.warningItem2')}</li>
                  <li>{t('backup.warningItem3')}</li>
                </ul>
              </div>
            </div>
          )}
        </GlassCard>
      </div>

      {/* Best Practices */}
      <div className="bg-emerald-500/10 border border-emerald-500/20 rounded-xl p-5">
        <div className="flex items-start gap-4">
          <div className="text-3xl">💡</div>
          <div className="flex-1 space-y-2">
            <h3 className="text-slate-100 font-semibold">{t('backup.bestPractices')}</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
              <div>
                <p className="font-medium text-slate-200 mb-2">{t('backup.doTitle')}</p>
                <ul className="list-disc list-inside space-y-1 text-slate-400">
                  <li>{t('backup.doItem1')}</li>
                  <li>{t('backup.doItem2')}</li>
                  <li>{t('backup.doItem3')}</li>
                  <li>{t('backup.doItem4')}</li>
                </ul>
              </div>
              <div>
                <p className="font-medium text-slate-200 mb-2">{t('backup.dontTitle')}</p>
                <ul className="list-disc list-inside space-y-1 text-slate-400">
                  <li>{t('backup.dontItem1')}</li>
                  <li>{t('backup.dontItem2')}</li>
                  <li>{t('backup.dontItem3')}</li>
                  <li>{t('backup.dontItem4')}</li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>
    </motion.div>
  );

  // About page
  const renderAbout = () => (
    <motion.div
      key="about"
      initial={{ opacity: 0, x: 20 }}
      animate={{ opacity: 1, x: 0 }}
      exit={{ opacity: 0 }}
      className="space-y-6"
    >
      <GlassCard title={t('settings.aboutSettings.title')} padding="lg">
        <div className="space-y-6 text-center">
          <div className="flex justify-center mb-4">
            <img src="/appicon.png" alt="Tyr Desktop" className="w-24 h-24" />
          </div>
          <h2 className="text-2xl font-semibold text-emerald-400">Tyr Desktop</h2>
          <Badge variant="info" size="lg">Version {version || '2.0.0'}</Badge>
          <p className="text-slate-400 text-sm px-4">{t('settings.aboutSettings.description')}</p>
          <div className="pt-4 space-y-3 flex flex-col items-stretch px-4">
            <Button
              variant="primary"
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
          <div className="pt-4 border-t border-slate-700">
            <p className="text-xs text-slate-500 px-4">{t('settings.aboutSettings.madeWith')}</p>
          </div>
        </div>
      </GlassCard>
    </motion.div>
  );

  return (
    <div className="space-y-6">
      {/* Back Button */}
      {currentPage !== 'hub' && (
        <Button variant="ghost" onClick={() => setCurrentPage('hub')}>
          ← {t('settings.backToSettings')}
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
          <div className="flex gap-3 justify-end pt-2">
            <Button variant="ghost" onClick={() => setShowPasswordModal(false)}>
              {t('action.cancel')}
            </Button>
            <Button variant="primary" onClick={handleChangePassword} disabled={isProcessing}>
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
          <div className="bg-red-500/10 border border-red-500/30 rounded-lg p-4">
            <p className="text-slate-200 font-medium mb-2">{t('settings.securitySettings.warningTitle')}</p>
            <p className="text-sm text-slate-300">{t('settings.securitySettings.warningDescription')}</p>
            <ul className="list-disc list-inside text-sm text-slate-400 mt-2 space-y-1">
              <li>{t('settings.securitySettings.warningItem1')}</li>
              <li>{t('settings.securitySettings.warningItem2')}</li>
              <li>{t('settings.securitySettings.warningItem3')}</li>
              <li>{t('settings.securitySettings.warningItem4')}</li>
            </ul>
          </div>
          <p className="text-slate-300">{t('settings.securitySettings.confirmWarning')}</p>
          <Input
            label={t('settings.securitySettings.currentPassword')}
            type="password"
            placeholder={t('settings.securitySettings.enterCurrentPassword')}
            value={currentPassword}
            onChange={(e) => setCurrentPassword(e.target.value)}
          />
          <div className="flex gap-3 justify-end pt-2">
            <Button variant="ghost" onClick={() => setShowKeysModal(false)}>
              {t('action.cancel')}
            </Button>
            <Button variant="danger" onClick={handleRegenerateKeys} disabled={isProcessing || !currentPassword}>
              {isProcessing ? t('settings.securitySettings.regenerating') : t('settings.securitySettings.yesRegenerateKeys')}
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}

export default Settings;
