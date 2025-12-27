import { ReactNode, useState, useEffect } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { motion } from 'framer-motion';
import { HolographicBorder } from './HolographicBorder';
import { Modal } from '../ui/Modal';
import { Button } from '../ui/Button';
import { useI18n } from '../../hooks/useI18n';
import { GetVersion, QuitApplication } from '../../../wailsjs/go/main/App';

interface LayoutProps {
  children: ReactNode;
}

/**
 * Layout Component - Main application layout with navigation
 *
 * Features:
 * - Sidebar navigation with Y2K styling
 * - Active route highlighting
 * - Responsive design
 * - Smooth transitions
 */
export function Layout({ children }: LayoutProps) {
  const location = useLocation();
  const { t } = useI18n();
  const [version, setVersion] = useState('...');
  const [showExitModal, setShowExitModal] = useState(false);

  // Load version on component mount
  useEffect(() => {
    GetVersion().then(setVersion).catch(() => setVersion('unknown'));
  }, []);

  // Handle application exit
  const handleExit = () => {
    try {
      QuitApplication();
    } catch (error) {
      console.error('Failed to quit application:', error);
    }
  };

  const navItems = [
    { path: '/', label: t('navigation.dashboard'), icon: '🏠' },
    { path: '/peers', label: t('navigation.peers'), icon: '🌐' },
    { path: '/backup', label: t('navigation.backup'), icon: '💾' },
    { path: '/logs', label: t('navigation.logs'), icon: '📋' },
    { path: '/settings', label: t('navigation.settings'), icon: '⚙️' },
  ];

  const isActive = (path: string) => {
    if (path === '/') return location.pathname === '/';
    return location.pathname.startsWith(path);
  };

  return (
    <div className="h-screen bg-gradient-to-br from-md-light-background via-md-light-surface to-[#E8F0E8] dark:from-md-dark-background dark:via-md-dark-surface dark:to-[#0F120F] flex overflow-hidden">
      {/* Sidebar Navigation */}
      <motion.aside
        initial={{ x: -100, opacity: 0 }}
        animate={{ x: 0, opacity: 1 }}
        transition={{ duration: 0.5, type: 'spring' }}
        className="w-56 lg:w-64 flex-shrink-0 p-4 lg:p-6 space-y-4 lg:space-y-6 overflow-y-auto"
      >
        {/* Logo */}
        <div className="text-center space-y-2">
          <div className="flex justify-center mb-2">
            <img src="/appicon.png" alt="Tyr" className="w-16 h-16" />
          </div>
          <h2 className="text-2xl font-display font-black text-transparent bg-clip-text bg-iridescent bg-[length:200%_100%]">
            Tyr Desktop
          </h2>
          <p className="text-xs text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant font-futuristic">P2P Email Client</p>
        </div>

        {/* Navigation Links */}
        <nav className="space-y-2">
          {navItems.map((item) => {
            const active = isActive(item.path);
            return (
              <HolographicBorder
                key={item.path}
                animated={active}
                borderWidth={active ? 2 : 1}
              >
                <Link
                  to={item.path}
                  className={`flex items-center gap-3 px-4 py-3 rounded-lg transition-all ${
                    active
                      ? 'bg-md-light-primaryContainer dark:bg-md-dark-primaryContainer scale-105'
                      : 'bg-md-light-surface/60 dark:bg-md-dark-surface/50 hover:bg-md-light-primaryContainer/50 dark:hover:bg-md-dark-primaryContainer/30'
                  }`}
                >
                  <span className="text-2xl">{item.icon}</span>
                  <span className={`font-futuristic ${active ? 'text-md-light-onPrimaryContainer dark:text-md-dark-onPrimaryContainer font-semibold' : 'text-md-light-onSurface dark:text-md-dark-onSurface'}`}>
                    {item.label}
                  </span>
                </Link>
              </HolographicBorder>
            );
          })}
        </nav>

        {/* Exit Button */}
        <div className="pt-4">
          <HolographicBorder borderWidth={1}>
            <button
              onClick={() => setShowExitModal(true)}
              className="flex items-center gap-3 px-4 py-3 rounded-lg w-full transition-all bg-md-light-errorContainer/30 dark:bg-md-dark-errorContainer/20 hover:bg-md-light-errorContainer/50 dark:hover:bg-md-dark-errorContainer/30"
            >
              <span className="text-2xl">🚪</span>
              <span className="font-futuristic text-md-light-error dark:text-md-dark-error">
                {t('app.quit')}
              </span>
            </button>
          </HolographicBorder>
        </div>

        {/* Version Info */}
        <div className="pt-6 mt-auto border-t border-md-light-outline/30 dark:border-md-dark-outline/30">
          <p className="text-xs text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant text-center">Version {version}</p>
          <p className="text-xs text-md-light-outline dark:text-md-dark-outline text-center mt-1">Wails v2 + React</p>
        </div>
      </motion.aside>

      {/* Main Content Area */}
      <main className="flex-1 p-4 lg:p-6 overflow-y-auto overflow-x-hidden min-w-0">
        <motion.div
          key={location.pathname}
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: -20 }}
          transition={{ duration: 0.3 }}
          className="h-full"
        >
          <div className="max-w-[1600px] mx-auto">
            {children}
          </div>
        </motion.div>
      </main>

      {/* Exit Confirmation Modal */}
      <Modal
        isOpen={showExitModal}
        onClose={() => setShowExitModal(false)}
        title={t('settings.exitConfirmation.title')}
        size="sm"
      >
        <div className="space-y-4">
          <div className="bg-md-light-primaryContainer/50 dark:bg-md-dark-primaryContainer/30 border border-md-light-primary/30 dark:border-md-dark-primary/30 rounded-lg p-4">
            <p className="text-md-light-onSurface dark:text-md-dark-onSurface">
              {t('settings.exitConfirmation.message')}
            </p>
          </div>
          <div className="flex gap-3 justify-end pt-4">
            <Button variant="ghost" onClick={() => setShowExitModal(false)}>
              {t('action.cancel')}
            </Button>
            <Button
              variant="danger"
              glow
              onClick={() => {
                setShowExitModal(false);
                handleExit();
              }}
            >
              {t('settings.exitConfirmation.confirm')}
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}

export default Layout;
