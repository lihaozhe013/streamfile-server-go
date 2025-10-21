/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./public/**/*.{html,js}",
    "./files/**/*.html",
    "!./src/frontend/markdown-viewer/**"
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
