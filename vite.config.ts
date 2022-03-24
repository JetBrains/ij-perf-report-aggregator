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

// export function ElementPlusResolver(): ComponentResolver {
//   return (name: string) => {
//     if (name.match(/^El[A-Z]/)) {
//       // ElTableColumn -> table-column
//       const partialName = kebabCase(name.slice(2))
//       return {
//         importName: name,
//         path: "element-plus/es",
//         sideEffects: `element-plus/theme-chalk/src/${partialName == "sub-menu" ? "submenu" : partialName}.scss`,
//       }
//     }
//   }
// }