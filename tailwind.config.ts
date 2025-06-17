import type { Config } from 'tailwindcss'

const config: Config = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        'yt-dark': '#0F0F0F',
        'yt-light': '#FFFFFF',
        'yt-sidebar': '#212121',
        'yt-hover': '#272727'
      }
    },
  },
  plugins: [],
}
export default config 