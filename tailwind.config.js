const defaultTheme = require("tailwindcss/defaultTheme")

module.exports = {
  content: ["./index.html", "./dashboard/**/*.vue"],
  theme: {
    extend: {
      fontFamily: {
        sans: ["InterVariable", ...defaultTheme.fontFamily.sans],
        mono: ["JetBrains MonoVariable", ...defaultTheme.fontFamily.mono],
      },
      colors: {
        // Lara Light Blue
        primary: "#3B82F6",
        darker: "#1D4ED8", // Lara Dark Blue
      },
    },
  },
  variants: {
    extend: {},
  },
  plugins: [require("@tailwindcss/typography"), require("@tailwindcss/forms")],
}
