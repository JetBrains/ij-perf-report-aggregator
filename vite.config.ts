// eslint-disable-next-line @typescript-eslint/ban-ts-comment
import vue from "@vitejs/plugin-vue"
// @ts-ignore
import path from "path"
import Components from "unplugin-vue-components/vite"
import { defineConfig, PluginOption, Rollup } from "vite"
import { writeFile } from "fs/promises"
import * as zlib from "zlib"
import { configDefaults } from "vitest/config"
import { viteStaticCopy } from "vite-plugin-static-copy"
import { PrimeVueResolver } from "@openvue/auto-import-resolver"
import tailwindcss from "@tailwindcss/vite"

// https://vitejs.dev/config/
// noinspection SpellCheckingInspection,TypeScriptUnresolvedVariable
export default defineConfig({
  test: {
    root: "dashboard/new-dashboard",
    include: [...configDefaults.include, "**/*.{test,spec}.ts"],
    globals: true,
    environment: "happy-dom",
    setupFiles: ["tests/setup.ts"],
    testTimeout: 10000,
  },
  plugins: [
    vue(),
    tailwindcss(),
    // visualizer({template: "sunburst"}),
    Components({
      directoryAsNamespace: true,
      dts: path.resolve(import.meta.dirname, "dashboard/new-dashboard/src/components.d.ts"),
      resolvers: [
        PrimeVueResolver(),
        // HeadlessUiResolver(),
        (name) => {
          // @ts-ignore
          const kind = process.env.NODE_ENV === "test" ? "" : "esm/"
          if (name.endsWith("Icon")) {
            return {
              path: `@heroicons/vue/24/outline/${kind}${name}.js`,
            }
          } else if (name.endsWith("IconSolid")) {
            return {
              path: `@heroicons/vue/20/solid/${kind}${name.substring(0, name.length - "Solid".length)}.js`,
            }
          } else {
            return null
          }
        },
      ],
    }),
    brotli(),
    viteStaticCopy({
      targets: [
        {
          dest: "../../degradation-analyzer/kodata",
          src: path.resolve(import.meta.dirname, "dashboard/new-dashboard/resources/projects"),
        },
      ],
    }),
  ],
  root: "dashboard/app",
  publicDir: path.resolve(import.meta.dirname, "dashboard/app/public"),
  server: {
    host: "localhost",
    port: 8080,
  },
  build: {
    // sourcemap: true,
    reportCompressedSize: false,
    emptyOutDir: true,
    chunkSizeWarningLimit: 600,
    rolldownOptions: {
      output: {
        codeSplitting: {
          groups: [
            { name: "zrender", test: /node_modules[\\/]zrender[\\/]/, priority: 20 },
            { name: "echarts", test: /node_modules[\\/]echarts[\\/]/, priority: 10 },
          ],
        },
      },
    },
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    outDir: path.resolve(import.meta.dirname, "cmd/frontend/resources"),
  },
  css: {
    preprocessorMaxWorkers: true,
  },
})

function brotli(): PluginOption {
  return {
    name: "offline-compression",
    writeBundle(outputOptions, bundle) {
      const outDir = outputOptions.dir!
      return Promise.all(Object.values(bundle).map((it) => brotliCompressFile(it, outDir))) as Promise<never>
    },
    apply: "build",
  }
}

async function brotliCompressFile(asset: Rollup.OutputAsset | Rollup.OutputChunk, outDir: string): Promise<void> {
  const file = path.join(outDir, asset.fileName)
  // woff2 is based on the Brotli compression algorithm - no need to compress
  // woff is zlib-compressed internally - brotli saves <1% and costs ~30ms per file
  if (file.endsWith(".png") || file.endsWith(".woff2") || file.endsWith(".woff")) {
    return
  }

  const data = Buffer.from("code" in asset ? asset.code : (asset as Rollup.OutputAsset).source)
  // https://github.com/google/ngx_brotli#brotli_min_length default is 20, so, we will compress any asset regardless of size
  if (data.length < 20) {
    throw new Error("Asset size is suspiciously small")
  }

  const mode = file.endsWith(".wasm") ? zlib.constants.BROTLI_MODE_GENERIC : zlib.constants.BROTLI_MODE_TEXT
  await new Promise((resolve, reject) => {
    zlib.brotliCompress(
      data,
      {
        params: {
          [zlib.constants.BROTLI_PARAM_MODE]: mode,
          [zlib.constants.BROTLI_PARAM_QUALITY]: zlib.constants.BROTLI_MAX_QUALITY,
        },
      },
      (error, buffer) => {
        if (error != null) {
          reject(error)
          return
        }

        writeFile(`${file}.br`, buffer).then(resolve, reject)
      }
    )
  })
}
