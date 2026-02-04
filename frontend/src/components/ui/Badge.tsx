import React from 'react';
import { motion } from 'framer-motion';

export type BadgeVariant = 'success' | 'error' | 'warning' | 'info' | 'default';
export type BadgeSize = 'sm' | 'md' | 'lg';

interface BadgeProps {
  children: React.ReactNode;
  variant?: BadgeVariant;
  size?: BadgeSize;
  animated?: boolean;
  className?: string;
}

const variantClasses: Record<BadgeVariant, string> = {
  success: 'bg-green-500/10 text-green-500 border-green-500/30',
  error: 'bg-red-500/10 text-red-500 border-red-500/30',
  warning: 'bg-amber-500/10 text-amber-500 border-amber-500/30',
  info: 'bg-blue-500/10 text-blue-500 border-blue-500/30',
  default: 'bg-slate-500/10 text-slate-400 border-slate-500/30',
};

const sizeClasses: Record<BadgeSize, string> = {
  sm: 'px-2 py-0.5 text-xs',
  md: 'px-3 py-1 text-sm',
  lg: 'px-4 py-1.5 text-sm',
};

export const Badge: React.FC<BadgeProps> = ({
  children,
  variant = 'default',
  size = 'md',
  animated = false,
  className = '',
}) => {
  const baseClasses = 'inline-flex items-center gap-1.5 font-medium rounded-full border';
  const variantClass = variantClasses[variant];
  const sizeClass = sizeClasses[size];

  return (
    <motion.span
      className={`${baseClasses} ${variantClass} ${sizeClass} ${className}`}
      animate={animated ? { scale: [1, 1.02, 1] } : undefined}
      transition={animated ? { duration: 2, repeat: Infinity, ease: 'easeInOut' } : undefined}
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
