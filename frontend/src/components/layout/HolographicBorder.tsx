import React from 'react';

interface HolographicBorderProps {
  children: React.ReactNode;
  animated?: boolean;
  borderWidth?: number;
  speed?: 'slow' | 'normal' | 'fast';
  rounded?: 'none' | 'sm' | 'md' | 'lg' | 'xl' | '2xl' | 'full';
  className?: string;
}

/**
 * HolographicBorder - Now a simple wrapper without visual effects
 * Kept for backward compatibility
 */
export const HolographicBorder: React.FC<HolographicBorderProps> = ({
  children,
  className = '',
}) => {
  return (
    <div className={className}>
      {children}
    </div>
  );
};
