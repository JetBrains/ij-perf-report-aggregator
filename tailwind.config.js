const defaultTheme = require("tailwindcss/defaultTheme")

module.exports = {
  darkMode: ["class", ".dark-mode"],
  content: ["./index.html", "./dashboard/**/*.vue"],
  theme: {
    extend: {
      fontFamily: {
        sans: ["InterVariable", ...defaultTheme.fontFamily.sans],
        mono: ["JetBrains MonoVariable", ...defaultTheme.fontFamily.mono],
      },
      colors: {
        // Lara Light Blue
        primary: {
          DEFAULT: "#3B82F6", // Light mode primary
          dark: "#6495ED", // Dark mode primary
        },
        darker: "#1D4ED8", // Lara Dark Blue
      },
    },
  },
  variants: {
    extend: {},
  },
  plugins: [require("@tailwindcss/typography"), require("@tailwindcss/forms")],
}
