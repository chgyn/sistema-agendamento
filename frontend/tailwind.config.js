/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        brand: {
          50: '#fff8f0',
          100: '#ffedd5',
          200: '#fed7aa',
          300: '#fdba74',
          400: '#fb923c',
          500: '#f97316', // Vibrant Barber Amber / Orange
          600: '#ea580c',
          700: '#c2410c',
          800: '#9a3412',
          900: '#7c2d12',
          950: '#431407',
        },
        surface: {
          950: '#0c0d11', // Deep Obsidian background
          900: '#121318', // Surface base (cards)
          850: '#16171f', // Surface elevated
          800: '#1c1e27', // Surface highlight
          750: '#232532', // Borders & dividers
          700: '#2b2d3d',
          600: '#3f4257',
        },
      },
      fontFamily: {
        sans: ['Plus Jakarta Sans', 'Inter', 'system-ui', 'sans-serif'],
        display: ['Outfit', 'Plus Jakarta Sans', 'system-ui', 'sans-serif'],
      },
      boxShadow: {
        'glow': '0 0 25px -4px rgba(249, 115, 22, 0.35)',
        'glow-sm': '0 0 15px -3px rgba(249, 115, 22, 0.25)',
        'glow-emerald': '0 0 25px -5px rgba(16, 185, 129, 0.3)',
        'glow-purple': '0 0 25px -5px rgba(168, 85, 247, 0.35)',
      },
    },
  },
  plugins: [],
}
