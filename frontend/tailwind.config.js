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
        // Slate + Emerald palette
        // Primary accent - Emerald
        'primary': '#10b981',           // emerald-500
        'primary-hover': '#34d399',     // emerald-400
        'primary-dim': '#065f46',       // emerald-800

        // Slate palette for backgrounds
        'slate': {
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

        // Emerald palette for accents
        'emerald': {
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

        // Status Colors
        status: {
          starting: '#f59e0b',    // amber-500
          running: '#22c55e',     // green-500
          stopping: '#f59e0b',    // amber-500
          stopped: '#6b7280',     // gray-500
          error: '#ef4444',       // red-500
        },

        // Background colors for theme
        'bg-primary': {
          light: '#f8fafc',       // slate-50
          dark: '#1e293b',        // slate-800
        },
        'bg-secondary': {
          light: '#f1f5f9',       // slate-100
          dark: '#0f172a',        // slate-900
        },
        'bg-tertiary': {
          light: '#e2e8f0',       // slate-200
          dark: '#334155',        // slate-700
        },

        // Text colors
        'text-primary': {
          light: '#1e293b',       // slate-800
          dark: '#f8fafc',        // slate-50
        },
        'text-secondary': {
          light: '#475569',       // slate-600
          dark: '#cbd5e1',        // slate-300
        },
        'text-tertiary': {
          light: '#94a3b8',       // slate-400
          dark: '#94a3b8',        // slate-400
        },

        // Glass surface colors
        glass: {
          light: 'rgba(248, 250, 252, 0.95)',
          dark: 'rgba(30, 41, 59, 0.95)',
        },
      },
      backgroundImage: {
        // Subtle gradient for elevated surfaces
        'surface-gradient': 'linear-gradient(180deg, rgba(255,255,255,0.05) 0%, rgba(255,255,255,0) 100%)',
      },
      backdropBlur: {
        xs: '2px',
      },
      animation: {
        'pulse-ring': 'pulse-ring 2s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'fade-in': 'fade-in 200ms ease-out',
        'slide-up': 'slide-up 200ms ease-out',
      },
      keyframes: {
        'pulse-ring': {
          '0%, 100%': { opacity: '0', transform: 'scale(1)' },
          '50%': { opacity: '0.3', transform: 'scale(1.1)' },
        },
        'fade-in': {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        'slide-up': {
          '0%': { opacity: '0', transform: 'translateY(10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
      },
      fontFamily: {
        'sans': ['Inter', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'Helvetica', 'Arial', 'sans-serif'],
        'mono': ['JetBrains Mono', 'SF Mono', 'Fira Code', 'Consolas', 'monospace'],
      },
      boxShadow: {
        'sm': '0 1px 2px rgba(0, 0, 0, 0.4)',
        'md': '0 4px 6px -1px rgba(0, 0, 0, 0.5), 0 2px 4px -2px rgba(0, 0, 0, 0.5)',
        'lg': '0 10px 15px -3px rgba(0, 0, 0, 0.5), 0 4px 6px -4px rgba(0, 0, 0, 0.5)',
        'glass': '0 4px 30px rgba(0, 0, 0, 0.5), 0 0 0 1px rgba(255, 255, 255, 0.05)',
        'glow': '0 0 20px rgba(16, 185, 129, 0.15)',
      },
      borderRadius: {
        'sm': '4px',
        'md': '8px',
        'lg': '12px',
        'xl': '16px',
      },
    },
  },
  plugins: [],
}
