import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"
import eslintPlugin from "vite-plugin-eslint"
// import visualizer from "rollup-plugin-visualizer"

// https://vitejs.dev/config/
// noinspection SpellCheckingInspection
export default defineConfig({
  plugins: [
    {...eslintPlugin({include: ["src/**/*.ts", "src/**/*.vue"], cache: false}), enforce: "pre"},
    vue()
  ],
  base: "/v2/",
  server: {
    host: "localhost",
    port: 8080,
  },
  build: {
    // sourcemap: true,
    emptyOutDir: true,
    chunkSizeWarningLimit: 600,
    outDir: "../cmd/frontend/kodata/v2",
    rollupOptions: {
      // plugins: [visualizer({filename: "/Volumes/data/foo/s.html"})],
      output: {
        manualChunks: {
          // element-plus is used in various chunks (because pages are loaded dynamically)
          echarts: ["echarts"],
          vue: ["vue"],
        }
      },
    },
  },
})
