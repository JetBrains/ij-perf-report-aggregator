// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import path from "path"
import vue from "@vitejs/plugin-vue"
import { ElementPlusResolver, PrimeVueResolver } from "unplugin-vue-components/resolvers"
import { defineConfig } from "vite"
import AutoImport from "unplugin-auto-import/vite"
import Components from "unplugin-vue-components/vite"

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
    AutoImport({
      resolvers: [ElementPlusResolver(), PrimeVueResolver() ],
    }),
    Components({
      resolvers: [ElementPlusResolver(), PrimeVueResolver()],
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