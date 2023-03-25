/** @type {import('tailwindcss').Config} */
module.exports = {
  prefix: 'tw-',
  important: true,
  content: ["./src/**/*.{html,js,ts,vue}"],
  theme: {
    fontFamily: {
      sans: ['ChakraPetch', 'sans-serif'],
      serif: ['Merriweather', 'serif'],
    },
    extend: {},
  },
  variants: {
    opacity: ['responsive', 'hover']
  },
  plugins: [],
};
