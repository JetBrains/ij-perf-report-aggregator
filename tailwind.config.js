const defaultTheme = require("tailwindcss/defaultTheme")
module.exports = {
  content: ["./index.html", "./dashboard/**/*.{vue,ts}"],
  darkMode: "class",
  mode: "jit",
  theme: {
    extend: {
      fontFamily: {
        sans: ["InterVariable", ...defaultTheme.fontFamily.sans],
        mono: ["JetBrains MonoVariable", ...defaultTheme.fontFamily.mono],
      },
    },
  },
  variants: {
    extend: {},
  },
  plugins: [
    require("@tailwindcss/line-clamp"),
    require('@tailwindcss/typography'),
  ],
}
