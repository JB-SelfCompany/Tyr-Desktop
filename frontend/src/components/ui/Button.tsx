import React from 'react';
import { motion, HTMLMotionProps } from 'framer-motion';
import { buttonVariants } from '../../styles/animations';

export type ButtonVariant = 'primary' | 'secondary' | 'danger' | 'ghost' | 'neon';
export type ButtonSize = 'sm' | 'md' | 'lg';

interface ButtonProps extends Omit<HTMLMotionProps<'button'>, 'size'> {
  variant?: ButtonVariant;
  size?: ButtonSize;
  glow?: boolean;
  loading?: boolean;
  fullWidth?: boolean;
  children: React.ReactNode;
}

const variantClasses: Record<ButtonVariant, string> = {
  primary: 'bg-md-light-primary hover:bg-[#005138] dark:bg-md-dark-primary dark:hover:bg-[#89F8C7] text-md-light-onPrimary dark:text-md-dark-onPrimary border-md-light-primary dark:border-md-dark-primary',
  secondary: 'bg-md-light-secondaryContainer hover:bg-[#B8D6C9] dark:bg-md-dark-secondaryContainer dark:hover:bg-[#405B50] text-md-light-onSecondaryContainer dark:text-md-dark-onSecondaryContainer border-md-light-secondary/30 dark:border-md-dark-secondary/30',
  danger: 'bg-md-light-error hover:bg-[#930012] dark:bg-md-dark-error dark:hover:bg-[#FFD9D4] text-md-light-onError dark:text-md-dark-onError border-md-light-error dark:border-md-dark-error',
  ghost: 'bg-transparent hover:bg-md-light-primaryContainer/50 dark:hover:bg-md-dark-primaryContainer/20 text-md-light-onSurface dark:text-md-dark-onSurface border-md-light-outline/30 dark:border-md-dark-outline/30',
  neon: 'bg-gradient-to-r from-neon-pink to-neon-purple text-white border-neon-pink shadow-neon-pink',
};

const sizeClasses: Record<ButtonSize, string> = {
  sm: 'px-3 py-1.5 text-sm',
  md: 'px-4 py-2 text-base',
  lg: 'px-6 py-3 text-lg',
};

const glowClasses: Record<ButtonVariant, string> = {
  primary: 'shadow-[0_0_20px_rgba(0,108,76,0.4)] dark:shadow-[0_0_20px_rgba(108,219,156,0.4)]',
  secondary: 'shadow-[0_0_15px_rgba(77,99,87,0.3)] dark:shadow-[0_0_15px_rgba(179,204,190,0.3)]',
  danger: 'shadow-[0_0_20px_rgba(186,26,26,0.4)] dark:shadow-[0_0_20px_rgba(255,180,171,0.4)]',
  ghost: '',
  neon: 'shadow-neon-pink animate-glow-pulse',
};

export const Button: React.FC<ButtonProps> = ({
  variant = 'primary',
  size = 'md',
  glow = false,
  loading = false,
  fullWidth = false,
  disabled,
  children,
  className = '',
  ...props
}) => {
  const baseClasses = 'font-futuristic font-semibold tracking-wide rounded-lg border transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed';
  const variantClass = variantClasses[variant];
  const sizeClass = sizeClasses[size];
  const glowClass = glow ? glowClasses[variant] : '';
  const widthClass = fullWidth ? 'w-full' : '';

  return (
    <motion.button
      className={`${baseClasses} ${variantClass} ${sizeClass} ${glowClass} ${widthClass} ${className}`}
      variants={buttonVariants}
      initial="initial"
      whileHover={!disabled && !loading ? "hover" : undefined}
      whileTap={!disabled && !loading ? "tap" : undefined}
      disabled={disabled || loading}
      aria-busy={loading}
      aria-disabled={disabled || loading}
      role="button"
      {...props}
    >
      {loading ? (
        <span className="flex items-center justify-center gap-2">
          <svg
            className="animate-spin h-5 w-5"
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
          <span>Loading...</span>
        </span>
      ) : (
        children
      )}
    </motion.button>
  );
};
