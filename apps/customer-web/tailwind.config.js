/** @type {import('tailwindcss').Config} */
import forms from '@tailwindcss/forms';
import typography from '@tailwindcss/typography';
import aspectRatio from '@tailwindcss/aspect-ratio';

const config = {
  darkMode: ['class'],
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    container: {
      center: true,
      padding: '2rem',
      screens: {
        '2xl': '1400px',
      },
    },
    extend: {
      colors: {
        primary: {
          DEFAULT: 'hsl(24, 94%, 50%)',
          foreground: 'hsl(0, 0%, 100%)',
          dark: 'hsl(21, 90%, 48%)',
          light: 'hsl(24, 100%, 55%)',
        },
        secondary: {
          DEFAULT: 'hsl(0, 0%, 100%)',
          foreground: 'hsl(0, 0%, 7%)',
        },
        accent: {
          DEFAULT: 'hsl(0, 0%, 96%)',
          foreground: 'hsl(0, 0%, 13%)',
        },
        destructive: {
          DEFAULT: 'hsl(0, 84.2%, 60.2%)',
          foreground: 'hsl(0, 0%, 100%)',
        },
        muted: {
          DEFAULT: 'hsl(0, 0%, 96.1%)',
          foreground: 'hsl(0, 0%, 45.1%)',
        },
        border: 'hsl(0, 0%, 89.8%)',
        input: 'hsl(0, 0%, 89.8%)',
        ring: 'hsl(24, 94%, 50%)',
        background: 'hsl(0, 0%, 100%)',
        foreground: 'hsl(0, 0%, 3.9%)',
        card: {
          DEFAULT: 'hsl(0, 0%, 100%)',
          foreground: 'hsl(0, 0%, 3.9%)',
        },
      },
      borderRadius: {
        lg: 'var(--radius)',
        md: 'calc(var(--radius) - 2px)',
        sm: 'calc(var(--radius) - 4px)',
      },
      keyframes: {
        'accordion-down': {
          from: { height: 0 },
          to: { height: 'var(--radix-accordion-content-height)' },
        },
        'accordion-up': {
          from: { height: 'var(--radix-accordion-content-height)' },
          to: { height: 0 },
        },
      },
      animation: {
        'accordion-down': 'accordion-down 0.2s ease-out',
        'accordion-up': 'accordion-up 0.2s ease-out',
      },
      boxShadow: {
        amazon: '0 2px 5px 0 rgba(213, 217, 217, 0.5)',
        'amazon-hover': '0 2px 8px 0 rgba(213, 217, 217, 0.8)',
      },
      spacing: {
        'nav-height': '60px',
        'header-height': '115px',
      },
    },
  },
  plugins: [forms, typography, aspectRatio],
};

export default config;
