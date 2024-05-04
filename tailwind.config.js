/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/**/*.html"],
  theme: {
    extend: {
      fontFamily: {
        harlow: ["Harlow", "sans-serif"],
      },
      colors: {
        primary: "#222222",
        secondary: "rgb(217 119 6)",
        hover: "rgb(241 245 249)",
      }
    },
  },
  plugins: [],
}

