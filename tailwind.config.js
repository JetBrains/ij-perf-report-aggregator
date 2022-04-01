module.exports = {
  content: ["./index.html", "./dashboard/**/*.{vue,ts}"],
  darkMode: "class",
  mode: "jit",
  theme: {
    fontSize: {
      // we set rem to 14px for PrimeVue, but tailwind designed to work as is,
      // so, we use this as a workaround to use Tailwind UI templates as is
      "sm": "1rem",
      "base": "1.125rem;",
      "lg": "1.25rem;",
    },
    fontFamily: {
      sans: ["InterVariable"],
      mono: ["JetBrains MonoVariable"],
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
