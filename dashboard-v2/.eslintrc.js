module.exports = {
  root: true,
  env: {
    node: true,
  },
  extends: [
    "plugin:vue/vue3-essential",
    "plugin:vue/vue3-recommended",
    "plugin:vue/vue3-strongly-recommended",
    "@vue/typescript/recommended",
  ],
  parserOptions: {
    ecmaVersion: 2020,
    warnOnUnsupportedTypeScriptVersion: false,
  },
  rules: {
    // "no-console": process.env.NODE_ENV === "production" ? ["warn", {allow: ["warn", "error"]}] : "off",
    "no-debugger": process.env.NODE_ENV === "production" ? "error" : "off",
    "max-len": ["error", {"code": 140}],
    "object-shorthand": ["error", "always", {"avoidExplicitReturnArrows": true}],
    "quotes": ["error", "double"],
    "@typescript-eslint/no-unused-vars": "off",
    "semi": "off",
    "arrow-parens": ["error", "as-needed"],
    "@typescript-eslint/semi": ["error", "never"],
    "@typescript-eslint/no-inferrable-types": ["error", {"ignoreParameters": true}],
    "@typescript-eslint/member-delimiter-style": ["error", {
      "multiline": {
        "delimiter": "none",
      },
    }],
    "vue/html-quotes": ["error", "double", {"avoidEscape": true}],
    "no-restricted-imports": ["error", "element-plus", "element-plus/es"]
  },
}
