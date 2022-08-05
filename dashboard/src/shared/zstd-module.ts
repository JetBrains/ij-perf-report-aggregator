import { from, map, shareReplay } from "rxjs"
import zstdWasmUrl from "./zstd.wasm?url"

function abort(what: string): void {
  throw new Error(what)
}

const UTF8Decoder = new TextDecoder()

export let HEAPU8: Uint8Array
export let buffer: ArrayBuffer

function UTF8ToString(ptr: number, maxBytesToRead: number = 0) {
  if (!ptr) {
    return ""
  }

  const maxPtr = ptr + maxBytesToRead
  let end = ptr
  for (; !(end >= maxPtr) && HEAPU8[end];) {
    ++end
  }
  return UTF8Decoder.decode(HEAPU8.subarray(ptr, end))
}

const asmLibraryArg = {
  a(condition: number, filename: number, line: number, func?: number) {
    abort(`Assertion failed: ${UTF8ToString(condition)}, at: ` +
      `[${filename ? UTF8ToString(filename) : "unknown filename"}, ${line}, ${func ? UTF8ToString(func) : "unknown function"}]`)
  },
  c(dest: number, src: number, num: number) {
    HEAPU8.copyWithin(dest, src, src + num)
  },
  b(requestedSize: number) {
    requestedSize = requestedSize >>> 0
    abort(`OOM (requestedSize: ${requestedSize})`)
  },
}

function initRuntime(asm: Record<string, () => void>) {
  asm["e"]()
}

const imports = {
  "a": asmLibraryArg,
}

export let malloc: (size: number) => number
export let free: (offset: number, size: number) => void

export let ZSTD_isError: (size: number) => boolean

export let ZSTD_compressBound: (size: number) => number

export let ZSTD_createCCtx: () => number
export let ZSTD_freeCCtx: (cCtx: number) => void

export let ZSTD_createCDict: (dictBuffer: number, dictSize: number, compressionLevel: number) => number
export let ZSTD_freeCDict: (cDict: number) => void

export let ZSTD_compress_usingCDict: (cCtx: number, dst: number, dstCapacity: number, src: number, srcSize: number, cDict: number) => number

// let stackSave
// let stackRestore
// let stackAlloc

export const zstdReady = from(WebAssembly.instantiateStreaming(fetch(zstdWasmUrl), imports)).pipe(
  map(output => {
    const asm = (output.instance || output).exports as Record<string, never>
    malloc = asm["f"]
    free = asm["g"]
    ZSTD_isError = asm["h"]
    ZSTD_compressBound = asm["i"]
    ZSTD_createCCtx = asm["j"]
    ZSTD_freeCCtx = asm["k"]
    ZSTD_freeCDict = asm["l"]
    ZSTD_createCDict = asm["m"]
    ZSTD_compress_usingCDict = asm["n"]
    // stackSave = asm["o"];
    // stackRestore = asm["p"];
    // stackAlloc = asm["q"];
    // wasmTable = asm["n"];
    const wasmMemory = asm["d"] as { buffer: ArrayBuffer }
    buffer = wasmMemory.buffer
    HEAPU8 = new Uint8Array(buffer)
    initRuntime(asm)
    return null
  }),
  shareReplay(1)
)