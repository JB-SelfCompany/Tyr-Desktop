import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';
import { motion } from 'framer-motion';
import {
  Button,
  Input,
  GlassCard,
  HolographicBorder,
  LoadingSpinner,
} from '../components';
import { useI18n } from '../hooks/useI18n';
import { useConfig } from '../hooks/useConfig';
import {
  CreateBackup,
  RestoreBackup,
  ShowSaveFileDialog,
  ShowOpenFileDialog,
} from '../../wailsjs/go/main/App';
import { toast } from '../components/ui/Toast';
// Import types from models
type BackupOptionsDTO = {
  includeDatabase: boolean;
  password: string;
};

type RestoreOptionsDTO = {
  backupPath: string;
  password: string;
};

/**
 * Backup Screen - Backup and restore functionality
 *
 * Features:
 * - Create encrypted backup
 *   - Include database checkbox
 *   - Password entry
 *   - Save file dialog
 * - Restore from backup
 *   - Open file dialog
 *   - Password entry
 *   - Confirmation
 * - Progress indicators
 * - Success/error notifications
 */
export function Backup() {
  const { t } = useI18n();
  const navigate = useNavigate();
  const { loadConfig } = useConfig({ loadOnMount: false });

  // Create backup state
  const [createPassword, setCreatePassword] = useState('');
  const [confirmCreatePassword, setConfirmCreatePassword] = useState('');
  const [includeDatabase, setIncludeDatabase] = useState(true);
  const [isCreating, setIsCreating] = useState(false);
  const [createProgress, setCreateProgress] = useState(0);

  // Restore backup state
  const [restoreFilePath, setRestoreFilePath] = useState('');
  const [restorePassword, setRestorePassword] = useState('');
  const [isRestoring, setIsRestoring] = useState(false);
  const [restoreProgress, setRestoreProgress] = useState(0);

  // Validation errors
  const [createPasswordError, setCreatePasswordError] = useState('');

  // Progress messages
  const [createProgressMessage, setCreateProgressMessage] = useState('');
  const [restoreProgressMessage, setRestoreProgressMessage] = useState('');

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
      // Configuration was restored successfully
      // Reload config from backend and navigate to dashboard
      await loadConfig();
      setIsRestoring(false);
      setRestoreProgress(0);
      setRestoreProgressMessage('');
      toast.success(t('backup.messages.backupRestoredMessage'));
      // Navigate to dashboard to see restored settings
      navigate('/');
    };

    EventsOn('backup:progress', backupProgressHandler);
    EventsOn('restore:progress', restoreProgressHandler);
    EventsOn('config:restored', configRestoredHandler);

    return () => {
      EventsOff('backup:progress');
      EventsOff('restore:progress');
      EventsOff('config:restored');
    };
  }, []);

  // Handle create backup
  const handleCreateBackup = async () => {
    // Validate password
    if (createPassword.length < 8) {
      setCreatePasswordError(t('backup.messages.passwordShort'));
      return;
    }
    if (createPassword !== confirmCreatePassword) {
      setCreatePasswordError(t('backup.messages.passwordMismatch'));
      return;
    }
    setCreatePasswordError('');

    // Show save file dialog
    try {
      // Generate filename in format: tbackup-dd-mm-yy.tb
      const now = new Date();
      const day = String(now.getDate()).padStart(2, '0');
      const month = String(now.getMonth() + 1).padStart(2, '0');
      const year = String(now.getFullYear()).slice(-2);
      const defaultFilename = `tbackup-${day}-${month}-${year}.tb`;

      const savePath = await ShowSaveFileDialog(t('backup.messages.saveBackup'), defaultFilename);
      if (!savePath) {
        return; // User cancelled
      }

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
    // Validate
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

      // Reset restore form state
      setRestoreFilePath('');
      setRestorePassword('');

      // NOTE: config:restored event handler will:
      // - Reload config from backend
      // - Navigate to dashboard
      // - Reset loading states
      // - Show success message
      // This avoids full page reload which causes system tray re-initialization issues
    } catch (error) {
      setIsRestoring(false);
      setRestoreProgress(0);
      setRestoreProgressMessage('');
      toast.error(error instanceof Error ? error.message : t('backup.messages.restoreFailedMessage'));
    }
  };

  return (
    <div className="space-y-6 pb-6">
      {/* Header */}
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, type: 'spring' }}
        className="text-center space-y-2"
      >
        <h1 className="text-5xl font-display font-black text-transparent bg-clip-text bg-iridescent bg-[length:200%_100%] animate-holographic-spin">
          {t('backup.title')}
        </h1>
        <p className="text-lg text-white/70 font-futuristic">
          {t('backup.subtitle')}
        </p>
      </motion.div>

      {/* Main Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Create Backup */}
        <motion.div
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.5, delay: 0.1 }}
        >
          <HolographicBorder animated={isCreating} borderWidth={isCreating ? 3 : 2}>
            <GlassCard
              title={t('backup.createBackup')}
              subtitle={t('backup.createBackupSubtitle')}
              padding="lg"
              variant="strong"
            >
              {isCreating ? (
                <div className="space-y-6 py-8">
                  <LoadingSpinner size="xl" variant="holographic" text={createProgressMessage || t('backup.creating')} />
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm text-white/70">
                      <span>{createProgressMessage || t('backup.progress')}</span>
                      <span>{createProgress}%</span>
                    </div>
                    <div className="h-2 bg-space-blue-dark rounded-full overflow-hidden">
                      <motion.div
                        animate={{ width: `${createProgress}%` }}
                        transition={{ duration: 0.3 }}
                        className="h-full bg-gradient-to-r from-neon-pink via-neon-cyan to-neon-green"
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
                    <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 p-4 bg-space-blue/30 rounded-lg border border-white/10">
                      <div className="flex-1">
                        <p className="text-white font-medium">{t('backup.includeDatabaseLabel')}</p>
                        <p className="text-sm text-white/60 mt-1">
                          {t('backup.includeDatabaseDescription')}
                        </p>
                      </div>
                      <button
                        onClick={() => setIncludeDatabase(!includeDatabase)}
                        className={`relative w-14 h-8 rounded-full transition-colors flex-shrink-0 ${
                          includeDatabase ? 'bg-neon-green' : 'bg-white/20'
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

                  <div className="bg-neon-cyan/10 border border-neon-cyan/30 rounded-lg p-4 mt-6">
                    <p className="text-sm text-white/90">
                      <strong>{t('backup.whatsIncluded')}</strong>
                    </p>
                    <ul className="list-disc list-inside text-sm text-white/70 mt-2 space-y-1">
                      <li>{t('backup.includedItem1')}</li>
                      <li>{t('backup.includedItem2')}</li>
                      <li>{t('backup.includedItem3')}</li>
                      {includeDatabase && <li>{t('backup.includedItem4')}</li>}
                    </ul>
                    <p className="text-xs text-white/60 mt-3">
                      {t('backup.encryptionInfo')}
                    </p>
                  </div>
                </div>
              )}
            </GlassCard>
          </HolographicBorder>
        </motion.div>

        {/* Restore Backup */}
        <motion.div
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.5, delay: 0.2 }}
        >
          <HolographicBorder animated={isRestoring} borderWidth={isRestoring ? 3 : 2}>
            <GlassCard
              title={t('backup.restoreBackup')}
              subtitle={t('backup.restoreBackupSubtitle')}
              padding="lg"
              variant="strong"
            >
              {isRestoring ? (
                <div className="space-y-6 py-8">
                  <LoadingSpinner size="xl" variant="holographic" text={restoreProgressMessage || t('backup.restoring')} />
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm text-white/70">
                      <span>{restoreProgressMessage || t('backup.progress')}</span>
                      <span>{restoreProgress}%</span>
                    </div>
                    <div className="h-2 bg-space-blue-dark rounded-full overflow-hidden">
                      <motion.div
                        animate={{ width: `${restoreProgress}%` }}
                        transition={{ duration: 0.3 }}
                        className="h-full bg-gradient-to-r from-neon-green via-neon-cyan to-neon-pink"
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
                      <label className="block text-sm font-medium text-white/90 mb-2">
                        {t('backup.backupFile')}
                      </label>
                      <div className="bg-space-blue/30 border border-white/20 rounded-lg p-4">
                        {restoreFilePath ? (
                          <div className="space-y-2">
                            <p className="text-xs text-white/50 uppercase tracking-wide">{t('backup.selectedFile')}</p>
                            <p className="text-white font-mono text-sm break-all">{restoreFilePath}</p>
                          </div>
                        ) : (
                          <p className="text-white/50 text-center py-2">{t('backup.noFileSelected')}</p>
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

                  <div className="bg-red-500/10 border border-red-500/30 rounded-lg p-4 mt-6">
                    <p className="text-sm text-white/90">
                      <strong>{t('backup.warningTitle')}</strong>
                    </p>
                    <ul className="list-disc list-inside text-sm text-white/70 mt-2 space-y-1">
                      <li>{t('backup.warningItem1')}</li>
                      <li>{t('backup.warningItem2')}</li>
                      <li>{t('backup.warningItem3')}</li>
                    </ul>
                  </div>
                </div>
              )}
            </GlassCard>
          </HolographicBorder>
        </motion.div>
      </div>

      {/* Info Section */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.3 }}
      >
        <HolographicBorder borderWidth={1}>
          <div className="bg-neon-pink/10 border border-neon-pink/30 rounded-lg p-6">
            <div className="flex items-start gap-4">
              <div className="text-4xl">💡</div>
              <div className="flex-1 space-y-2">
                <h3 className="text-white font-semibold text-lg">{t('backup.bestPractices')}</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm text-white/80">
                  <div>
                    <p className="font-medium mb-2">{t('backup.doTitle')}</p>
                    <ul className="list-disc list-inside space-y-1 text-white/70">
                      <li>{t('backup.doItem1')}</li>
                      <li>{t('backup.doItem2')}</li>
                      <li>{t('backup.doItem3')}</li>
                      <li>{t('backup.doItem4')}</li>
                    </ul>
                  </div>
                  <div>
                    <p className="font-medium mb-2">{t('backup.dontTitle')}</p>
                    <ul className="list-disc list-inside space-y-1 text-white/70">
                      <li>{t('backup.dontItem1')}</li>
                      <li>{t('backup.dontItem2')}</li>
                      <li>{t('backup.dontItem3')}</li>
                      <li>{t('backup.dontItem4')}</li>
                    </ul>
                  </div>
                </div>
                <div className="pt-3 mt-3 border-t border-white/10">
                  <p className="text-xs text-white/60">
                    {t('backup.encryptionNotice')}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </HolographicBorder>
      </motion.div>
    </div>
  );
}

export default Backup;
