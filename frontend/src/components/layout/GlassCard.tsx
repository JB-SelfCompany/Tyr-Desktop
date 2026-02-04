import React from 'react';
import { motion, HTMLMotionProps } from 'framer-motion';
import { cardHoverVariants } from '../../styles/animations';

interface GlassCardProps extends Omit<HTMLMotionProps<'div'>, 'title'> {
  title?: string;
  subtitle?: string;
  accentColor?: string;
  variant?: 'default' | 'strong' | 'subtle';
  hoverable?: boolean;
  padding?: 'none' | 'sm' | 'md' | 'lg' | 'xl';
  children: React.ReactNode;
  headerAction?: React.ReactNode;
}

const paddingClasses: Record<string, string> = {
  none: '',
  sm: 'p-3 md:p-4',
  md: 'p-4 md:p-5',
  lg: 'p-5 md:p-6',
  xl: 'p-6 md:p-8',
};

export const GlassCard: React.FC<GlassCardProps> = ({
  title,
  subtitle,
  accentColor,
  variant = 'default',
  hoverable = false,
  padding = 'md',
  children,
  headerAction,
  className = '',
  ...props
}) => {
  const variantClasses = {
    default: 'glass',
    strong: 'bg-slate-800',
    subtle: 'bg-slate-800/90',
  };

  const baseClasses = 'backdrop-blur-lg border border-slate-700/50 rounded-2xl relative overflow-hidden shadow-glass';
  const paddingClass = paddingClasses[padding];
  const hoverClasses = hoverable ? 'transition-colors hover:border-slate-600' : '';

  return (
    <motion.div
      className={`${baseClasses} ${variantClasses[variant]} ${paddingClass} ${hoverClasses} ${className}`}
      variants={hoverable ? cardHoverVariants : undefined}
      initial={hoverable ? 'initial' : undefined}
      whileHover={hoverable ? 'hover' : undefined}
      whileTap={hoverable ? 'tap' : undefined}
      {...props}
    >
      {/* Accent bar (optional) */}
      {accentColor && (
        <div
          className="absolute top-0 left-0 right-0 h-1 rounded-t-xl"
          style={{ background: accentColor }}
        />
      )}

      {/* Header */}
      {(title || subtitle || headerAction) && (
        <div className="mb-4 flex items-start justify-between">
          <div>
            {title && (
              <h3 className="text-lg font-semibold text-slate-100">
                {title}
              </h3>
            )}
            {subtitle && (
              <p className="mt-1 text-sm text-slate-400">
                {subtitle}
              </p>
            )}
          </div>
          {headerAction && (
            <div className="ml-4 flex-shrink-0">
              {headerAction}
            </div>
          )}
        </div>
      )}

      {/* Content */}
      <div className="relative z-10">
        {children}
      </div>
    </motion.div>
  );
};
