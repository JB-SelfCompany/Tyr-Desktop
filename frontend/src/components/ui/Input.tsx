import React, { forwardRef, InputHTMLAttributes } from 'react';
import { motion, AnimatePresence } from 'framer-motion';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  helperText?: string;
  holographic?: boolean;
  leftIcon?: React.ReactNode;
  rightIcon?: React.ReactNode;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  (
    {
      label,
      error,
      helperText,
      holographic = false,
      leftIcon,
      rightIcon,
      className = '',
      ...props
    },
    ref
  ) => {
    // Generate unique IDs for accessibility
    const inputId = props.id || `input-${Math.random().toString(36).substr(2, 9)}`;
    const errorId = `${inputId}-error`;
    const helperId = `${inputId}-helper`;

    const baseInputClasses = 'w-full px-4 py-2.5 rounded-lg font-body text-base transition-all duration-200 outline-none';

    const themeClasses = 'bg-md-light-surfaceVariant dark:bg-md-dark-surfaceVariant backdrop-blur-lg border text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant placeholder-md-light-outline dark:placeholder-md-dark-outline';

    const borderClasses = error
      ? 'border-md-light-error dark:border-md-dark-error focus:border-md-light-error dark:focus:border-md-dark-error'
      : holographic
      ? 'border-transparent bg-gradient-to-r from-neon-pink/20 via-neon-cyan/20 to-neon-green/20 focus:from-neon-pink/40 focus:via-neon-cyan/40 focus:to-neon-green/40'
      : 'border-md-light-outline/50 dark:border-md-dark-outline/50 focus:border-md-light-primary dark:focus:border-md-dark-primary';

    const paddingClasses = leftIcon ? 'pl-10' : rightIcon ? 'pr-10' : '';

    return (
      <div className="w-full">
        {label && (
          <label htmlFor={inputId} className="block mb-2 text-sm font-futuristic font-semibold text-md-light-onSurface dark:text-md-dark-onSurface">
            {label}
          </label>
        )}

        <div className="relative">
          {leftIcon && (
            <div className="absolute left-3 top-1/2 -translate-y-1/2 text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant">
              {leftIcon}
            </div>
          )}

          <input
            ref={ref}
            id={inputId}
            className={`${baseInputClasses} ${themeClasses} ${borderClasses} ${paddingClasses} ${className}`}
            aria-invalid={!!error}
            aria-describedby={
              error ? errorId : helperText ? helperId : undefined
            }
            {...props}
          />

          {rightIcon && (
            <div className="absolute right-3 top-1/2 -translate-y-1/2 text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant">
              {rightIcon}
            </div>
          )}
        </div>

        {error && (
          <motion.p
            id={errorId}
            initial={{ opacity: 0, y: -5 }}
            animate={{ opacity: 1, y: 0 }}
            className="mt-1.5 text-sm text-md-light-error dark:text-md-dark-error font-body"
            role="alert"
            aria-live="polite"
          >
            {error}
          </motion.p>
        )}

        {helperText && !error && (
          <p id={helperId} className="mt-1.5 text-sm text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant font-body">
            {helperText}
          </p>
        )}
      </div>
    );
  }
);

Input.displayName = 'Input';
