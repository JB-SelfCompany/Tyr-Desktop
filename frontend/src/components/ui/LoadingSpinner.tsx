import React from 'react';
import { motion } from 'framer-motion';

export type SpinnerSize = 'sm' | 'md' | 'lg' | 'xl';

interface LoadingSpinnerProps {
  size?: SpinnerSize;
  text?: string;
  fullScreen?: boolean;
}

const sizeClasses: Record<SpinnerSize, string> = {
  sm: 'w-5 h-5',
  md: 'w-8 h-8',
  lg: 'w-12 h-12',
  xl: 'w-16 h-16',
};

const borderSizes: Record<SpinnerSize, string> = {
  sm: 'border-2',
  md: 'border-[3px]',
  lg: 'border-4',
  xl: 'border-4',
};

const textSizeClasses: Record<SpinnerSize, string> = {
  sm: 'text-sm',
  md: 'text-sm',
  lg: 'text-base',
  xl: 'text-lg',
};

export const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({
  size = 'md',
  text,
  fullScreen = false,
}) => {
  const sizeClass = sizeClasses[size];
  const borderSize = borderSizes[size];
  const textSizeClass = textSizeClasses[size];

  const spinnerContent = (
    <div className="flex flex-col items-center gap-3">
      {/* Spinner */}
      <motion.div
        className={`${sizeClass} rounded-full ${borderSize} border-slate-700 border-t-emerald-500`}
        animate={{ rotate: 360 }}
        transition={{ duration: 0.8, repeat: Infinity, ease: 'linear' }}
      />

      {/* Loading text */}
      {text && (
        <p className={`${textSizeClass} text-slate-400`}>
          {text}
        </p>
      )}
    </div>
  );

  if (fullScreen) {
    return (
      <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/80 backdrop-blur-sm">
        {spinnerContent}
      </div>
    );
  }

  return spinnerContent;
};

// Simple inline spinner (for buttons, etc.)
export const InlineSpinner: React.FC<{ className?: string }> = ({ className = '' }) => (
  <motion.div
    className={`rounded-full border-2 border-current/25 border-t-current ${className}`}
    style={{ width: '1em', height: '1em' }}
    animate={{ rotate: 360 }}
    transition={{ duration: 0.8, repeat: Infinity, ease: 'linear' }}
  />
);
