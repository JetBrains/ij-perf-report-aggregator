const defaultTheme = require("tailwindcss/defaultTheme")
const colors = require("tailwindcss/colors")
module.exports = {
  content: ["./index.html", "./dashboard/**/*.{vue,ts}"],
  theme: {
    extend: {
      fontFamily: {
        sans: ["InterVariable", ...defaultTheme.fontFamily.sans],
        mono: ["JetBrains MonoVariable", ...defaultTheme.fontFamily.mono],
      },
      // colors: {
      //   // Lara Light Blue
      //   primary: {
      //     // active tab border color, input focus border (TailwindUI uses border-indigo-500)
      //     500: colors.indigo[500],
      //     // active tab text color (TailwindUI uses text-indigo-600, but PrimeVue uses the same primary color)
      //     600: colors.indigo[600],
      //   },
      // },
    },
  },
  plugins: [
    require("@tailwindcss/line-clamp"),
    require("@tailwindcss/typography"),
    require("@tailwindcss/forms"),
  ],
}
