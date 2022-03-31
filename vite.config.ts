// eslint-disable-next-line @typescript-eslint/ban-ts-comment
import vue from "@vitejs/plugin-vue"
// @ts-ignore
import path from "path"
import AutoImport from "unplugin-auto-import/vite"
import { PrimeVueResolver, HeadlessUiResolver } from "unplugin-vue-components/resolvers"
import Components from "unplugin-vue-components/vite"
import { defineConfig } from "vite"
import svgLoader from "vite-svg-loader"
// import visualizer from "rollup-plugin-visualizer"

// https://vitejs.dev/config/
// noinspection SpellCheckingInspection,TypeScriptUnresolvedVariable
export default defineConfig({
  plugins: [
    // {
    //   ...eslintPlugin({
    //     include: ["dashboard/**/*.ts", "jb/dashboard/jb/**/*.vue"],
    //     cache: false,
    //   }),
    //   enforce: "pre",
    // },
    vue(),
    svgLoader(),
    AutoImport({
      resolvers: [PrimeVueResolver()],
    }),
    Components({
      resolvers: [
        PrimeVueResolver(),
        HeadlessUiResolver(),
        name => {
          // @ts-ignore
          const kind = process.env.NODE_ENV === "test" ? "" : "esm/"
          if (name.endsWith("Icon")) {
            return {
              path: `@heroicons/vue/outline/${kind}${name}.js`,
            }
          }
          else if (name.endsWith("IconSolid")) {
            return {
              path: `@heroicons/vue/solid/${kind}${name.substring(0, name.length - "Solid".length)}.js`,
            }
          }
          else {
            return null
          }
        }],
    }),
  ],
  root: "dashboard/app",
  publicDir: "dashboard/app/public",
  server: {
    host: "localhost",
    port: 8080,
  },
  build: {
    // sourcemap: true,
    emptyOutDir: true,
    chunkSizeWarningLimit: 600,
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    outDir: path.resolve(__dirname, "cmd/frontend/kodata"),
  },
})