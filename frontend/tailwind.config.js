/** @type {import('tailwindcss').Config} */
// Make optional plugins robust during local dev
const plugins = []
try { plugins.push(require('@tailwindcss/line-clamp')) } catch {}

export default {
  content: [
    "./index.html",
    "./src/**/*.{ts,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        brand: {
          DEFAULT: '#0ea5e9',
          dark: '#0b83b8',
        },
      },
      fontFamily: {
        display: ['Inter', 'ui-sans-serif', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'SFMono-Regular', 'Menlo', 'monospace']
      },
      boxShadow: {
        'soft': '0 10px 25px -10px rgba(0,0,0,0.25)'
      }
    },
  },
  plugins,
}
