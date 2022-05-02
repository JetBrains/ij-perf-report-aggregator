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

export let _malloc: (size: number) => number
export let _free: (offset: number, size: number) => void
export let _ZSTD_isError: (size: number) => boolean
export let _ZSTD_compressBound: (size: number) => number
// export let _ZSTD_compress_usingDict
export let _ZSTD_compress: (dst: number, dstCapacity: number, src: number, srcSize: number, compressionLevel: number) => number
// let stackSave
// let stackRestore
// let stackAlloc

export const zstdReady = WebAssembly.instantiateStreaming(fetch(zstdWasmUrl), imports).then(function (output) {
  const asm = (output.instance || output).exports as Record<string, never>
  _malloc = asm["f"]
  _free = asm["g"]
  _ZSTD_isError = asm["h"]
  _ZSTD_compressBound = asm["i"]
  // _ZSTD_compress_usingDict = asm["j"]
  _ZSTD_compress = asm["k"]
  buffer = (asm["d"] as { buffer: ArrayBuffer }).buffer
  HEAPU8 = new Uint8Array(buffer)
  initRuntime(asm)
})