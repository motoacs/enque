/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // Dark cinema palette
        base: {
          DEFAULT: '#0a0a0f',
          50: '#f5f5f7',
        },
        surface: {
          DEFAULT: '#12121a',
          light: '#1a1a25',
        },
        elevated: {
          DEFAULT: '#1e1e2a',
          light: '#252533',
        },
        panel: '#16161f',
        // Accent colors
        amber: {
          DEFAULT: '#e8a849',
          hover: '#f0b85c',
          dim: '#d4922a',
          muted: 'rgba(232, 168, 73, 0.12)',
          glow: 'rgba(232, 168, 73, 0.25)',
        },
        cyan: {
          DEFAULT: '#4ecdc4',
          hover: '#5fd9d1',
          dim: '#3bb5ad',
          muted: 'rgba(78, 205, 196, 0.12)',
          glow: 'rgba(78, 205, 196, 0.3)',
        },
        // Semantic
        success: { DEFAULT: '#34d399', muted: 'rgba(52, 211, 153, 0.12)' },
        danger: { DEFAULT: '#f87171', muted: 'rgba(248, 113, 113, 0.12)' },
        warn: { DEFAULT: '#fbbf24', muted: 'rgba(251, 191, 36, 0.12)' },
        // Borders
        edge: {
          DEFAULT: 'rgba(255, 255, 255, 0.07)',
          strong: 'rgba(255, 255, 255, 0.12)',
          accent: 'rgba(232, 168, 73, 0.3)',
        },
        // Text
        txt: {
          DEFAULT: '#e8e6e3',
          secondary: '#9d9da7',
          tertiary: '#5c5c68',
          muted: '#3c3c48',
        },
      },
      fontFamily: {
        display: ['Syne', 'system-ui', 'sans-serif'],
        body: ['DM Sans', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'Consolas', 'monospace'],
      },
      boxShadow: {
        'glow-amber': '0 0 20px rgba(232, 168, 73, 0.15), 0 0 60px rgba(232, 168, 73, 0.05)',
        'glow-cyan': '0 0 20px rgba(78, 205, 196, 0.15), 0 0 60px rgba(78, 205, 196, 0.05)',
        'glow-sm-amber': '0 0 10px rgba(232, 168, 73, 0.2)',
        'glow-sm-cyan': '0 0 10px rgba(78, 205, 196, 0.2)',
        'panel': '0 4px 24px rgba(0, 0, 0, 0.3), 0 1px 3px rgba(0, 0, 0, 0.2)',
        'dialog': '0 8px 40px rgba(0, 0, 0, 0.5), 0 2px 8px rgba(0, 0, 0, 0.3)',
      },
      animation: {
        'pulse-glow': 'pulse-glow 2.5s ease-in-out infinite',
        'fade-in': 'fade-in 0.2s ease-out',
        'slide-up': 'slide-up 0.3s ease-out',
      },
      keyframes: {
        'pulse-glow': {
          '0%, 100%': { boxShadow: '0 0 15px rgba(232, 168, 73, 0.15), 0 0 45px rgba(232, 168, 73, 0.05)' },
          '50%': { boxShadow: '0 0 25px rgba(232, 168, 73, 0.3), 0 0 70px rgba(232, 168, 73, 0.1)' },
        },
        'fade-in': {
          from: { opacity: '0' },
          to: { opacity: '1' },
        },
        'slide-up': {
          from: { opacity: '0', transform: 'translateY(8px)' },
          to: { opacity: '1', transform: 'translateY(0)' },
        },
      },
    },
  },
  plugins: [],
}
