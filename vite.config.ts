// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import path from "path"
import vue from "@vitejs/plugin-vue"
import { defineConfig } from "vite"
import viteComponents from "unplugin-vue-components/vite"
import { ElementPlusResolver } from "unplugin-vue-components/resolvers"
// import eslint from "@rollup/plugin-eslint"

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
    viteComponents({
      resolvers: [ElementPlusResolver({importStyle: "sass"})],
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
    rollupOptions: {
      // plugins: [visualizer({filename: "/Volumes/data/foo/s.html"})],
      // output: {
      //   manualChunks: {
      //     // element-plus is used in various chunks (because pages are loaded dynamically)
      //     echarts: ["echarts"],
      //     vue: ["vue"],
      //   }
      // },
    },
  },
})