import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
  Button,
  Input,
  Modal,
  GlassCard,
  BentoGrid,
  BentoCard,
  HolographicBorder,
  Badge,
} from '../components';
import { useConfig } from '../hooks/useConfig';
import { useI18n } from '../hooks/useI18n';
import {
  SetTheme,
  SetAutoStart,
  OpenURL,
  GetVersion,
  ChangePassword,
  RegenerateKeys,
} from '../../wailsjs/go/main/App';
import { showSuccess, showError } from '../store/uiStore';

type SettingsPage = 'hub' | 'general' | 'security' | 'about';

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

  // Load version on component mount
  useEffect(() => {
    GetVersion().then(setVersion).catch(() => setVersion('unknown'));
  }, []);

  // Handle language change
  const handleLanguageChange = async (lang: 'en' | 'ru') => {
    try {
      setIsProcessing(true);
      // changeLanguage will handle both i18n update and backend config save
      await changeLanguage(lang);
      showSuccess(t('dialog.success'), t('settings.messages.languageChanged'));
    } catch (error) {
      showError(t('dialog.error'), error instanceof Error ? error.message : t('settings.messages.languageChangeFailed'));
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
      showSuccess(t('dialog.success'), t('settings.messages.themeChanged'));
    } catch (error) {
      showError(t('dialog.error'), error instanceof Error ? error.message : t('settings.messages.themeChangeFailed'));
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
      showSuccess(t('dialog.success'), enabled ? t('settings.messages.autostartEnabled') : t('settings.messages.autostartDisabled'));
    } catch (error) {
      showError(t('dialog.error'), error instanceof Error ? error.message : t('settings.messages.autostartChangeFailed'));
    } finally {
      setIsProcessing(false);
    }
  };

  // Handle password change
  const handleChangePassword = async () => {
    if (!currentPassword) {
      showError(t('settings.messages.validationError'), t('settings.messages.passwordRequired'));
      return;
    }
    if (newPassword.length < 8) {
      showError(t('settings.messages.validationError'), t('settings.messages.passwordShort'));
      return;
    }
    if (newPassword !== confirmNewPassword) {
      showError(t('settings.messages.validationError'), t('settings.messages.passwordMismatch'));
      return;
    }
    setIsProcessing(true);
    try {
      await ChangePassword(currentPassword, newPassword);
      showSuccess(t('settings.messages.passwordChangeTitle'), t('settings.messages.passwordChanged'));
      setShowPasswordModal(false);
      setNewPassword('');
      setConfirmNewPassword('');
      setCurrentPassword('');
    } catch (error) {
      showError(t('settings.messages.passwordChangeFailedTitle'), error instanceof Error ? error.message : t('settings.messages.passwordChangeFailed'));
    } finally {
      setIsProcessing(false);
    }
  };

  // Handle key regeneration
  const handleRegenerateKeys = async () => {
    if (!currentPassword) {
      showError(t('settings.messages.passwordRequiredTitle'), t('settings.messages.passwordRequired'));
      return;
    }
    setIsProcessing(true);
    try {
      await RegenerateKeys(currentPassword);
      showSuccess(t('settings.messages.keysRegenerated'), t('settings.messages.keysRegeneratedMessage'));
      setShowKeysModal(false);
      setCurrentPassword('');
    } catch (error) {
      showError(t('settings.messages.keyRegenerationFailedTitle'), error instanceof Error ? error.message : t('settings.messages.keyRegenerationFailed'));
    } finally {
      setIsProcessing(false);
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
