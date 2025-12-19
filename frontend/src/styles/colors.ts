// Y2K Futurism Color Palette

export const colors = {
  // Y2K Neon Colors
  neon: {
    pink: '#FF00FF',
    cyan: '#00FFFF',
    green: '#00FF88',
    yellow: '#FFFF00',
    purple: '#9D00FF',
    blue: '#0080FF',
  },

  // Light Theme Colors
  light: {
    background: {
      primary: '#F0F4FF',
      secondary: '#FFFFFF',
      tertiary: '#E8ECFF',
    },
    glass: {
      surface: 'rgba(255, 255, 255, 0.7)',
      border: 'rgba(255, 255, 255, 0.3)',
      hover: 'rgba(255, 255, 255, 0.85)',
    },
    text: {
      primary: '#0A0E1A',
      secondary: '#4A5568',
      tertiary: '#A0AEC0',
    },
    primary: '#6366F1',
    primaryHover: '#4F46E5',
  },

  // Dark Theme Colors
  dark: {
    background: {
      primary: '#0A0E1A',
      secondary: '#1A1F36',
      tertiary: '#2A2F46',
    },
    glass: {
      surface: 'rgba(26, 31, 54, 0.7)',
      border: 'rgba(255, 255, 255, 0.1)',
      hover: 'rgba(26, 31, 54, 0.85)',
    },
    text: {
      primary: '#FFFFFF',
      secondary: '#CBD5E0',
      tertiary: '#A0AEC0',
    },
    primary: '#818CF8',
    primaryHover: '#6366F1',
  },

  // Metallic Colors
  metallic: {
    chrome: {
      light: 'linear-gradient(180deg, #E0E0E0 0%, #FAFAFA 50%, #D0D0D0 100%)',
      dark: 'linear-gradient(180deg, #3A3A3A 0%, #6A6A6A 50%, #3A3A3A 100%)',
    },
    silver: '#C0C0C0',
    gold: '#FFD700',
  },

  // Status Colors
  status: {
    success: '#00FF88',
    error: '#FF0055',
    warning: '#FFD700',
    info: '#00FFFF',
  },

  // Gradient Definitions
  gradients: {
    holographic: 'linear-gradient(135deg, #FF00FF 0%, #00FFFF 25%, #00FF00 50%, #FFFF00 75%, #FF00FF 100%)',
    holographicBorder: 'linear-gradient(90deg, #FF00FF, #00FFFF, #00FF00, #FFFF00, #FF00FF)',
    iridescent: 'linear-gradient(90deg, #667eea 0%, #764ba2 25%, #f093fb 50%, #4facfe 75%, #00f2fe 100%)',
    neonPink: 'radial-gradient(circle, rgba(255,0,255,0.5) 0%, transparent 70%)',
    neonCyan: 'radial-gradient(circle, rgba(0,255,255,0.5) 0%, transparent 70%)',
    neonGreen: 'radial-gradient(circle, rgba(0,255,136,0.5) 0%, transparent 70%)',
  },
} as const;

// Type-safe color getter
export type ColorTheme = 'light' | 'dark';
export type NeonColor = keyof typeof colors.neon;
export type StatusColor = keyof typeof colors.status;

// Helper functions
export const getGlassColor = (theme: ColorTheme, variant: 'surface' | 'border' | 'hover') => {
  return colors[theme].glass[variant];
};

export const getTextColor = (theme: ColorTheme, variant: 'primary' | 'secondary' | 'tertiary') => {
  return colors[theme].text[variant];
};

export const getBackgroundColor = (theme: ColorTheme, variant: 'primary' | 'secondary' | 'tertiary') => {
  return colors[theme].background[variant];
};

export const getNeonColor = (color: NeonColor) => {
  return colors.neon[color];
};

export const getStatusColor = (status: StatusColor) => {
  return colors.status[status];
};
