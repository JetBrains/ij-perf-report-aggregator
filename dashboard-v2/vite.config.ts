import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"
import eslintPlugin from "vite-plugin-eslint"
// import visualizer from "rollup-plugin-visualizer"

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    {...eslintPlugin({include: ["src/**/*.ts", "src/**/*.vue"]}), enforce: "pre"},
    vue()
  ],
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
      // plugins: [visualizer({filename: "/Volumes/data/s.html"})],
      output: {
        manualChunks: function (id: string): string | undefined {
          if (id.includes("element-plus")) {
            return "element-plus"
          }
          else if (id.includes("echarts")) {
            return "echarts"
          }
          else {
            return undefined
          }
        },
      },
    },
  },
})
