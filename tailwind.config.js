const defaultTheme = require("tailwindcss/defaultTheme")

module.exports = {
  content: ["./index.html", "./dashboard/**/*.{vue,ts}"],
  theme: {
    extend: {
      fontFamily: {
        sans: ["InterVariable", ...defaultTheme.fontFamily.sans],
        mono: ["JetBrains MonoVariable", ...defaultTheme.fontFamily.mono],
      },
      colors: {
        // Lara Light Blue
        primary: {
          // active tab border color, input focus border (TailwindUI uses border-indigo-500)
          500: "#3B82F6",
          // active tab text color (TailwindUI uses text-indigo-600, but PrimeVue uses the same primary color)
          600: "#3B82F6",
        },
      },
    },
  },
  variants: {
    extend: {},
  },
  plugins: [
    require("@tailwindcss/line-clamp"),
    require("@tailwindcss/typography"),
    require("@tailwindcss/forms"),
  ],
}
