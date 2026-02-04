// Slate + Emerald Color Palette

export const colors = {
  // Primary accent - Emerald
  emerald: {
    50: '#ecfdf5',
    100: '#d1fae5',
    200: '#a7f3d0',
    300: '#6ee7b7',
    400: '#34d399',
    500: '#10b981',
    600: '#059669',
    700: '#047857',
    800: '#065f46',
    900: '#064e3b',
  },

  // Slate palette for neutral colors
  slate: {
    50: '#f8fafc',
    100: '#f1f5f9',
    200: '#e2e8f0',
    300: '#cbd5e1',
    400: '#94a3b8',
    500: '#64748b',
    600: '#475569',
    700: '#334155',
    800: '#1e293b',
    900: '#0f172a',
    950: '#020617',
  },

  // Light Theme Colors
  light: {
    background: {
      primary: '#f8fafc',    // slate-50 - cards
      secondary: '#f1f5f9',  // slate-100 - page background
      tertiary: '#e2e8f0',   // slate-200 - inputs
    },
    glass: {
      surface: 'rgba(248, 250, 252, 0.95)',
      border: 'rgba(203, 213, 225, 0.6)',
      hover: 'rgba(248, 250, 252, 0.98)',
    },
    text: {
      primary: '#1e293b',    // slate-800
      secondary: '#475569',  // slate-600
      tertiary: '#94a3b8',   // slate-400
    },
    primary: '#10b981',      // emerald-500
    primaryHover: '#34d399', // emerald-400
  },

  // Dark Theme Colors
  dark: {
    background: {
      primary: '#1e293b',    // slate-800 - cards
      secondary: '#0f172a',  // slate-900 - page background
      tertiary: '#334155',   // slate-700 - inputs
    },
    glass: {
      surface: 'rgba(30, 41, 59, 0.95)',
      border: 'rgba(51, 65, 85, 0.6)',
      hover: 'rgba(30, 41, 59, 0.98)',
    },
    text: {
      primary: '#f8fafc',    // slate-50
      secondary: '#cbd5e1',  // slate-300
      tertiary: '#94a3b8',   // slate-400
    },
    primary: '#10b981',      // emerald-500
    primaryHover: '#34d399', // emerald-400
  },

  // Status Colors
  status: {
    success: '#22c55e',      // green-500
    error: '#ef4444',        // red-500
    warning: '#f59e0b',      // amber-500
    info: '#3b82f6',         // blue-500
    stopped: '#6b7280',      // gray-500
  },

  // Border colors
  border: {
    light: '#cbd5e1',        // slate-300
    dark: '#334155',         // slate-700
    focus: '#10b981',        // emerald-500
  },
} as const;

// Type-safe color getter
export type ColorTheme = 'light' | 'dark';
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

export const getStatusColor = (status: StatusColor) => {
  return colors.status[status];
};

export const getPrimaryColor = (theme: ColorTheme) => {
  return colors[theme].primary;
};

export const getBorderColor = (theme: ColorTheme) => {
  return theme === 'light' ? colors.border.light : colors.border.dark;
};
