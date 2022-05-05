module.exports = {
  root: true,
  env: {
    browser: true,
    "vue/setup-compiler-macros": true,
  },
  extends: [
    "plugin:import/recommended",
    "plugin:import/errors",
    "plugin:import/warnings",
    "plugin:import/typescript",
    "plugin:vue/vue3-essential",
    "plugin:vue/vue3-recommended",
    "plugin:vue/vue3-strongly-recommended",
    "eslint:recommended",
    "@vue/typescript/recommended",
    "plugin:@typescript-eslint/eslint-recommended",
    "plugin:@typescript-eslint/recommended",
    "plugin:@typescript-eslint/recommended-requiring-type-checking"
  ],
  // plugins: ["simple-import-sort"],
  parser: "vue-eslint-parser",
  parserOptions: {
    ecmaVersion: 2020,
    project: ["./dashboard/**/tsconfig.json"],
    parser: "@typescript-eslint/parser"
  },
  rules: {
    // "no-console": process.env.NODE_ENV === "production" ? ["warn", {allow: ["warn", "error"]}] : "off",
    "no-debugger": process.env.NODE_ENV === "production" ? "error" : "off",
    "max-len": ["error", {"code": 180}],
    "object-shorthand": ["error", "always", {"avoidExplicitReturnArrows": true}],
    "quotes": ["error", "double", {avoidEscape: true}],
    "@typescript-eslint/no-unused-vars": "off",
    "semi": "off",
    "import/order": [process.env.NODE_ENV === "production" ? "error" : "warn", {alphabetize: {order: "asc"}}],
    "import/no-unresolved": "off",
    "import/no-extraneous-dependencies": "error",
    "arrow-parens": ["error", "as-needed"],
    "@typescript-eslint/semi": ["error", "never"],
    "@typescript-eslint/restrict-template-expressions": ["error", {allowNullish: true}],
    "@typescript-eslint/no-inferrable-types": ["error", {"ignoreParameters": true}],
    "@typescript-eslint/member-delimiter-style": ["error", {
      "multiline": {
        "delimiter": "none",
      },
    }],
    "vue/html-quotes": ["error", "double", {"avoidEscape": true}],
    "no-restricted-imports": ["error",  "echarts", "../shared", "../../shared", "rxjs/operators"],
    "vue/multi-word-component-names": ["error", {
      "ignores": ["Dashboard", "Report"]
    }]
  },
}