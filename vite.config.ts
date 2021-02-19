// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import path from "path"
import vue from "@vitejs/plugin-vue"
import { defineConfig } from "vite"
import viteComponents from "vite-plugin-components"
import eslintPlugin from "vite-plugin-eslint"

// import visualizer from "rollup-plugin-visualizer"

// https://vitejs.dev/config/
// noinspection SpellCheckingInspection,TypeScriptUnresolvedVariable
export default defineConfig({
  plugins: [
    {
      ...eslintPlugin({
        include: ["dashboard/**/*.ts", "jb/dashboard/jb/**/*.vue"],
        cache: false,
      }),
      enforce: "pre",
    },
    vue(),
    viteComponents({
      deep: false,
      customComponentResolvers: [
        name => {
          // eslint-disable-next-line @typescript-eslint/ban-ts-comment
          // @ts-ignore
          if (name.startsWith("El")) {
            const partialName = name[2].toLowerCase() + name.substring(3).replace(/[A-Z]/g, l => `-${l.toLowerCase()}`)
            return {
              path: `element-plus/es/el-${partialName}`,
              sideEffects: [`element-plus/packages/theme-chalk/src/${partialName}.scss`],
            }
          }
          return null
        },
      ],
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
      // remove once migration from amcharts to echarts will be completed
      external: ["xlsx", "pdfmake"]
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
