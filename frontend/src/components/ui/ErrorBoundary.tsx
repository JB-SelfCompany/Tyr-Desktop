import React, { Component, ErrorInfo, ReactNode } from 'react';
import { Button } from './Button';
import { GlassCard } from '../layout/GlassCard';

interface Props {
  children: ReactNode;
  fallback?: ReactNode;
}

interface State {
  hasError: boolean;
  error: Error | null;
  errorInfo: ErrorInfo | null;
}

/**
 * ErrorBoundary Component
 *
 * Catches JavaScript errors anywhere in the child component tree,
 * logs those errors, and displays a fallback UI.
 *
 * Features:
 * - Catches rendering errors
 * - Shows user-friendly error message with Y2K styling
 * - Provides reload button
 * - Logs error details to console
 *
 * Usage:
 * ```tsx
 * <ErrorBoundary>
 *   <YourComponent />
 * </ErrorBoundary>
 * ```
 *
 * With custom fallback:
 * ```tsx
 * <ErrorBoundary fallback={<CustomErrorUI />}>
 *   <YourComponent />
 * </ErrorBoundary>
 * ```
 */
export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {
      hasError: false,
      error: null,
      errorInfo: null,
    };
  }

  static getDerivedStateFromError(error: Error): State {
    // Update state so the next render will show the fallback UI
    return {
      hasError: true,
      error,
      errorInfo: null,
    };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    // Log error details to console
    console.error('ErrorBoundary caught an error:', error, errorInfo);

    // Update state with error info
    this.setState({
      error,
      errorInfo,
    });

    // You can also log the error to an error reporting service here
    // logErrorToService(error, errorInfo);
  }

  handleReload = () => {
    // Clear error state and reload
    this.setState({
      hasError: false,
      error: null,
      errorInfo: null,
    });

    // Optionally, reload the entire page
    // window.location.reload();
  };

  render() {
    if (this.state.hasError) {
      // If custom fallback is provided, use it
      if (this.props.fallback) {
        return this.props.fallback;
      }

      // Default error UI
      return (
        <div className="min-h-screen bg-slate-900 flex items-center justify-center p-6">
          <GlassCard className="max-w-2xl w-full" padding="xl">
            <div className="text-center space-y-6">
              {/* Error Icon */}
              <div className="flex justify-center">
                <div className="relative">
                  <div className="absolute inset-0 bg-red-500/30 blur-xl rounded-full animate-pulse" />
                  <svg
                    className="w-20 h-20 text-red-400 relative z-10"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    aria-hidden="true"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                    />
                  </svg>
                </div>
              </div>

              {/* Error Title */}
              <h1 className="text-3xl font-bold text-red-400">
                Oops! Something went wrong
              </h1>

              {/* Error Message */}
              <p className="text-slate-300 text-lg">
                We're sorry, but an unexpected error occurred. Don't worry, your data is safe.
              </p>

              {/* Error Details (Expandable) */}
              {this.state.error && (
                <details className="text-left mt-4">
                  <summary className="cursor-pointer text-slate-400 hover:text-slate-200 transition-colors text-sm font-mono mb-2">
                    Show error details
                  </summary>
                  <div className="bg-red-500/10 border border-red-500/30 rounded-xl p-4 text-xs font-mono text-red-400 overflow-x-auto">
                    <div className="mb-2">
                      <strong>Error:</strong> {this.state.error.toString()}
                    </div>
                    {this.state.errorInfo && (
                      <div>
                        <strong>Stack trace:</strong>
                        <pre className="mt-2 text-xs whitespace-pre-wrap">
                          {this.state.errorInfo.componentStack}
                        </pre>
                      </div>
                    )}
                  </div>
                </details>
              )}

              {/* Action Buttons */}
              <div className="flex gap-4 justify-center pt-4">
                <Button
                  variant="primary"
                  onClick={this.handleReload}
                  className="min-w-[140px]"
                  aria-label="Try again"
                >
                  Try Again
                </Button>

                <Button
                  variant="secondary"
                  onClick={() => window.location.reload()}
                  className="min-w-[140px]"
                  aria-label="Reload page"
                >
                  Reload Page
                </Button>
              </div>
            </div>
          </GlassCard>
        </div>
      );
    }

    // No error, render children normally
    return this.props.children;
  }
}

/**
 * Error Boundary Hook (for functional components wrapper)
 *
 * This is a convenience function to wrap components with ErrorBoundary
 *
 * Usage:
 * ```tsx
 * export default withErrorBoundary(MyComponent);
 * ```
 */
export const withErrorBoundary = <P extends object>(
  Component: React.ComponentType<P>,
  fallback?: ReactNode
) => {
  return (props: P) => (
    <ErrorBoundary fallback={fallback}>
      <Component {...props} />
    </ErrorBoundary>
  );
};
