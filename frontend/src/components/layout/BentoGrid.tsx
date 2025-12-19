import React from 'react';

interface BentoGridProps {
  children: React.ReactNode;
  columns?: 2 | 3 | 4;
  gap?: 'sm' | 'md' | 'lg';
  className?: string;
}

const columnClasses: Record<number, string> = {
  2: 'grid-cols-1 md:grid-cols-2',
  3: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3',
  4: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-4',
};

const gapClasses: Record<string, string> = {
  sm: 'gap-3',
  md: 'gap-4',
  lg: 'gap-6',
};

export const BentoGrid: React.FC<BentoGridProps> = ({
  children,
  columns = 3,
  gap = 'md',
  className = '',
}) => {
  const columnClass = columnClasses[columns];
  const gapClass = gapClasses[gap];

  return (
    <div className={`grid ${columnClass} ${gapClass} ${className}`}>
      {children}
    </div>
  );
};

interface BentoCardProps {
  children: React.ReactNode;
  span?: 1 | 2 | 3 | 4;
  rowSpan?: 1 | 2 | 3;
  className?: string;
}

const spanClasses: Record<number, string> = {
  1: 'col-span-1',
  2: 'col-span-1 md:col-span-2',
  3: 'col-span-1 md:col-span-2 lg:col-span-3',
  4: 'col-span-1 md:col-span-2 lg:col-span-4',
};

const rowSpanClasses: Record<number, string> = {
  1: 'row-span-1',
  2: 'row-span-2',
  3: 'row-span-3',
};

export const BentoCard: React.FC<BentoCardProps> = ({
  children,
  span = 1,
  rowSpan = 1,
  className = '',
}) => {
  const spanClass = spanClasses[span];
  const rowSpanClass = rowSpanClasses[rowSpan];

  return (
    <div className={`${spanClass} ${rowSpanClass} ${className}`}>
      {children}
    </div>
  );
};
