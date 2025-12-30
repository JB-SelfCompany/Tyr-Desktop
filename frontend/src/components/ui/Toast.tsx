import { Toaster, toast as hotToast, ToastOptions } from 'react-hot-toast';

/**
 * Toast Provider Component
 *
 * Wraps react-hot-toast with Y2K Futurism styling
 *
 * Features:
 * - Glassmorphism backgrounds
 * - Holographic borders
 * - Neon glow effects
 * - Spring animations
 * - Auto-dismiss after 4 seconds
 */
export const ToastProvider = () => {
  return (
    <Toaster
      position="bottom-center"
      toastOptions={{
        // Default options
        duration: 2000,
        style: {
          background: 'rgba(26, 31, 54, 0.9)',
          backdropFilter: 'blur(20px)',
          color: '#fff',
          border: '1px solid rgba(255, 255, 255, 0.1)',
          borderRadius: '12px',
          padding: '16px',
          fontSize: '14px',
          fontFamily: 'Inter, system-ui, sans-serif',
          boxShadow: '0 8px 32px rgba(0, 0, 0, 0.3)',
        },

        // Success
        success: {
          duration: 2000,
          style: {
            background: 'rgba(16, 185, 129, 0.15)',
            backdropFilter: 'blur(20px)',
            border: '1px solid rgba(16, 185, 129, 0.5)',
            boxShadow: '0 0 20px rgba(16, 185, 129, 0.3), 0 8px 32px rgba(0, 0, 0, 0.3)',
          },
          iconTheme: {
            primary: '#10b981',
            secondary: '#ffffff',
          },
        },

        // Error
        error: {
          duration: 2000,
          style: {
            background: 'rgba(239, 68, 68, 0.15)',
            backdropFilter: 'blur(20px)',
            border: '1px solid rgba(239, 68, 68, 0.5)',
            boxShadow: '0 0 20px rgba(239, 68, 68, 0.3), 0 8px 32px rgba(0, 0, 0, 0.3)',
          },
          iconTheme: {
            primary: '#ef4444',
            secondary: '#ffffff',
          },
        },

        // Loading
        loading: {
          style: {
            background: 'rgba(99, 102, 241, 0.15)',
            backdropFilter: 'blur(20px)',
            border: '1px solid rgba(99, 102, 241, 0.5)',
            boxShadow: '0 0 20px rgba(99, 102, 241, 0.3), 0 8px 32px rgba(0, 0, 0, 0.3)',
          },
          iconTheme: {
            primary: '#6366f1',
            secondary: '#ffffff',
          },
        },
      }}
    />
  );
};

/**
 * Toast Utility Functions
 *
 * Convenient wrappers around react-hot-toast with consistent styling
 */
export const toast = {
  success: (message: string, options?: ToastOptions) => {
    return hotToast.success(message, options);
  },

  error: (message: string, options?: ToastOptions) => {
    return hotToast.error(message, options);
  },

  loading: (message: string, options?: ToastOptions) => {
    return hotToast.loading(message, options);
  },

  promise: <T,>(
    promise: Promise<T>,
    messages: {
      loading: string;
      success: string | ((data: T) => string);
      error: string | ((err: Error) => string);
    },
    options?: ToastOptions
  ) => {
    return hotToast.promise(promise, messages, options);
  },

  custom: (message: string, options?: ToastOptions) => {
    return hotToast(message, options);
  },

  dismiss: (toastId?: string) => {
    return hotToast.dismiss(toastId);
  },
};

/**
 * Example Usage:
 *
 * import { toast } from '@/components/ui/Toast';
 *
 * // Simple success
 * toast.success('Service started successfully!');
 *
 * // Error with custom duration
 * toast.error('Failed to connect to peer', { duration: 6000 });
 *
 * // Loading state
 * const toastId = toast.loading('Starting service...');
 * // Later dismiss it
 * toast.dismiss(toastId);
 *
 * // Promise-based (automatic loading -> success/error)
 * toast.promise(
 *   startService(),
 *   {
 *     loading: 'Starting service...',
 *     success: 'Service started!',
 *     error: (err) => `Failed: ${err.message}`,
 *   }
 * );
 */
