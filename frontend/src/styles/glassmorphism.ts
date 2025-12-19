// Glassmorphism utilities for Y2K Futurism design

export type BlurStrength = 'xs' | 'sm' | 'md' | 'lg' | 'xl';
export type GlassVariant = 'default' | 'strong' | 'subtle';

// Blur values mapping
export const blurValues: Record<BlurStrength, string> = {
  xs: 'blur(4px)',
  sm: 'blur(8px)',
  md: 'blur(12px)',
  lg: 'blur(20px)',
  xl: 'blur(40px)',
};

// Glass styles generator
export const glassStyles = {
  // Light theme glass
  light: (variant: GlassVariant = 'default', blur: BlurStrength = 'lg') => ({
    background: variant === 'strong'
      ? 'rgba(255, 255, 255, 0.85)'
      : variant === 'subtle'
      ? 'rgba(255, 255, 255, 0.5)'
      : 'rgba(255, 255, 255, 0.7)',
    backdropFilter: blurValues[blur],
    WebkitBackdropFilter: blurValues[blur],
    border: '1px solid rgba(255, 255, 255, 0.3)',
    boxShadow: '0 8px 32px 0 rgba(31, 38, 135, 0.37)',
  }),

  // Dark theme glass
  dark: (variant: GlassVariant = 'default', blur: BlurStrength = 'lg') => ({
    background: variant === 'strong'
      ? 'rgba(26, 31, 54, 0.85)'
      : variant === 'subtle'
      ? 'rgba(26, 31, 54, 0.5)'
      : 'rgba(26, 31, 54, 0.7)',
    backdropFilter: blurValues[blur],
    WebkitBackdropFilter: blurValues[blur],
    border: '1px solid rgba(255, 255, 255, 0.1)',
    boxShadow: '0 8px 32px 0 rgba(0, 0, 0, 0.5)',
  }),
};

// CSS class generator for Tailwind
export const getGlassClasses = (theme: 'light' | 'dark', variant: GlassVariant = 'default') => {
  const baseClasses = 'backdrop-blur-lg';

  if (theme === 'light') {
    const bgClass = variant === 'strong'
      ? 'bg-white/85'
      : variant === 'subtle'
      ? 'bg-white/50'
      : 'bg-white/70';

    return `${baseClasses} ${bgClass} border border-white/30 shadow-glass`;
  } else {
    const bgClass = variant === 'strong'
      ? 'bg-space-blue-light/85'
      : variant === 'subtle'
      ? 'bg-space-blue-light/50'
      : 'bg-space-blue-light/70';

    return `${baseClasses} ${bgClass} border border-white/10 shadow-xl`;
  }
};

// Frosted glass effect with saturation boost
export const frostedGlass = (theme: 'light' | 'dark') => ({
  backdropFilter: 'blur(20px) saturate(180%)',
  WebkitBackdropFilter: 'blur(20px) saturate(180%)',
  background: theme === 'light'
    ? 'rgba(255, 255, 255, 0.7)'
    : 'rgba(26, 31, 54, 0.7)',
});

// Dynamic glass effect with hover
export const interactiveGlass = {
  base: (theme: 'light' | 'dark') => glassStyles[theme]('default', 'lg'),
  hover: (theme: 'light' | 'dark') => ({
    ...glassStyles[theme]('strong', 'xl'),
    transform: 'translateY(-2px)',
    boxShadow: theme === 'light'
      ? '0 12px 40px 0 rgba(31, 38, 135, 0.45)'
      : '0 12px 40px 0 rgba(0, 0, 0, 0.6)',
  }),
};

// Glass card preset
export const glassCard = {
  light: 'backdrop-blur-lg bg-white/70 border border-white/30 rounded-2xl shadow-glass',
  dark: 'backdrop-blur-lg bg-space-blue-light/70 border border-white/10 rounded-2xl shadow-xl',
};

// Performance-optimized glass (reduced blur for weaker GPUs)
export const optimizedGlass = {
  light: 'backdrop-blur-sm bg-white/80 border border-white/30 rounded-2xl shadow-glass',
  dark: 'backdrop-blur-sm bg-space-blue-light/80 border border-white/10 rounded-2xl shadow-xl',
};
