import React, { forwardRef, InputHTMLAttributes } from 'react';
import { motion } from 'framer-motion';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  helperText?: string;
  leftIcon?: React.ReactNode;
  rightIcon?: React.ReactNode;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  (
    {
      label,
      error,
      helperText,
      leftIcon,
      rightIcon,
      className = '',
      ...props
    },
    ref
  ) => {
    const inputId = props.id || `input-${Math.random().toString(36).substr(2, 9)}`;
    const errorId = `${inputId}-error`;
    const helperId = `${inputId}-helper`;

    const baseInputClasses = 'w-full px-4 py-2.5 rounded-xl text-sm transition-all duration-200 outline-none';

    const themeClasses = 'bg-slate-800 border text-slate-100 placeholder-slate-500';

    const borderClasses = error
      ? 'border-red-500 focus:border-red-500 focus:ring-2 focus:ring-red-500/50'
      : 'border-slate-600 focus:border-emerald-500 focus:ring-2 focus:ring-emerald-500/50';

    const paddingClasses = leftIcon ? 'pl-10' : rightIcon ? 'pr-10' : '';

    return (
      <div className="w-full">
        {label && (
          <label
            htmlFor={inputId}
            className="block mb-2 text-sm font-medium text-slate-200"
          >
            {label}
          </label>
        )}

        <div className="relative">
          {leftIcon && (
            <div className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400">
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
            <div className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400">
              {rightIcon}
            </div>
          )}
        </div>

        {error && (
          <motion.p
            id={errorId}
            initial={{ opacity: 0, y: -5 }}
            animate={{ opacity: 1, y: 0 }}
            className="mt-1.5 text-sm text-red-500"
            role="alert"
            aria-live="polite"
          >
            {error}
          </motion.p>
        )}

        {helperText && !error && (
          <p id={helperId} className="mt-1.5 text-sm text-slate-400">
            {helperText}
          </p>
        )}
      </div>
    );
  }
);

Input.displayName = 'Input';
