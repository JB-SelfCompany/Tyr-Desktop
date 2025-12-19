import React from 'react';
import { motion } from 'framer-motion';

export type SpinnerSize = 'sm' | 'md' | 'lg' | 'xl';
export type SpinnerVariant = 'default' | 'neon' | 'holographic';

interface LoadingSpinnerProps {
  size?: SpinnerSize;
  variant?: SpinnerVariant;
  text?: string;
  fullScreen?: boolean;
}

const sizeClasses: Record<SpinnerSize, string> = {
  sm: 'w-6 h-6',
  md: 'w-10 h-10',
  lg: 'w-16 h-16',
  xl: 'w-24 h-24',
};

const textSizeClasses: Record<SpinnerSize, string> = {
  sm: 'text-sm',
  md: 'text-base',
  lg: 'text-lg',
  xl: 'text-xl',
};

export const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({
  size = 'md',
  variant = 'default',
  text,
  fullScreen = false,
}) => {
  const sizeClass = sizeClasses[size];
  const textSizeClass = textSizeClasses[size];

  const spinnerContent = (
    <div className="flex flex-col items-center gap-4">
      {/* Spinner */}
      <div className="relative">
        {variant === 'neon' ? (
          // Neon variant - glowing rings
          <motion.div
            className={`${sizeClass} relative`}
            animate={{ rotate: 360 }}
            transition={{ duration: 1.5, repeat: Infinity, ease: 'linear' }}
          >
            <div className="absolute inset-0 rounded-full border-4 border-neon-pink/30" />
            <div className="absolute inset-0 rounded-full border-t-4 border-neon-pink shadow-neon-pink" />
          </motion.div>
        ) : variant === 'holographic' ? (
          // Holographic variant - rainbow gradient with proper rotation
          <motion.div
            className={`${sizeClass} relative`}
            animate={{ rotate: 360 }}
            transition={{ duration: 2, repeat: Infinity, ease: 'linear' }}
          >
            <svg className="w-full h-full" viewBox="0 0 50 50">
              <defs>
                <linearGradient id="holographic-gradient" x1="0%" y1="0%" x2="100%" y2="100%">
                  <stop offset="0%" stopColor="#FF00FF" />
                  <stop offset="33%" stopColor="#00FFFF" />
                  <stop offset="66%" stopColor="#00FF00" />
                  <stop offset="100%" stopColor="#FFFF00" />
                </linearGradient>
              </defs>
              <circle
                cx="25"
                cy="25"
                r="20"
                fill="none"
                stroke="rgba(255,255,255,0.1)"
                strokeWidth="4"
              />
              <circle
                cx="25"
                cy="25"
                r="20"
                fill="none"
                stroke="url(#holographic-gradient)"
                strokeWidth="4"
                strokeLinecap="round"
                strokeDasharray="90 150"
              />
            </svg>
          </motion.div>
        ) : (
          // Default variant - simple spinner
          <motion.div
            className={`${sizeClass}`}
            animate={{ rotate: 360 }}
            transition={{ duration: 1, repeat: Infinity, ease: 'linear' }}
          >
            <svg
              className="w-full h-full text-md-light-primary dark:text-md-dark-primary"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                className="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                strokeWidth="4"
              />
              <path
                className="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              />
            </svg>
          </motion.div>
        )}

        {/* Pulsing background glow */}
        {variant !== 'default' && (
          <motion.div
            className={`absolute inset-0 rounded-full blur-xl ${
              variant === 'neon' ? 'bg-neon-pink/30' : 'bg-neon-cyan/30'
            }`}
            animate={{
              opacity: [0.3, 0.6, 0.3],
              scale: [0.8, 1.2, 0.8],
            }}
            transition={{
              duration: 2,
              repeat: Infinity,
              ease: 'easeInOut',
            }}
          />
        )}
      </div>

      {/* Loading text */}
      {text && (
        <motion.p
          className={`${textSizeClass} font-futuristic font-semibold text-md-light-onSurface dark:text-md-dark-onSurface`}
          animate={{ opacity: [0.5, 1, 0.5] }}
          transition={{
            duration: 1.5,
            repeat: Infinity,
            ease: 'easeInOut',
          }}
        >
          {text}
        </motion.p>
      )}
    </div>
  );

  if (fullScreen) {
    return (
      <div className="fixed inset-0 z-50 flex items-center justify-center bg-space-blue/80 backdrop-blur-lg">
        {spinnerContent}
      </div>
    );
  }

  return spinnerContent;
};

// Simple inline spinner (for buttons, etc.)
export const InlineSpinner: React.FC<{ className?: string }> = ({ className = '' }) => (
  <motion.svg
    className={`animate-spin ${className}`}
    xmlns="http://www.w3.org/2000/svg"
    fill="none"
    viewBox="0 0 24 24"
  >
    <circle
      className="opacity-25"
      cx="12"
      cy="12"
      r="10"
      stroke="currentColor"
      strokeWidth="4"
    />
    <path
      className="opacity-75"
      fill="currentColor"
      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
    />
  </motion.svg>
);
