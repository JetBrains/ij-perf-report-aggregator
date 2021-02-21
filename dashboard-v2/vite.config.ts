import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"
import eslintPlugin from "vite-plugin-eslint"

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    eslintPlugin({include: ["src/**/*.ts", "src/**/*.vue"]}),
    vue(),
  ],
  server: {
    host: "localhost",
    port: 8080,
  },
  build: {
    sourcemap: true,
    emptyOutDir: true,
    chunkSizeWarningLimit: 600,
    outDir: "../cmd/frontend/kodata/v2",
    rollupOptions: {
      output: {
        manualChunks: {
          "element-plus": ["element-plus"],
          "echarts": ["echarts"],
        },
      },
    },
  },
})
