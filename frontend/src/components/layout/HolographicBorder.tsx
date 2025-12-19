import React from 'react';
import { motion } from 'framer-motion';

interface HolographicBorderProps {
  children: React.ReactNode;
  animated?: boolean;
  borderWidth?: number;
  speed?: 'slow' | 'normal' | 'fast';
  rounded?: 'none' | 'sm' | 'md' | 'lg' | 'xl' | '2xl' | 'full';
  className?: string;
}

const speedDurations: Record<string, number> = {
  slow: 4,
  normal: 3,
  fast: 2,
};

const roundedClasses: Record<string, string> = {
  none: 'rounded-none',
  sm: 'rounded-sm',
  md: 'rounded-md',
  lg: 'rounded-lg',
  xl: 'rounded-xl',
  '2xl': 'rounded-2xl',
  full: 'rounded-full',
};

export const HolographicBorder: React.FC<HolographicBorderProps> = ({
  children,
  animated = true,
  borderWidth = 2,
  speed = 'normal',
  rounded = 'xl',
  className = '',
}) => {
  const duration = speedDurations[speed];
  const roundedClass = roundedClasses[rounded];

  return (
    <div className={`relative ${className}`}>
      {/* Holographic border */}
      <motion.div
        className={`absolute inset-0 ${roundedClass} p-[${borderWidth}px]`}
        style={{
          background: 'linear-gradient(90deg, #FF00FF 0%, #00FFFF 25%, #00FF00 50%, #FFFF00 75%, #FF00FF 100%)',
          backgroundSize: animated ? '200% 100%' : '100% 100%',
          WebkitMask: 'linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0)',
          WebkitMaskComposite: 'xor',
          maskComposite: 'exclude',
        }}
        animate={
          animated
            ? {
                backgroundPosition: ['0% 50%', '200% 50%'],
              }
            : {}
        }
        transition={
          animated
            ? {
                duration,
                repeat: Infinity,
                ease: 'linear',
              }
            : {}
        }
      />

      {/* Content */}
      <div className={`relative z-10 ${roundedClass} overflow-hidden`}>
        {children}
      </div>
    </div>
  );
};
