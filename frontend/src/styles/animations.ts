// Framer Motion animation variants

import { Variants, Transition } from 'framer-motion';

// Transition configurations
export const transitions = {
  fast: { duration: 0.15, ease: 'easeOut' as const },
  normal: { duration: 0.2, ease: 'easeOut' as const },
  slow: { duration: 0.3, ease: 'easeOut' as const },
};

// Spring configurations
export const springs = {
  gentle: {
    type: 'spring' as const,
    stiffness: 120,
    damping: 14,
    mass: 0.8,
  },
  snappy: {
    type: 'spring' as const,
    stiffness: 400,
    damping: 25,
    mass: 0.5,
  },
  smooth: {
    type: 'spring' as const,
    stiffness: 80,
    damping: 15,
    mass: 1,
  },
};

// Fade variants
export const fadeVariants: Variants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: transitions.normal,
  },
  exit: {
    opacity: 0,
    transition: transitions.fast,
  },
};

// Slide up variants (for page transitions)
export const slideUpVariants: Variants = {
  hidden: { opacity: 0, y: 10 },
  visible: {
    opacity: 1,
    y: 0,
    transition: transitions.normal,
  },
  exit: {
    opacity: 0,
    transition: transitions.fast,
  },
};

// Slide variants (horizontal)
export const slideVariants: Variants = {
  hidden: { opacity: 0, x: -10 },
  visible: {
    opacity: 1,
    x: 0,
    transition: transitions.normal,
  },
  exit: {
    opacity: 0,
    x: 10,
    transition: transitions.fast,
  },
};

// Scale variants (for modals, popovers)
export const scaleVariants: Variants = {
  hidden: { opacity: 0, scale: 0.95 },
  visible: {
    opacity: 1,
    scale: 1,
    transition: springs.gentle,
  },
  exit: {
    opacity: 0,
    scale: 0.95,
    transition: transitions.fast,
  },
};

// Stagger children animation
export const staggerContainer: Variants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      staggerChildren: 0.05,
      delayChildren: 0.1,
    },
  },
};

export const staggerItem: Variants = {
  hidden: { opacity: 0, y: 10 },
  visible: {
    opacity: 1,
    y: 0,
    transition: transitions.normal,
  },
};

// Card hover animation (subtle)
export const cardHoverVariants: Variants = {
  initial: { scale: 1 },
  hover: {
    borderColor: '#475569', // slate-600
    transition: transitions.fast,
  },
  tap: {
    scale: 0.99,
    transition: transitions.fast,
  },
};

// Button variants
export const buttonVariants: Variants = {
  initial: { scale: 1 },
  hover: {
    scale: 1.02,
    transition: transitions.fast,
  },
  tap: {
    scale: 0.98,
    transition: transitions.fast,
  },
};

// Modal backdrop variants
export const backdropVariants: Variants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: transitions.normal,
  },
  exit: {
    opacity: 0,
    transition: transitions.fast,
  },
};

// Toast notification variants
export const toastVariants: Variants = {
  hidden: {
    opacity: 0,
    x: 50,
  },
  visible: {
    opacity: 1,
    x: 0,
    transition: transitions.normal,
  },
  exit: {
    opacity: 0,
    x: 50,
    transition: transitions.fast,
  },
};

// Page transition variants
export const pageTransitionVariants: Variants = {
  initial: {
    opacity: 0,
    y: 10,
  },
  animate: {
    opacity: 1,
    y: 0,
    transition: transitions.normal,
  },
  exit: {
    opacity: 0,
    transition: transitions.fast,
  },
};

// Rotate animation (for loaders)
export const rotateVariants: Variants = {
  animate: {
    rotate: 360,
    transition: {
      duration: 0.8,
      repeat: Infinity,
      ease: 'linear',
    },
  },
};

// List item variants
export const listItemVariants: Variants = {
  hidden: { opacity: 0, x: -10 },
  visible: {
    opacity: 1,
    x: 0,
    transition: transitions.normal,
  },
  exit: {
    opacity: 0,
    x: 10,
    transition: transitions.fast,
  },
};

// Custom transition creator
export const createCustomTransition = (
  duration: number = 0.2,
  ease: 'easeIn' | 'easeOut' | 'easeInOut' | 'linear' = 'easeOut'
) => ({
  duration,
  ease,
});

// Entrance animation for components
export const entranceVariants: Variants = {
  hidden: {
    opacity: 0,
    y: 10,
  },
  visible: {
    opacity: 1,
    y: 0,
    transition: transitions.normal,
  },
};
