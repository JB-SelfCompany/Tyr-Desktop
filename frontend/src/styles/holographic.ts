// Subtle visual effects

// Shimmer effect for loading states
export const shimmerEffect = {
  background: 'linear-gradient(90deg, transparent 0%, rgba(255, 255, 255, 0.1) 50%, transparent 100%)',
  backgroundSize: '200% 100%',
  animation: 'shimmer 2s linear infinite',
};

// Glow effect for accent elements
export const accentGlow = (color: string = '#10b981') => ({
  boxShadow: `0 0 20px ${color}25`,
});

// Focus glow
export const focusGlow = {
  boxShadow: '0 0 0 3px rgba(16, 185, 129, 0.2)',
};

// CSS classes for Tailwind usage
export const effectClasses = {
  shimmer: 'bg-gradient-to-r from-transparent via-white/10 to-transparent bg-[length:200%_100%] animate-shimmer',
  glow: 'shadow-glow',
  focusRing: 'focus:ring-2 focus:ring-emerald-500 focus:ring-offset-2 focus:ring-offset-slate-900',
};

// Card shine effect on hover
export const cardShine = {
  base: 'relative overflow-hidden',
  shine: 'absolute inset-0 bg-gradient-to-br from-white/0 via-white/5 to-white/0 opacity-0 group-hover:opacity-100 transition-opacity duration-300',
};

// Keyframes for shimmer animation
export const shimmerKeyframes = `
@keyframes shimmer {
  0% {
    background-position: -200% 0;
  }
  100% {
    background-position: 200% 0;
  }
}
`;
