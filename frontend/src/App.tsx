import { useEffect, useState, lazy, Suspense } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import './style.css';

// Import Layout
import { Layout } from './components/layout';
import { LoadingSpinner, ToastProvider, ErrorBoundary } from './components';

// Lazy load screens for better performance
const Dashboard = lazy(() => import('./screens/Dashboard'));
const Onboarding = lazy(() => import('./screens/Onboarding'));
const Settings = lazy(() => import('./screens/Settings'));
const Peers = lazy(() => import('./screens/Peers'));
const Logs = lazy(() => import('./screens/Logs'));

// Import hooks
import { useEventStream } from './hooks/useEventStream';
import { useUIStore } from './store/uiStore';
import { useThemeManager } from './hooks/useThemeManager';
import { useI18n } from './hooks/useI18n';

// Import Wails bindings
import { IsOnboardingComplete } from '../wailsjs/go/main/App';
import { LogPrint } from './wailsjs/runtime/runtime';

/**
 * App Component - Main application root
 *
 * Features:
 * - React Router navigation
 * - Onboarding flow check
 * - Event stream initialization
 * - Layout wrapper with sidebar
 * - Loading states
 * - Dark mode management
 */
function App() {
  const [onboardingComplete, setOnboardingComplete] = useState<boolean>(false);
  const [isCheckingOnboarding, setIsCheckingOnboarding] = useState(true);

  // Initialize event stream for real-time updates
  useEventStream();

  // Initialize theme management with system detection
  useThemeManager();

  // Initialize i18n synchronization
  useI18n();

  // UI state
  const isAppLoading = useUIStore((state) => state.isAppLoading);
  const setAppLoading = useUIStore((state) => state.setAppLoading);

  // Function to complete onboarding and update state
  const completeOnboarding = async () => {
    LogPrint('[App] completeOnboarding called');

    // Double-check with backend that onboarding is actually complete
    try {
      const complete = await IsOnboardingComplete();
      LogPrint('[App] Backend onboarding status after completion: ' + complete);
      setOnboardingComplete(complete);

      if (!complete) {
        LogPrint('[App] ERROR: Backend says onboarding is NOT complete!');
      } else {
        LogPrint('[App] SUCCESS: Onboarding confirmed complete, UI will update');
      }
    } catch (err) {
      LogPrint('[App] Failed to verify onboarding status: ' + (err instanceof Error ? err.message : String(err)));
      // Optimistically set to true anyway
      setOnboardingComplete(true);
    }
  };

  useEffect(() => {
    // Check onboarding status on mount
    const checkOnboarding = async () => {
      LogPrint('[App] Checking onboarding status...');
      try {
        const complete = await IsOnboardingComplete();
        LogPrint('[App] Onboarding status from backend: ' + complete);
        setOnboardingComplete(complete);
      } catch (err) {
        LogPrint('[App] Failed to check onboarding status: ' + (err instanceof Error ? err.message : String(err)));
        // If check fails, assume onboarding is needed
        setOnboardingComplete(false);
      } finally {
        setIsCheckingOnboarding(false);
        setAppLoading(false);
        LogPrint('[App] checkOnboarding finished');
      }
    };

    checkOnboarding();
  }, [setAppLoading]);

  // Show loading screen while checking onboarding
  if (isCheckingOnboarding || isAppLoading) {
    LogPrint(`[App] Showing loading screen - isCheckingOnboarding: ${isCheckingOnboarding}, isAppLoading: ${isAppLoading}`);
    return (
      <>
        <ToastProvider />
        <div className="min-h-screen bg-gradient-to-br from-space-blue via-space-blue-light to-space-blue flex items-center justify-center">
          <LoadingSpinner size="xl" variant="holographic" text="Loading Tyr Desktop..." fullScreen />
        </div>
      </>
    );
  }

  // Loading fallback component
  const ScreenLoadingFallback = (
    <div className="flex items-center justify-center min-h-[400px]">
      <LoadingSpinner size="lg" variant="holographic" text="Loading..." />
    </div>
  );

  // Show onboarding if not complete
  LogPrint('[App] Render decision - onboardingComplete: ' + onboardingComplete);
  if (!onboardingComplete) {
    LogPrint('[App] Showing onboarding screen');
    return (
      <>
        <ToastProvider />
        <ErrorBoundary>
          <Router>
            <Suspense fallback={ScreenLoadingFallback}>
              <Routes>
                <Route path="/onboarding" element={<Onboarding onComplete={completeOnboarding} />} />
                <Route path="*" element={<Navigate to="/onboarding" replace />} />
              </Routes>
            </Suspense>
          </Router>
        </ErrorBoundary>
      </>
    );
  }

  // Show main app with navigation
  LogPrint('[App] Showing main app (Dashboard)');
  return (
    <>
      <ToastProvider />
      <ErrorBoundary>
        <Router>
          <Layout>
            <Suspense fallback={ScreenLoadingFallback}>
              <Routes>
                <Route path="/" element={<Dashboard />} />
                <Route path="/peers" element={<Peers />} />
                <Route path="/logs" element={<Logs />} />
                <Route path="/settings/*" element={<Settings />} />
                <Route path="*" element={<Navigate to="/" replace />} />
              </Routes>
            </Suspense>
          </Layout>
        </Router>
      </ErrorBoundary>
    </>
  );
}

export default App;
