// Typography system

// Font families
export const fonts = {
  // Sans-serif for all text
  sans: "'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', sans-serif",

  // Monospace for code/technical data
  mono: "'JetBrains Mono', 'SF Mono', 'Fira Code', 'Consolas', monospace",
};

// Font sizes
export const fontSizes = {
  xs: '0.75rem',     // 12px
  sm: '0.875rem',    // 14px
  base: '1rem',      // 16px
  lg: '1.125rem',    // 18px
  xl: '1.25rem',     // 20px
  '2xl': '1.5rem',   // 24px
  '3xl': '1.875rem', // 30px
  '4xl': '2.25rem',  // 36px
};

// Font weights
export const fontWeights = {
  normal: 400,
  medium: 500,
  semibold: 600,
  bold: 700,
};

// Line heights
export const lineHeights = {
  tight: 1.25,
  snug: 1.375,
  normal: 1.5,
  relaxed: 1.625,
};

// Text styles (presets)
export const textStyles = {
  // Heading styles
  h1: {
    fontFamily: fonts.sans,
    fontSize: fontSizes['4xl'],
    fontWeight: fontWeights.bold,
    lineHeight: lineHeights.tight,
  },
  h2: {
    fontFamily: fonts.sans,
    fontSize: fontSizes['2xl'],
    fontWeight: fontWeights.semibold,
    lineHeight: lineHeights.tight,
  },
  h3: {
    fontFamily: fonts.sans,
    fontSize: fontSizes.xl,
    fontWeight: fontWeights.semibold,
    lineHeight: lineHeights.snug,
  },
  h4: {
    fontFamily: fonts.sans,
    fontSize: fontSizes.lg,
    fontWeight: fontWeights.semibold,
    lineHeight: lineHeights.snug,
  },
  h5: {
    fontFamily: fonts.sans,
    fontSize: fontSizes.base,
    fontWeight: fontWeights.medium,
    lineHeight: lineHeights.normal,
  },
  h6: {
    fontFamily: fonts.sans,
    fontSize: fontSizes.sm,
    fontWeight: fontWeights.medium,
    lineHeight: lineHeights.normal,
  },

  // Body text styles
  bodyLarge: {
    fontFamily: fonts.sans,
    fontSize: fontSizes.lg,
    fontWeight: fontWeights.normal,
    lineHeight: lineHeights.relaxed,
  },
  body: {
    fontFamily: fonts.sans,
    fontSize: fontSizes.base,
    fontWeight: fontWeights.normal,
    lineHeight: lineHeights.normal,
  },
  bodySmall: {
    fontFamily: fonts.sans,
    fontSize: fontSizes.sm,
    fontWeight: fontWeights.normal,
    lineHeight: lineHeights.normal,
  },

  // Special styles
  caption: {
    fontFamily: fonts.sans,
    fontSize: fontSizes.xs,
    fontWeight: fontWeights.normal,
    lineHeight: lineHeights.tight,
  },
  label: {
    fontFamily: fonts.sans,
    fontSize: fontSizes.sm,
    fontWeight: fontWeights.medium,
    lineHeight: lineHeights.normal,
  },
  button: {
    fontFamily: fonts.sans,
    fontSize: fontSizes.sm,
    fontWeight: fontWeights.medium,
    lineHeight: lineHeights.normal,
  },
  code: {
    fontFamily: fonts.mono,
    fontSize: '0.9em',
    fontWeight: fontWeights.normal,
    lineHeight: lineHeights.normal,
  },
};

// Tailwind CSS classes for text styles
export const textClasses = {
  h1: 'text-4xl font-bold leading-tight',
  h2: 'text-2xl font-semibold leading-tight',
  h3: 'text-xl font-semibold leading-snug',
  h4: 'text-lg font-semibold leading-snug',
  h5: 'text-base font-medium',
  h6: 'text-sm font-medium',

  bodyLarge: 'text-lg leading-relaxed',
  body: 'text-base',
  bodySmall: 'text-sm',

  caption: 'text-xs leading-tight',
  label: 'text-sm font-medium',
  button: 'text-sm font-medium',
  code: 'font-mono text-[0.9em]',
};

// Truncate text utilities
export const truncate = {
  single: 'truncate',
  twoLines: 'line-clamp-2',
  threeLines: 'line-clamp-3',
};

// Text alignment
export const textAlign = {
  left: 'text-left',
  center: 'text-center',
  right: 'text-right',
};

// Google Fonts import URL
export const fontImportURL = "https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap";
