import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';
import { motion } from 'framer-motion';
import {
  Button,
  Input,
  GlassCard,
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

type BackupOptionsDTO = {
  backupPath: string;
  includeDatabase: boolean;
  password: string;
};

type RestoreOptionsDTO = {
  backupPath: string;
  password: string;
};

/**
 * Backup Screen - Backup and restore functionality
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
      try {
        await loadConfig();
      } catch (error) {
        console.error('Failed to reload config after restore:', error);
      }
      setIsRestoring(false);
      setRestoreProgress(0);
      setRestoreProgressMessage('');
      toast.success(t('backup.messages.backupRestoredMessage'));
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
      if (!savePath) {
        return;
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
      const result = await RestoreBackup(options);

      if (!result.success) {
        throw new Error(result.message || t('backup.messages.restoreFailedMessage'));
      }

      setRestoreFilePath('');
      setRestorePassword('');
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
        initial={{ opacity: 0, y: -10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.2 }}
        className="text-center"
      >
        <h1 className="text-2xl font-semibold text-slate-100">
          {t('backup.title')}
        </h1>
        <p className="text-sm text-slate-400 mt-1">
          {t('backup.subtitle')}
        </p>
      </motion.div>

      {/* Main Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Create Backup */}
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.2, delay: 0.05 }}
        >
          <GlassCard
            title={t('backup.createBackup')}
            subtitle={t('backup.createBackupSubtitle')}
            padding="lg"
          >
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
                      transition={{ duration: 0.3 }}
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
        </motion.div>

        {/* Restore Backup */}
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.2, delay: 0.1 }}
        >
          <GlassCard
            title={t('backup.restoreBackup')}
            subtitle={t('backup.restoreBackupSubtitle')}
            padding="lg"
          >
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
                      transition={{ duration: 0.3 }}
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
        </motion.div>
      </div>

      {/* Best Practices */}
      <motion.div
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.2, delay: 0.15 }}
      >
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
    </div>
  );
}

export default Backup;
