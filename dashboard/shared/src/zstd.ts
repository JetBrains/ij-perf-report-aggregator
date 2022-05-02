import { from, shareReplay } from "rxjs"
import { _free, _malloc, _ZSTD_compress, _ZSTD_compressBound, _ZSTD_isError, HEAPU8, zstdReady } from "./zstd-module"

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function initZstd(): Promise<void> {
  return zstdReady
}

export const initZstdObservable = from(initZstd()).pipe(shareReplay(1))

function isError(code: number): boolean {
  return _ZSTD_isError(code)
}

function compressBound(size: number): number {
  return _ZSTD_compressBound(size)
}

const zstdCompressionLevel = 3

// https://github.com/bokuweb/zstd-wasm
export function compressZstdToUrlSafeBase64(s: string) {
  // see https://developer.mozilla.org/en-US/docs/Web/API/TextEncoder/encodeInto#buffer_sizing about computing the output space needed for full conversion of string to bytes
  const maxUncompressedSize = s.length * 3
  // compute maximum compressed size in worst case single-pass scenario - https://zstd.docsforge.com/dev/api/ZSTD_compressBound/
  const uncompressedOffset = _malloc(maxUncompressedSize)

  try {
    const heapU8 = HEAPU8
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    const sourceSize = new TextEncoder().encodeInto(s, heapU8.subarray(uncompressedOffset, uncompressedOffset + maxUncompressedSize)).written!

    const maxCompressedSize = compressBound(sourceSize)
    const compressedOffset = _malloc(maxCompressedSize)
    try {
      // compress - https://zstd.docsforge.com/dev/api/ZSTD_compress/
      // size_t ZSTD_compress(void *dst, size_t dstCapacity, const void *src, size_t srcSize, int compressionLevel)
      // console.time("zstd")
      const sizeOrError = _ZSTD_compress(compressedOffset, maxCompressedSize, uncompressedOffset, sourceSize, zstdCompressionLevel)
      if (isError(sizeOrError)) {
        // noinspection ExceptionCaughtLocallyJS
        throw new Error(`Failed to compress with code ${sizeOrError}`)
      }

      const result = bytesToBase64(heapU8, compressedOffset, sizeOrError)
      _free(compressedOffset, maxCompressedSize)
      _free(uncompressedOffset, maxUncompressedSize)
      // console.timeEnd("zstd")
      return result
    }
    finally {
      _free(compressedOffset, maxCompressedSize)
    }
  }
  catch (e) {
    _free(uncompressedOffset, maxUncompressedSize)
    throw e
  }
}

const base64UrlSafe = [
  "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
  "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
  "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
  "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
  "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "-", "_",
]

function bytesToBase64(bytes: Uint8Array, offset: number, size: number) {
  let result = ""
  let i = offset + 2
  const end = offset + size
  for (; i < end; i += 3) {
    result += base64UrlSafe[bytes[i - 2] >> 2]
    result += base64UrlSafe[((bytes[i - 2] & 0x03) << 4) | (bytes[i - 1] >> 4)]
    result += base64UrlSafe[((bytes[i - 1] & 0x0F) << 2) | (bytes[i] >> 6)]
    result += base64UrlSafe[bytes[i] & 0x3F]
  }
  // 1 octet yet to write
  if (i === end + 1) {
    result += base64UrlSafe[bytes[i - 2] >> 2]
    result += base64UrlSafe[(bytes[i - 2] & 0x03) << 4]
  }
  // 2 octets yet to write
  if (i === end) {
    result += base64UrlSafe[bytes[i - 2] >> 2]
    result += base64UrlSafe[((bytes[i - 2] & 0x03) << 4) | (bytes[i - 1] >> 4)]
    result += base64UrlSafe[(bytes[i - 1] & 0x0F) << 2]
  }
  return result
}