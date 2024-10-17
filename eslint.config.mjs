import eslint from "@eslint/js"
import tseslint from "typescript-eslint"
import pluginVue from "eslint-plugin-vue"
import eslintConfigPrettier from "eslint-config-prettier"
import eslintPluginUnicorn from "eslint-plugin-unicorn"

export default tseslint.config(
  eslint.configs.recommended,
  ...tseslint.configs.strictTypeChecked,
  ...tseslint.configs.stylisticTypeChecked,
  ...pluginVue.configs["flat/recommended"],
  {
    plugins: {
      "typescript-eslint": tseslint.plugin,
      unicorn: eslintPluginUnicorn,
    },
    languageOptions: {
      parserOptions: {
        parser: tseslint.parser,
        project: "./tsconfig.json",
        extraFileExtensions: [".vue"],
        sourceType: "module",
      },
    },
  },
  eslintConfigPrettier,
  {
    ignores: [
      "**/components.d.ts",
      "**/auto-imports.d.ts",
      "**/vite.config.ts",
      "**/postcss.config.js",
      "**/tailwind.config.js",
      "**/eslint.config.mjs",
      "cmd/frontend/resources/**/*",
      "dashboard/new-dashboard/src/components/common/BranchIcon.vue",
      "dashboard/new-dashboard/src/components/common/SpaceIcon.vue",
      "dashboard/new-dashboard/tests/unit/dataquery.test.js",
    ],
  },
  {
    rules: {
      "no-debugger": "off",
      "max-len": [
        "error",
        {
          code: 300,
        },
      ],
      "object-shorthand": [
        "error",
        "always",
        {
          avoidExplicitReturnArrows: true,
        },
      ],
      quotes: [
        "error",
        "double",
        {
          avoidEscape: true,
        },
      ],
      "@typescript-eslint/no-empty-function": [
        "error",
        {
          allow: ["arrowFunctions"],
        },
      ],
      "@typescript-eslint/no-unused-vars": "off",
      "@typescript-eslint/prefer-regexp-exec": "off",
      "@typescript-eslint/restrict-template-expressions": [
        "error",
        {
          allowNullish: true,
        },
      ],
      "@typescript-eslint/no-inferrable-types": [
        "error",
        {
          ignoreParameters: true,
        },
      ],
      "@typescript-eslint/no-unsafe-enum-comparison": "off",
      "@typescript-eslint/non-nullable-type-assertion-style": "off",
      "vue/html-quotes": [
        "error",
        "double",
        {
          avoidEscape: true,
        },
      ],
      "vue/multi-word-component-names": [
        "error",
        {
          ignores: ["Dashboard", "Report", "Divider"],
        },
      ],
      "vue/no-setup-props-destructure": "off",
      "vue/no-deprecated-filter": "off",
      "unicorn/prefer-global-this": "off",
      "unicorn/prevent-abbreviations": "off",
      "unicorn/filename-case": "off",
      "unicorn/switch-case-braces": "off",
      "unicorn/no-null": "off",
      "unicorn/no-magic-array-flat-depth": "off",
      "unicorn/numeric-separators-style": "off",
      "unicorn/consistent-function-scoping": [
        "error",
        {
          checkArrowFunctions: false,
        },
      ],
      "unicorn/no-new-array": "off",
    },
  }
)
