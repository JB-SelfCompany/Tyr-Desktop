import React from 'react';
import { motion } from 'framer-motion';

export type BadgeVariant = 'success' | 'error' | 'warning' | 'info' | 'default';
export type BadgeSize = 'sm' | 'md' | 'lg';

interface BadgeProps {
  children: React.ReactNode;
  variant?: BadgeVariant;
  size?: BadgeSize;
  animated?: boolean;
  glow?: boolean;
  className?: string;
}

const variantClasses: Record<BadgeVariant, string> = {
  success: 'bg-md-light-primaryContainer dark:bg-md-dark-primaryContainer text-md-light-onPrimaryContainer dark:text-md-dark-onPrimaryContainer border-md-light-primary/30 dark:border-md-dark-primary/30',
  error: 'bg-md-light-errorContainer dark:bg-md-dark-errorContainer text-md-light-onErrorContainer dark:text-md-dark-onErrorContainer border-md-light-error/30 dark:border-md-dark-error/30',
  warning: 'bg-[#FFF4E0] dark:bg-[#4D3800] text-[#663C00] dark:text-[#FFD89C] border-[#FFB74D]/30 dark:border-[#FFB74D]/30',
  info: 'bg-md-light-tertiaryContainer dark:bg-md-dark-tertiaryContainer text-md-light-onTertiaryContainer dark:text-md-dark-onTertiaryContainer border-md-light-tertiary/30 dark:border-md-dark-tertiary/30',
  default: 'bg-md-light-surfaceVariant dark:bg-md-dark-surfaceVariant text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant border-md-light-outline/30 dark:border-md-dark-outline/30',
};

const glowClasses: Record<BadgeVariant, string> = {
  success: 'shadow-[0_0_15px_rgba(0,108,76,0.5)] dark:shadow-[0_0_15px_rgba(108,219,156,0.5)]',
  error: 'shadow-[0_0_15px_rgba(186,26,26,0.5)] dark:shadow-[0_0_15px_rgba(255,180,171,0.5)]',
  warning: 'shadow-[0_0_15px_rgba(255,152,0,0.5)]',
  info: 'shadow-[0_0_15px_rgba(61,99,115,0.5)] dark:shadow-[0_0_15px_rgba(165,204,222,0.5)]',
  default: '',
};

const sizeClasses: Record<BadgeSize, string> = {
  sm: 'px-2 py-0.5 text-xs',
  md: 'px-3 py-1 text-sm',
  lg: 'px-4 py-1.5 text-base',
};

export const Badge: React.FC<BadgeProps> = ({
  children,
  variant = 'default',
  size = 'md',
  animated = false,
  glow = false,
  className = '',
}) => {
  const baseClasses = 'inline-flex items-center gap-1.5 font-futuristic font-semibold rounded-full border backdrop-blur-lg';
  const variantClass = variantClasses[variant];
  const sizeClass = sizeClasses[size];
  const glowClass = glow ? glowClasses[variant] : '';

  const pulseAnimation = animated
    ? {
        scale: [1, 1.05, 1],
        opacity: [1, 0.9, 1],
      }
    : {};

  const pulseTransition = animated
    ? {
        duration: 2,
        repeat: Infinity,
        ease: 'easeInOut' as const,
      }
    : undefined;

  return (
    <motion.span
      className={`${baseClasses} ${variantClass} ${sizeClass} ${glowClass} ${className}`}
      animate={pulseAnimation}
      transition={pulseTransition}
    >
      {animated && (
        <span className="relative flex h-2 w-2">
          <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-current opacity-75" />
          <span className="relative inline-flex rounded-full h-2 w-2 bg-current" />
        </span>
      )}
      {children}
    </motion.span>
  );
};
