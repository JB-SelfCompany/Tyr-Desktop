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

      // Default error UI with Y2K styling
      return (
        <div className="min-h-screen bg-gradient-to-br from-space-blue via-space-blue-light to-space-blue flex items-center justify-center p-6">
          <GlassCard className="max-w-2xl w-full">
            <div className="text-center space-y-6">
              {/* Error Icon with Holographic Glow */}
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
              <h1 className="text-3xl font-bold font-display bg-gradient-to-r from-red-400 via-pink-400 to-red-400 bg-clip-text text-transparent">
                Oops! Something went wrong
              </h1>

              {/* Error Message */}
              <p className="text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant text-lg font-body">
                We're sorry, but an unexpected error occurred. Don't worry, your data is safe.
              </p>

              {/* Error Details (Expandable) */}
              {this.state.error && (
                <details className="text-left mt-4">
                  <summary className="cursor-pointer text-md-light-outline dark:text-md-dark-outline hover:text-md-light-onSurfaceVariant dark:hover:text-md-dark-onSurfaceVariant transition-colors text-sm font-mono mb-2">
                    Show error details
                  </summary>
                  <div className="bg-md-light-errorContainer/30 dark:bg-md-dark-errorContainer/30 rounded-lg p-4 text-xs font-mono text-md-light-error dark:text-md-dark-error overflow-x-auto border border-md-light-error/30 dark:border-md-dark-error/30">
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
                  <svg
                    className="w-5 h-5 mr-2"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    aria-hidden="true"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
                    />
                  </svg>
                  Try Again
                </Button>

                <Button
                  variant="secondary"
                  onClick={() => window.location.reload()}
                  className="min-w-[140px]"
                  aria-label="Reload page"
                >
                  <svg
                    className="w-5 h-5 mr-2"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    aria-hidden="true"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
                    />
                  </svg>
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
