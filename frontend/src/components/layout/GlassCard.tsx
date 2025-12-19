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
  md: 'p-4 md:p-6',
  lg: 'p-5 md:p-8',
  xl: 'p-6 md:p-10',
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
    default: 'bg-md-light-surface/92 dark:bg-[#1F2B25]/93',
    strong: 'bg-md-light-surface dark:bg-[#1F2B25]',
    subtle: 'bg-md-light-surface/88 dark:bg-[#1F2B25]/88',
  };

  const baseClasses = 'backdrop-blur-lg border border-[#C7CDC7] dark:border-[#404943] rounded-2xl shadow-lg relative overflow-hidden';
  const paddingClass = paddingClasses[padding];

  return (
    <motion.div
      className={`${baseClasses} ${variantClasses[variant]} ${paddingClass} ${className}`}
      variants={hoverable ? cardHoverVariants : undefined}
      initial={hoverable ? 'initial' : undefined}
      whileHover={hoverable ? 'hover' : undefined}
      whileTap={hoverable ? 'tap' : undefined}
      {...props}
    >
      {/* Accent gradient (optional) */}
      {accentColor && (
        <div
          className="absolute top-0 left-0 right-0 h-1 rounded-t-2xl"
          style={{ background: accentColor }}
        />
      )}

      {/* Header */}
      {(title || subtitle || headerAction) && (
        <div className="mb-4 flex items-start justify-between">
          <div>
            {title && (
              <h3 className="text-xl font-display font-bold text-md-light-onSurface dark:text-md-dark-onSurface">
                {title}
              </h3>
            )}
            {subtitle && (
              <p className="mt-1 text-sm text-md-light-onSurfaceVariant dark:text-md-dark-onSurfaceVariant font-body">
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

      {/* Shine effect (for hoverable cards) */}
      {hoverable && (
        <div className="absolute inset-0 bg-gradient-to-br from-white/0 via-white/10 to-white/0 opacity-0 group-hover:opacity-100 transition-opacity duration-500 pointer-events-none" />
      )}
    </motion.div>
  );
};
