// eslint-disable-next-line @typescript-eslint/ban-ts-comment
import vue from "@vitejs/plugin-vue"
// @ts-ignore
import { execFileSync } from "child_process"
// @ts-ignore
import path from "path"
// @ts-ignore
import { ComponentResolver } from "unplugin-vue-components"
import Components from "unplugin-vue-components/vite"
import { defineConfig, PluginOption } from "vite"
// import visualizer from "rollup-plugin-visualizer"
import { writeFile } from "fs/promises"
import * as zlib from "zlib"
import { OutputAsset, OutputChunk } from "rollup"
import { configDefaults } from "vitest/config"
import { viteStaticCopy } from "vite-plugin-static-copy"
import { PrimeVueResolver } from "@primevue/auto-import-resolver"
import tailwindcss from "@tailwindcss/vite"

// Resolve the build version exposed to the client error reporter for deploy correlation.
// On TeamCity the git checkout may be shallow/detached (or git unavailable in the build step), so
// `git rev-parse` fails there — TeamCity injects the revision and build number as environment
// variables instead. Order: explicit override → TeamCity VCS commit → TeamCity build number →
// local git → "unknown".
function resolveAppVersion(): string {
  // @ts-ignore
  const env = process.env as Record<string, string | undefined>

  const explicit = firstNonEmpty(env.APP_VERSION, env.BUILD_VCS_NUMBER, env.BUILD_NUMBER)
  if (explicit != null) {
    // Short-hash a full 40-char git SHA so the value matches the local-git format below and stays
    // well under the server's version-label length cap; leave human-readable build numbers intact.
    return /^[0-9a-f]{40}$/i.test(explicit) ? explicit.slice(0, 12) : explicit
  }

  try {
    return execFileSync("git", ["rev-parse", "--short", "HEAD"], { stdio: ["ignore", "pipe", "ignore"] })
      .toString()
      .trim()
  } catch {
    return "unknown"
  }
}

function firstNonEmpty(...values: (string | undefined)[]): string | undefined {
  for (const value of values) {
    const trimmed = value?.trim()
    if (trimmed != null && trimmed !== "") {
      return trimmed
    }
  }
  return undefined
}

// https://vitejs.dev/config/
// noinspection SpellCheckingInspection,TypeScriptUnresolvedVariable
export default defineConfig({
  define: {
    __APP_VERSION__: JSON.stringify(resolveAppVersion()),
  },
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
      dts: path.resolve(__dirname, "dashboard/new-dashboard/src/components.d.ts"),
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
          src: path.resolve(__dirname, "dashboard/new-dashboard/resources/projects"),
        },
      ],
    }),
  ],
  root: "dashboard/app",
  publicDir: path.resolve(__dirname, "dashboard/app/public"),
  server: {
    host: "localhost",
    port: 8080,
  },
  build: {
    // sourcemap: true,
    reportCompressedSize: false,
    emptyOutDir: true,
    chunkSizeWarningLimit: 600,
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    outDir: path.resolve(__dirname, "cmd/frontend/resources"),
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

async function brotliCompressFile(asset: OutputAsset | OutputChunk, outDir: string): Promise<void> {
  const file = path.join(outDir, asset.fileName)
  // woff2 is based on the Brotli compression algorithm - no need to compress
  if (file.endsWith(".png") || file.endsWith(".woff2")) {
    return
  }

  const data = Buffer.from("code" in asset ? asset.code : (asset as OutputAsset).source)
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
