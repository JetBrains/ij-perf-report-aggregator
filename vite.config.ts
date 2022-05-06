// eslint-disable-next-line @typescript-eslint/ban-ts-comment
import vue from "@vitejs/plugin-vue"
// @ts-ignore
import path from "path"
// @ts-ignore
import { ComponentResolver } from "unplugin-vue-components"
import Components from "unplugin-vue-components/vite"
import { defineConfig, PluginOption } from "vite"
import svgLoader from "vite-svg-loader"
// import visualizer from "rollup-plugin-visualizer"
import eslintPlugin from "@nabla/vite-plugin-eslint"
import { writeFile } from "fs/promises"
import * as zlib from "zlib"
import { OutputAsset, OutputChunk } from "rollup"

// https://vitejs.dev/config/
// noinspection SpellCheckingInspection,TypeScriptUnresolvedVariable
export default defineConfig({
  plugins: [
    eslintPlugin(),
    vue(),
    // visualizer({template: "sunburst"}),
    svgLoader(),
    Components({
      dts: path.resolve(__dirname, "dashboard/shared/src/components.d.ts"),
      resolvers: [
        PrimeVueResolver(),
        // HeadlessUiResolver(),
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
        },
      ],
    }),
    brotli(),
  ],
  root:  "dashboard/app",
  publicDir: path.resolve(__dirname, "dashboard/app/public"),
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
    outDir: path.resolve(__dirname, "cmd/frontend/resources"),
  },
})

function brotli(): PluginOption {
  let dir = ""
  return {
    name: "offline-compression",
    writeBundle(outputOptions, bundle) {
      const outDir = outputOptions.dir
      return Promise.all(Object.values(bundle).map(it => brotliCompressFile(it, outDir))) as Promise<never>
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
    zlib.brotliCompress(data, {
      params: {
        [zlib.constants.BROTLI_PARAM_MODE]: mode,
        [zlib.constants.BROTLI_PARAM_QUALITY]: zlib.constants.BROTLI_MAX_QUALITY,
      },
    }, (error, buffer) => {
      if (error != null) {
        reject(error)
        return
      }

      writeFile(`${file}.br`, buffer).then(resolve, reject)
    })
  })
}

// noinspection SpellCheckingInspection
const components = new Set<string>([
  "Accordion",
  "AccordionTab",
  "AutoComplete",
  "Avatar",
  "AvatarGroup",
  "Badge",
  "BlockUI",
  "Breadcrumb",
  "Button",
  "Calendar",
  "Card",
  "Carousel",
  "CascadeSelect",
  "Chart",
  "Checkbox",
  "Chip",
  "Chips",
  "ColorPicker",
  "Column",
  "ColumnGroup",
  // 'ConfirmDialog',
  // 'ConfirmPopup',
  // These must be registered globally in order for the confirm service to work properly
  "ContextMenu",
  "DataTable",
  "DataView",
  "DataViewLayoutOptions",
  "DeferredContent",
  "Dialog",
  "Divider",
  "Dock",
  "Dropdown",
  "Editor",
  "Fieldset",
  "FileUpload",
  "FullCalendar",
  "Galleria",
  "Image",
  "InlineMessage",
  "Inplace",
  "InputMask",
  "InputNumber",
  "InputSwitch",
  "InputText",
  "Knob",
  "Listbox",
  "MegaMenu",
  "Menu",
  "Menubar",
  "Message",
  "MultiSelect",
  "OrderList",
  "OrganizationChart",
  "OverlayPanel",
  "Paginator",
  "Panel",
  "PanelMenu",
  "Password",
  "PickList",
  "ProgressBar",
  "ProgressSpinner",
  "RadioButton",
  "Rating",
  "ScrollPanel",
  "ScrollTop",
  "SelectButton",
  "Sidebar",
  "Skeleton",
  "Slider",
  "SpeedDial",
  "SplitButton",
  "Splitter",
  "SplitterPanel",
  "Steps",
  "TabMenu",
  "TabPanel",
  "TabView",
  "Tag",
  "Terminal",
  "TerminalService",
  "Textarea",
  "TieredMenu",
  "Timeline",
  "Timelist",
  // 'Toast',
  // Toast must be registered globally in order for the toast service to work properly
  "ToggleButton",
  "Toolbar",
  // 'Tooltip',
  // Tooltip must be registered globally in order for the tooltip service to work properly
  "Tree",
  "TreeSelect",
  "TreeTable",
  "TriStateCheckbox",
  "VirtualScroller",
])

function PrimeVueResolver(): ComponentResolver {
  return {
    type: "component",
    resolve: (name: string) => {
      if (components.has(name)) {
        return {
          // SFC compilation saves 200KB (2.65 vs 2.67 MB)
          path: `primevue/${name.toLowerCase()}/sfc`,
        }
      }
      else {
        return null
      }
    },
  }
}
