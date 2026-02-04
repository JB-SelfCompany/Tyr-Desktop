import { ReactNode, useState, useEffect } from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import { motion } from 'framer-motion';
import { Modal } from '../ui/Modal';
import { Button } from '../ui/Button';
import { useI18n } from '../../hooks/useI18n';
import { GetVersion, QuitApplication } from '../../../wailsjs/go/main/App';

interface LayoutProps {
  children: ReactNode;
}

/**
 * Layout Component - Main application layout with navigation
 */
export function Layout({ children }: LayoutProps) {
  const location = useLocation();
  const { t } = useI18n();
  const [version, setVersion] = useState('...');
  const [showExitModal, setShowExitModal] = useState(false);

  useEffect(() => {
    GetVersion().then(setVersion).catch(() => setVersion('unknown'));
  }, []);

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
    { path: '/logs', label: t('navigation.logs'), icon: '📋' },
    { path: '/settings', label: t('navigation.settings'), icon: '⚙️' },
  ];

  return (
    <div className="flex h-screen bg-slate-900">
      {/* Sidebar Navigation */}
      <aside className="w-64 flex-shrink-0 glass border-r border-slate-700/50">
        <div className="flex flex-col h-full">
          {/* Logo */}
          <div className="p-6 border-b border-slate-700/50">
            <div className="flex flex-col items-center text-center space-y-2">
              <img src="/appicon.png" alt="Tyr" className="w-16 h-16" />
              <h1 className="text-xl font-bold text-emerald-400">
                Tyr Desktop
              </h1>
              <p className="text-xs text-slate-400">
                P2P Email Client
              </p>
            </div>
          </div>

          {/* Navigation Links */}
          <nav className="flex-1 p-4 space-y-1">
            {navItems.map((item) => (
              <NavLink
                key={item.path}
                to={item.path}
                end={item.path === '/'}
                className={({ isActive }) =>
                  `flex items-center gap-3 px-4 py-3 rounded-xl transition-all duration-200 ${
                    isActive
                      ? 'bg-emerald-500/20 text-emerald-400 font-medium'
                      : 'text-slate-300 hover:bg-slate-700/50 hover:text-slate-100'
                  }`
                }
              >
                <span className="text-lg">{item.icon}</span>
                <span>{item.label}</span>
              </NavLink>
            ))}
          </nav>

          {/* Exit Button */}
          <div className="p-4 border-t border-slate-700/50">
            <button
              onClick={() => setShowExitModal(true)}
              className="flex items-center gap-3 px-4 py-3 rounded-xl w-full transition-all duration-200 text-red-400 hover:bg-red-500/10"
            >
              <span className="text-lg">🚪</span>
              <span>{t('app.quit')}</span>
            </button>
          </div>

          {/* Version Info */}
          <div className="p-4 border-t border-slate-700/50">
            <p className="text-xs text-center text-slate-500">
              Tyr Desktop v{version}
            </p>
          </div>
        </div>
      </aside>

      {/* Main Content Area */}
      <main className="flex-1 flex flex-col min-h-0 overflow-hidden">
        <motion.div
          key={location.pathname}
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.2 }}
          className="flex-1 min-h-0 p-6 pb-6 overflow-y-auto overflow-x-hidden"
        >
          {children}
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
          <div className="bg-slate-700 rounded-xl p-4">
            <p className="text-slate-200">
              {t('settings.exitConfirmation.message')}
            </p>
          </div>
          <div className="flex gap-3 justify-end pt-2">
            <Button variant="ghost" onClick={() => setShowExitModal(false)}>
              {t('action.cancel')}
            </Button>
            <Button
              variant="danger"
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
