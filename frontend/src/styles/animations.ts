// Framer Motion animation variants for Y2K Futurism

import { Variants, Transition } from 'framer-motion';

// Spring configurations (no linear!)
export const springs = {
  // Gentle spring (for subtle animations)
  gentle: {
    type: 'spring' as const,
    stiffness: 120,
    damping: 14,
    mass: 0.8,
  },

  // Bouncy spring (for playful interactions)
  bouncy: {
    type: 'spring' as const,
    stiffness: 300,
    damping: 10,
    mass: 0.5,
  },

  // Smooth spring (for large elements)
  smooth: {
    type: 'spring' as const,
    stiffness: 80,
    damping: 15,
    mass: 1,
  },

  // Snappy spring (for quick interactions)
  snappy: {
    type: 'spring' as const,
    stiffness: 400,
    damping: 25,
    mass: 0.5,
  },

  // Wobbly spring (for attention-grabbing)
  wobbly: {
    type: 'spring' as const,
    stiffness: 180,
    damping: 8,
    mass: 0.8,
  },
};

// Fade variants
export const fadeVariants: Variants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: springs.gentle,
  },
  exit: {
    opacity: 0,
    transition: springs.smooth,
  },
};

// Slide variants
export const slideVariants: Variants = {
  hidden: { opacity: 0, x: -50 },
  visible: {
    opacity: 1,
    x: 0,
    transition: springs.bouncy,
  },
  exit: {
    opacity: 0,
    x: 50,
    transition: springs.smooth,
  },
};

// Scale variants (for modals, popovers)
export const scaleVariants: Variants = {
  hidden: { opacity: 0, scale: 0.8 },
  visible: {
    opacity: 1,
    scale: 1,
    transition: springs.bouncy,
  },
  exit: {
    opacity: 0,
    scale: 0.8,
    transition: springs.smooth,
  },
};

// Float animation (for floating elements)
export const floatVariants: Variants = {
  animate: {
    y: [-10, 10, -10],
    transition: {
      duration: 3,
      repeat: Infinity,
      ease: 'easeInOut',
    },
  },
};

// Glow pulse animation
export const glowPulseVariants: Variants = {
  animate: {
    filter: [
      'drop-shadow(0 0 8px currentColor)',
      'drop-shadow(0 0 16px currentColor)',
      'drop-shadow(0 0 8px currentColor)',
    ],
    transition: {
      duration: 2,
      repeat: Infinity,
      ease: 'easeInOut',
    },
  },
};

// Slide up variants (for page transitions)
export const slideUpVariants: Variants = {
  hidden: { opacity: 0, y: 50 },
  visible: {
    opacity: 1,
    y: 0,
    transition: springs.gentle,
  },
  exit: {
    opacity: 0,
    y: -50,
    transition: springs.smooth,
  },
};

// Stagger children animation
export const staggerContainer: Variants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      staggerChildren: 0.1,
      delayChildren: 0.2,
    },
  },
};

export const staggerItem: Variants = {
  hidden: { opacity: 0, y: 20 },
  visible: {
    opacity: 1,
    y: 0,
    transition: springs.bouncy,
  },
};

// Card hover animation
export const cardHoverVariants: Variants = {
  initial: { scale: 1 },
  hover: {
    scale: 1.02,
    y: -4,
    transition: springs.snappy,
  },
  tap: {
    scale: 0.98,
    transition: springs.snappy,
  },
};

// Button variants
export const buttonVariants: Variants = {
  initial: { scale: 1 },
  hover: {
    scale: 1.05,
    transition: springs.bouncy,
  },
  tap: {
    scale: 0.95,
    transition: springs.snappy,
  },
};

// Modal backdrop variants
export const backdropVariants: Variants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: { duration: 0.3 },
  },
  exit: {
    opacity: 0,
    transition: { duration: 0.2 },
  },
};

// Shimmer animation (for loading)
export const shimmerVariants: Variants = {
  animate: {
    backgroundPosition: ['0% 50%', '100% 50%', '0% 50%'],
    transition: {
      duration: 2,
      repeat: Infinity,
      ease: 'linear',
    },
  },
};

// Rotate animation (for loaders)
export const rotateVariants: Variants = {
  animate: {
    rotate: 360,
    transition: {
      duration: 1,
      repeat: Infinity,
      ease: 'linear',
    },
  },
};

// Page transition variants
export const pageTransitionVariants: Variants = {
  initial: {
    opacity: 0,
    scale: 0.95,
  },
  animate: {
    opacity: 1,
    scale: 1,
    transition: springs.gentle,
  },
  exit: {
    opacity: 0,
    scale: 1.05,
    transition: springs.smooth,
  },
};

// Neon glow animation
export const neonGlowVariants: Variants = {
  animate: {
    boxShadow: [
      '0 0 10px currentColor',
      '0 0 20px currentColor',
      '0 0 30px currentColor',
      '0 0 20px currentColor',
      '0 0 10px currentColor',
    ],
    transition: {
      duration: 2,
      repeat: Infinity,
      ease: 'easeInOut',
    },
  },
};

// Holographic border animation
export const holographicBorderVariants: Variants = {
  animate: {
    backgroundPosition: ['0% 50%', '200% 50%'],
    transition: {
      duration: 3,
      repeat: Infinity,
      ease: 'linear',
    },
  },
};

// Toast notification variants
export const toastVariants: Variants = {
  hidden: {
    opacity: 0,
    x: 100,
    scale: 0.8,
  },
  visible: {
    opacity: 1,
    x: 0,
    scale: 1,
    transition: springs.bouncy,
  },
  exit: {
    opacity: 0,
    x: 100,
    scale: 0.8,
    transition: springs.smooth,
  },
};

// List item variants (for animated lists)
export const listItemVariants: Variants = {
  hidden: { opacity: 0, x: -20 },
  visible: {
    opacity: 1,
    x: 0,
    transition: springs.gentle,
  },
  exit: {
    opacity: 0,
    x: 20,
    transition: springs.smooth,
  },
};

// Custom spring transition for specific use cases
export const createCustomSpring = (
  stiffness: number = 120,
  damping: number = 14,
  mass: number = 0.8
): Transition => ({
  type: 'spring',
  stiffness,
  damping,
  mass,
});

// Entrance animation for components
export const entranceVariants: Variants = {
  hidden: {
    opacity: 0,
    y: 30,
    scale: 0.9,
  },
  visible: {
    opacity: 1,
    y: 0,
    scale: 1,
    transition: {
      ...springs.bouncy,
      delay: 0.1,
    },
  },
};

// Hover lift animation (for interactive cards)
export const hoverLiftVariants: Variants = {
  initial: {
    y: 0,
    boxShadow: '0 4px 12px rgba(0, 0, 0, 0.1)',
  },
  hover: {
    y: -8,
    boxShadow: '0 12px 24px rgba(0, 0, 0, 0.2)',
    transition: springs.snappy,
  },
};
