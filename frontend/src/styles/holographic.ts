// Holographic and iridescent effects for Y2K Futurism

export type HolographicSpeed = 'slow' | 'normal' | 'fast';

// Holographic gradient definitions
export const holographicGradients = {
  // Rainbow holographic
  rainbow: 'linear-gradient(135deg, #FF00FF 0%, #00FFFF 25%, #00FF00 50%, #FFFF00 75%, #FF00FF 100%)',

  // Border holographic (horizontal)
  border: 'linear-gradient(90deg, #FF00FF 0%, #00FFFF 25%, #00FF00 50%, #FFFF00 75%, #FF00FF 100%)',

  // Iridescent (softer colors)
  iridescent: 'linear-gradient(90deg, #667eea 0%, #764ba2 25%, #f093fb 50%, #4facfe 75%, #00f2fe 100%)',

  // Cosmic (deep space colors)
  cosmic: 'linear-gradient(135deg, #667eea 0%, #764ba2 33%, #f093fb 66%, #667eea 100%)',

  // Neon shimmer
  neonShimmer: 'linear-gradient(90deg, #FF00FF 0%, #00FFFF 50%, #FF00FF 100%)',
};

// Animation durations
const speedDurations: Record<HolographicSpeed, string> = {
  slow: '4s',
  normal: '3s',
  fast: '2s',
};

// Animated holographic border
export const holographicBorder = (speed: HolographicSpeed = 'normal', width: number = 2) => ({
  position: 'relative' as const,
  '&::before': {
    content: '""',
    position: 'absolute' as const,
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    borderRadius: 'inherit',
    padding: `${width}px`,
    background: holographicGradients.border,
    WebkitMask: 'linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0)',
    WebkitMaskComposite: 'xor',
    maskComposite: 'exclude',
    backgroundSize: '200% 100%',
    animation: `holographic-slide ${speedDurations[speed]} linear infinite`,
  },
});

// Holographic text effect (CSS-in-JS)
export const holographicTextEffect = {
  background: holographicGradients.iridescent,
  WebkitBackgroundClip: 'text',
  WebkitTextFillColor: 'transparent',
  backgroundClip: 'text',
  backgroundSize: '200% 100%',
  animation: 'holographic-slide 3s linear infinite',
};

// Iridescent background with animation
export const iridescentBackground = (speed: HolographicSpeed = 'normal') => ({
  background: holographicGradients.iridescent,
  backgroundSize: '200% 100%',
  animation: `holographic-slide ${speedDurations[speed]} linear infinite`,
});

// Shimmer effect (for loading states)
export const shimmerEffect = {
  background: 'linear-gradient(90deg, transparent 0%, rgba(255, 255, 255, 0.4) 50%, transparent 100%)',
  backgroundSize: '200% 100%',
  animation: 'shimmer 2s linear infinite',
};

// CSS classes for Tailwind usage
export const holographicClasses = {
  border: 'relative before:absolute before:inset-0 before:rounded-[inherit] before:p-[2px] before:bg-holographic-border before:[-webkit-mask:linear-gradient(#fff_0_0)_content-box,linear-gradient(#fff_0_0)] before:[mask-composite:exclude] before:animate-holographic-spin',

  text: 'bg-iridescent bg-clip-text text-transparent bg-[length:200%_100%] animate-holographic-spin',

  background: 'bg-iridescent bg-[length:200%_100%] animate-holographic-spin',

  shimmer: 'bg-gradient-to-r from-transparent via-white/40 to-transparent bg-[length:200%_100%] animate-shimmer',
};

// Neon glow effect
export type NeonGlowColor = 'pink' | 'cyan' | 'green' | 'yellow' | 'purple' | 'blue';

const neonColorValues: Record<NeonGlowColor, string> = {
  pink: '#FF00FF',
  cyan: '#00FFFF',
  green: '#00FF88',
  yellow: '#FFFF00',
  purple: '#9D00FF',
  blue: '#0080FF',
};

export const neonGlow = (color: NeonGlowColor, intensity: 'low' | 'medium' | 'high' = 'medium') => {
  const colorValue = neonColorValues[color];
  const sizes = {
    low: '0 0 10px',
    medium: '0 0 20px',
    high: '0 0 40px',
  };

  return {
    boxShadow: `${sizes[intensity]} ${colorValue}`,
    animation: 'glow-pulse 2s ease-in-out infinite',
  };
};

// Pulsing neon glow
export const pulsingGlow = (color: NeonGlowColor) => {
  const colorValue = neonColorValues[color];
  return {
    animation: 'glow-pulse 2s ease-in-out infinite',
    '--glow-color': colorValue,
  } as React.CSSProperties;
};

// Holographic card effect
export const holographicCard = {
  base: 'relative overflow-hidden',
  shine: 'absolute inset-0 bg-gradient-to-br from-white/0 via-white/20 to-white/0 opacity-0 group-hover:opacity-100 transition-opacity duration-500',
};

// Rainbow border animation keyframes (for injection into global CSS)
export const holographicKeyframes = `
@keyframes holographic-slide {
  0% {
    background-position: 0% 50%;
  }
  100% {
    background-position: 200% 50%;
  }
}

@keyframes shimmer {
  0% {
    background-position: -200% 0;
  }
  100% {
    background-position: 200% 0;
  }
}

@keyframes glow-pulse {
  0%, 100% {
    filter: drop-shadow(0 0 8px var(--glow-color, currentColor));
  }
  50% {
    filter: drop-shadow(0 0 16px var(--glow-color, currentColor));
  }
}
`;
