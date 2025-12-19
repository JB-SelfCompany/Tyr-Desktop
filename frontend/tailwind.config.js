/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        // Material Design 3 Colors - Flat structure for Tailwind
        // Light Theme
        'md-light-primary': '#006C4C',
        'md-light-onPrimary': '#FFFFFF',
        'md-light-primaryContainer': '#89F8C7',
        'md-light-onPrimaryContainer': '#002114',
        'md-light-secondary': '#4D6357',
        'md-light-onSecondary': '#FFFFFF',
        'md-light-secondaryContainer': '#CFE9D9',
        'md-light-onSecondaryContainer': '#0A1F16',
        'md-light-tertiary': '#3D6373',
        'md-light-onTertiary': '#FFFFFF',
        'md-light-tertiaryContainer': '#C1E8FB',
        'md-light-onTertiaryContainer': '#001F29',
        'md-light-error': '#BA1A1A',
        'md-light-errorContainer': '#FFDAD6',
        'md-light-onError': '#FFFFFF',
        'md-light-onErrorContainer': '#410002',
        'md-light-background': '#F5F7F5',
        'md-light-onBackground': '#1A1C1A',
        'md-light-surface': '#FEFFFE',
        'md-light-onSurface': '#1A1C1A',
        'md-light-surfaceVariant': '#DDE5DD',
        'md-light-onSurfaceVariant': '#414942',
        'md-light-outline': '#707973',

        // Dark Theme
        'md-dark-primary': '#6CDB9C',
        'md-dark-onPrimary': '#003826',
        'md-dark-primaryContainer': '#005138',
        'md-dark-onPrimaryContainer': '#89F8C7',
        'md-dark-secondary': '#B3CCBE',
        'md-dark-onSecondary': '#1F352A',
        'md-dark-secondaryContainer': '#354B40',
        'md-dark-onSecondaryContainer': '#CFE9D9',
        'md-dark-tertiary': '#A5CCDE',
        'md-dark-onTertiary': '#073543',
        'md-dark-tertiaryContainer': '#244C5B',
        'md-dark-onTertiaryContainer': '#C1E8FB',
        'md-dark-error': '#FFB4AB',
        'md-dark-errorContainer': '#93000A',
        'md-dark-onError': '#690005',
        'md-dark-onErrorContainer': '#FFDAD6',
        'md-dark-background': '#191C1A',
        'md-dark-onBackground': '#E8F0E8',
        'md-dark-surface': '#191C1A',
        'md-dark-onSurface': '#E8F0E8',
        'md-dark-surfaceVariant': '#404943',
        'md-dark-onSurfaceVariant': '#C1CFC1',
        'md-dark-outline': '#899389',
        // Status Colors (from Android)
        status: {
          starting: '#FF9800',
          running: '#4CAF50',
          stopping: '#FF9800',
          stopped: '#9E9E9E',
          error: '#F44336',
        },
        // Legacy Y2K Neon Colors (for backward compatibility)
        neon: {
          pink: '#FF00FF',
          cyan: '#00FFFF',
          green: '#00FF88',
          yellow: '#FFFF00',
          purple: '#9D00FF',
          blue: '#0080FF',
        },
        // Updated Background colors with better contrast
        'space-blue': {
          DEFAULT: '#191C1A',      // Dark theme background (from MD dark)
          light: '#2A3530',        // Slightly lighter for cards
          lighter: '#354B40',      // Even lighter variant
        },
        'soft-blue': {
          DEFAULT: '#E8F0E8',      // Light theme background (subtle green tint)
          light: '#FBFDF8',        // Almost white (from MD light)
        },
        // Glass surface colors
        glass: {
          light: 'rgba(255, 255, 255, 0.85)',
          dark: 'rgba(42, 53, 48, 0.85)',
        },
      },
      backgroundImage: {
        // Holographic gradient
        'holographic': 'linear-gradient(135deg, #FF00FF 0%, #00FFFF 25%, #00FF00 50%, #FFFF00 75%, #FF00FF 100%)',
        'holographic-border': 'linear-gradient(90deg, #FF00FF, #00FFFF, #00FF00, #FFFF00, #FF00FF)',
        // Iridescent gradient
        'iridescent': 'linear-gradient(90deg, #667eea 0%, #764ba2 25%, #f093fb 50%, #4facfe 75%, #00f2fe 100%)',
        // Chrome gradients
        'chrome-light': 'linear-gradient(180deg, #E0E0E0 0%, #FAFAFA 50%, #D0D0D0 100%)',
        'chrome-dark': 'linear-gradient(180deg, #3A3A3A 0%, #6A6A6A 50%, #3A3A3A 100%)',
        // Glow gradients
        'glow-pink': 'radial-gradient(circle, rgba(255,0,255,0.3) 0%, transparent 70%)',
        'glow-cyan': 'radial-gradient(circle, rgba(0,255,255,0.3) 0%, transparent 70%)',
        'glow-green': 'radial-gradient(circle, rgba(0,255,136,0.3) 0%, transparent 70%)',
      },
      backdropBlur: {
        xs: '2px',
      },
      animation: {
        'glow-pulse': 'glow-pulse 2s ease-in-out infinite',
        'holographic-spin': 'holographic-spin 3s linear infinite',
        'float': 'float 3s ease-in-out infinite',
        'shimmer': 'shimmer 2s linear infinite',
      },
      keyframes: {
        'glow-pulse': {
          '0%, 100%': { filter: 'drop-shadow(0 0 8px currentColor)' },
          '50%': { filter: 'drop-shadow(0 0 16px currentColor)' },
        },
        'holographic-spin': {
          '0%': { backgroundPosition: '0% 50%' },
          '100%': { backgroundPosition: '200% 50%' },
        },
        'float': {
          '0%, 100%': { transform: 'translateY(0px)' },
          '50%': { transform: 'translateY(-10px)' },
        },
        'shimmer': {
          '0%': { backgroundPosition: '-200% 0' },
          '100%': { backgroundPosition: '200% 0' },
        },
      },
      fontFamily: {
        'futuristic': ['Rajdhani', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'Helvetica', 'Arial', 'sans-serif'],
        'display': ['Orbitron', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'Helvetica', 'Arial', 'sans-serif'],
        'body': ['-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'Helvetica', 'Arial', 'sans-serif'],
        'mono': ['JetBrains Mono', 'Consolas', 'Monaco', 'Courier New', 'monospace'],
      },
      boxShadow: {
        'neon-pink': '0 0 20px rgba(255, 0, 255, 0.5)',
        'neon-cyan': '0 0 20px rgba(0, 255, 255, 0.5)',
        'neon-green': '0 0 20px rgba(0, 255, 136, 0.5)',
        'glass': '0 8px 32px 0 rgba(31, 38, 135, 0.37)',
      },
    },
  },
  plugins: [],
}

