// Typography system for Y2K Futurism

// Font families
export const fonts = {
  // Display font for headings (futuristic, geometric)
  display: "'Orbitron', 'Rajdhani', sans-serif",

  // Body font (clean, readable)
  body: "'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', sans-serif",

  // Futuristic font for special elements
  futuristic: "'Rajdhani', 'Orbitron', sans-serif",

  // Monospace for code/technical data
  mono: "'JetBrains Mono', 'Fira Code', 'Courier New', monospace",
};

// Font sizes (following a modular scale)
export const fontSizes = {
  xs: '0.75rem',     // 12px
  sm: '0.875rem',    // 14px
  base: '1rem',      // 16px
  lg: '1.125rem',    // 18px
  xl: '1.25rem',     // 20px
  '2xl': '1.5rem',   // 24px
  '3xl': '1.875rem', // 30px
  '4xl': '2.25rem',  // 36px
  '5xl': '3rem',     // 48px
  '6xl': '3.75rem',  // 60px
  '7xl': '4.5rem',   // 72px
};

// Font weights
export const fontWeights = {
  thin: 100,
  extralight: 200,
  light: 300,
  normal: 400,
  medium: 500,
  semibold: 600,
  bold: 700,
  extrabold: 800,
  black: 900,
};

// Line heights
export const lineHeights = {
  none: 1,
  tight: 1.25,
  snug: 1.375,
  normal: 1.5,
  relaxed: 1.625,
  loose: 2,
};

// Letter spacing
export const letterSpacings = {
  tighter: '-0.05em',
  tight: '-0.025em',
  normal: '0',
  wide: '0.025em',
  wider: '0.05em',
  widest: '0.1em',
};

// Text styles (presets)
export const textStyles = {
  // Heading styles
  h1: {
    fontFamily: fonts.display,
    fontSize: fontSizes['5xl'],
    fontWeight: fontWeights.bold,
    lineHeight: lineHeights.tight,
    letterSpacing: letterSpacings.tight,
  },
  h2: {
    fontFamily: fonts.display,
    fontSize: fontSizes['4xl'],
    fontWeight: fontWeights.bold,
    lineHeight: lineHeights.tight,
    letterSpacing: letterSpacings.tight,
  },
  h3: {
    fontFamily: fonts.display,
    fontSize: fontSizes['3xl'],
    fontWeight: fontWeights.semibold,
    lineHeight: lineHeights.snug,
    letterSpacing: letterSpacings.normal,
  },
  h4: {
    fontFamily: fonts.display,
    fontSize: fontSizes['2xl'],
    fontWeight: fontWeights.semibold,
    lineHeight: lineHeights.snug,
  },
  h5: {
    fontFamily: fonts.display,
    fontSize: fontSizes.xl,
    fontWeight: fontWeights.medium,
    lineHeight: lineHeights.normal,
  },
  h6: {
    fontFamily: fonts.display,
    fontSize: fontSizes.lg,
    fontWeight: fontWeights.medium,
    lineHeight: lineHeights.normal,
  },

  // Body text styles
  bodyLarge: {
    fontFamily: fonts.body,
    fontSize: fontSizes.lg,
    fontWeight: fontWeights.normal,
    lineHeight: lineHeights.relaxed,
  },
  body: {
    fontFamily: fonts.body,
    fontSize: fontSizes.base,
    fontWeight: fontWeights.normal,
    lineHeight: lineHeights.normal,
  },
  bodySmall: {
    fontFamily: fonts.body,
    fontSize: fontSizes.sm,
    fontWeight: fontWeights.normal,
    lineHeight: lineHeights.normal,
  },

  // Special styles
  caption: {
    fontFamily: fonts.body,
    fontSize: fontSizes.xs,
    fontWeight: fontWeights.normal,
    lineHeight: lineHeights.tight,
  },
  overline: {
    fontFamily: fonts.futuristic,
    fontSize: fontSizes.xs,
    fontWeight: fontWeights.semibold,
    lineHeight: lineHeights.none,
    letterSpacing: letterSpacings.widest,
    textTransform: 'uppercase' as const,
  },
  button: {
    fontFamily: fonts.futuristic,
    fontSize: fontSizes.base,
    fontWeight: fontWeights.semibold,
    lineHeight: lineHeights.none,
    letterSpacing: letterSpacings.wide,
  },
  code: {
    fontFamily: fonts.mono,
    fontSize: fontSizes.sm,
    fontWeight: fontWeights.normal,
    lineHeight: lineHeights.normal,
  },

  // Futuristic display text
  displayLarge: {
    fontFamily: fonts.display,
    fontSize: fontSizes['7xl'],
    fontWeight: fontWeights.black,
    lineHeight: lineHeights.none,
    letterSpacing: letterSpacings.tight,
  },
  displayMedium: {
    fontFamily: fonts.display,
    fontSize: fontSizes['6xl'],
    fontWeight: fontWeights.bold,
    lineHeight: lineHeights.tight,
    letterSpacing: letterSpacings.tight,
  },
};

// Tailwind CSS classes for text styles
export const textClasses = {
  h1: 'font-display text-5xl font-bold leading-tight tracking-tight',
  h2: 'font-display text-4xl font-bold leading-tight tracking-tight',
  h3: 'font-display text-3xl font-semibold leading-snug',
  h4: 'font-display text-2xl font-semibold leading-snug',
  h5: 'font-display text-xl font-medium',
  h6: 'font-display text-lg font-medium',

  bodyLarge: 'font-body text-lg leading-relaxed',
  body: 'font-body text-base',
  bodySmall: 'font-body text-sm',

  caption: 'font-body text-xs leading-tight',
  overline: 'font-futuristic text-xs font-semibold uppercase tracking-widest',
  button: 'font-futuristic text-base font-semibold tracking-wide',
  code: 'font-mono text-sm',

  displayLarge: 'font-display text-7xl font-black leading-none tracking-tight',
  displayMedium: 'font-display text-6xl font-bold leading-tight tracking-tight',
};

// Holographic text style (static, can be combined with any text style)
export const holographicTextStyle = {
  background: 'linear-gradient(90deg, #667eea 0%, #764ba2 25%, #f093fb 50%, #4facfe 75%, #00f2fe 100%)',
  WebkitBackgroundClip: 'text',
  WebkitTextFillColor: 'transparent',
  backgroundClip: 'text',
  backgroundSize: '200% 100%',
};

// Neon text glow effect
export const neonTextGlow = (color: string) => ({
  textShadow: `0 0 10px ${color}, 0 0 20px ${color}, 0 0 30px ${color}`,
});

// Gradient text utilities
export const gradientText = {
  holographic: 'bg-iridescent bg-clip-text text-transparent bg-[length:200%_100%]',
  neonPink: 'bg-gradient-to-r from-neon-pink to-neon-purple bg-clip-text text-transparent',
  neonCyan: 'bg-gradient-to-r from-neon-cyan to-neon-blue bg-clip-text text-transparent',
  neonGreen: 'bg-gradient-to-r from-neon-green to-neon-cyan bg-clip-text text-transparent',
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
  justify: 'text-justify',
};

// Google Fonts import URL (add to HTML head)
export const fontImportURL = "https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&family=Orbitron:wght@400;500;600;700;900&family=Rajdhani:wght@300;400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap";
