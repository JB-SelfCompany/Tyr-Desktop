// Glassmorphism utilities

export type BlurStrength = 'xs' | 'sm' | 'md' | 'lg' | 'xl';
export type GlassVariant = 'default' | 'strong' | 'subtle';

// Blur values mapping
export const blurValues: Record<BlurStrength, string> = {
  xs: 'blur(4px)',
  sm: 'blur(8px)',
  md: 'blur(12px)',
  lg: 'blur(16px)',
  xl: 'blur(24px)',
};

// Glass styles generator
export const glassStyles = {
  // Light theme glass
  light: (variant: GlassVariant = 'default', blur: BlurStrength = 'lg') => ({
    background: variant === 'strong'
      ? 'rgba(248, 250, 252, 0.98)'
      : variant === 'subtle'
      ? 'rgba(248, 250, 252, 0.85)'
      : 'rgba(248, 250, 252, 0.95)',
    backdropFilter: blurValues[blur],
    WebkitBackdropFilter: blurValues[blur],
    border: '1px solid rgba(203, 213, 225, 0.6)',
    boxShadow: '0 4px 20px rgba(0, 0, 0, 0.08)',
  }),

  // Dark theme glass
  dark: (variant: GlassVariant = 'default', blur: BlurStrength = 'lg') => ({
    background: variant === 'strong'
      ? 'rgba(30, 41, 59, 0.98)'
      : variant === 'subtle'
      ? 'rgba(30, 41, 59, 0.85)'
      : 'rgba(30, 41, 59, 0.95)',
    backdropFilter: blurValues[blur],
    WebkitBackdropFilter: blurValues[blur],
    border: '1px solid rgba(51, 65, 85, 0.6)',
    boxShadow: '0 4px 30px rgba(0, 0, 0, 0.5), 0 0 0 1px rgba(255, 255, 255, 0.05)',
  }),
};

// CSS class generator for Tailwind
export const getGlassClasses = (theme: 'light' | 'dark', variant: GlassVariant = 'default') => {
  const baseClasses = 'backdrop-blur-lg';

  if (theme === 'light') {
    const bgClass = variant === 'strong'
      ? 'bg-slate-50/[0.98]'
      : variant === 'subtle'
      ? 'bg-slate-50/[0.85]'
      : 'bg-slate-50/[0.95]';

    return `${baseClasses} ${bgClass} border border-slate-300/60 shadow-md`;
  } else {
    const bgClass = variant === 'strong'
      ? 'bg-slate-800/[0.98]'
      : variant === 'subtle'
      ? 'bg-slate-800/[0.85]'
      : 'bg-slate-800/[0.95]';

    return `${baseClasses} ${bgClass} border border-slate-700/60 shadow-glass`;
  }
};

// Frosted glass effect
export const frostedGlass = (theme: 'light' | 'dark') => ({
  backdropFilter: 'blur(16px)',
  WebkitBackdropFilter: 'blur(16px)',
  background: theme === 'light'
    ? 'rgba(248, 250, 252, 0.95)'
    : 'rgba(30, 41, 59, 0.95)',
});

// Dynamic glass effect with hover
export const interactiveGlass = {
  base: (theme: 'light' | 'dark') => glassStyles[theme]('default', 'lg'),
  hover: (theme: 'light' | 'dark') => ({
    ...glassStyles[theme]('strong', 'lg'),
    borderColor: theme === 'light' ? '#94a3b8' : '#475569',
  }),
};

// Glass card preset classes
export const glassCard = {
  light: 'backdrop-blur-lg bg-slate-50/[0.95] border border-slate-300/60 rounded-xl shadow-md',
  dark: 'backdrop-blur-lg bg-slate-800/[0.95] border border-slate-700/60 rounded-xl shadow-glass',
};

// Performance-optimized glass (reduced blur)
export const optimizedGlass = {
  light: 'backdrop-blur-sm bg-slate-50/[0.95] border border-slate-300/60 rounded-xl shadow-md',
  dark: 'backdrop-blur-sm bg-slate-800/[0.95] border border-slate-700/60 rounded-xl shadow-glass',
};
