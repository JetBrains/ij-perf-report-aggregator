module.exports = {
  root: true,
  env: {
    node: true,
  },
  settings: {
    "import/extensions": [
      ".ts",
      ".vue",
    ],
  },
  extends: [
    "plugin:vue/vue3-essential",
    "plugin:vue/vue3-recommended",
    "plugin:vue/vue3-strongly-recommended",
    "plugin:import/recommended",
    "plugin:import/errors",
    "plugin:import/warnings",
    "plugin:import/typescript",
    "@vue/typescript/recommended",
  ],
  parserOptions: {
    ecmaVersion: 2020,
    warnOnUnsupportedTypeScriptVersion: false,
  },
  rules: {
    // "no-console": process.env.NODE_ENV === "production" ? ["warn", {allow: ["warn", "error"]}] : "off",
    "no-debugger": process.env.NODE_ENV === "production" ? "error" : "off",
    "max-len": ["error", {"code": 180}],
    "object-shorthand": ["error", "always", {"avoidExplicitReturnArrows": true}],
    "quotes": ["error", "double"],
    "@typescript-eslint/no-unused-vars": "off",
    "semi": "off",
    "import/order": [process.env.NODE_ENV === "production" ? "error" : "warn", {alphabetize: {order: "asc"}}],
    "import/no-unresolved": "off",
    "import/no-extraneous-dependencies": "error",
    "arrow-parens": ["error", "as-needed"],
    "@typescript-eslint/semi": ["error", "never"],
    "@typescript-eslint/no-inferrable-types": ["error", {"ignoreParameters": true}],
    "@typescript-eslint/member-delimiter-style": ["error", {
      "multiline": {
        "delimiter": "none",
      },
    }],
    "vue/html-quotes": ["error", "double", {"avoidEscape": true}],
    "no-restricted-imports": ["error", "element-plus", "echarts", "element-plus/es", "element-plus/lib/theme-chalk/index.css"],
  },
}
