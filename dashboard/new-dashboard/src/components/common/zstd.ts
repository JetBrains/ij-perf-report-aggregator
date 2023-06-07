import { forkJoin, map, shareReplay } from "rxjs"
import zstdDictionaryUrl from "../../../../../pkg/data-query/zstd.dictionary?url"
import { fromFetchWithRetryAndErrorHandling } from "../../configurators/rxjs"
import {
  free,
  malloc,
  ZSTD_compress_usingCDict,
  ZSTD_compressBound,
  ZSTD_createCCtx,
  ZSTD_createCDict,
  ZSTD_freeCCtx,
  ZSTD_freeCDict,
  ZSTD_isError,
  HEAPU8,
  zstdReady,
} from "./zstd-module"

let compressor: CompressorUsingDictionary | null = null

export function getCompressor(): CompressorUsingDictionary {
  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  return compressor!
}

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-expect-error
if (import.meta.hot) {
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-expect-error
  import.meta.hot.dispose(() => {
    if (compressor !== null) {
      compressor.dispose()
      compressor = null
    }
  })
}

// zstdDictionaryUrl - if zstd dictionary will be changed on a server side, then server must introduce a new endpoint for query (currently, `/api/q`)
export const initZstdObservable = forkJoin([fromFetchWithRetryAndErrorHandling<ArrayBuffer>(zstdDictionaryUrl, (it) => it.arrayBuffer()), zstdReady]).pipe(
  map(([dictionaryData, _]) => {
    if (compressor !== null) {
      compressor.dispose()
    }
    compressor = new CompressorUsingDictionary(dictionaryData)
    return null
  }),
  shareReplay(1)
)

function isError(code: number): boolean {
  return ZSTD_isError(code)
}

function compressBound(size: number): number {
  return ZSTD_compressBound(size)
}

const zstdCompressionLevel = 7

export class CompressorUsingDictionary {
  private readonly zstdDictionaryPointer: number
  private readonly zstdDictionarySize: number
  private readonly dict: number
  private readonly context: number

  constructor(dictionaryData: ArrayBuffer) {
    this.zstdDictionaryPointer = malloc(dictionaryData.byteLength)
    this.zstdDictionarySize = dictionaryData.byteLength
    HEAPU8.set(new Uint8Array(dictionaryData), this.zstdDictionaryPointer)
    this.dict = ZSTD_createCDict(this.zstdDictionaryPointer, this.zstdDictionarySize, zstdCompressionLevel)

    this.context = ZSTD_createCCtx()
  }

  compress(s: string): string {
    // see https://developer.mozilla.org/en-US/docs/Web/API/TextEncoder/encodeInto#buffer_sizing about computing the output space needed for full conversion of string to bytes
    const maxUncompressedSize = s.length * 3
    // compute maximum compressed size in worst case single-pass scenario - https://zstd.docsforge.com/dev/api/ZSTD_compressBound/
    const uncompressedOffset = malloc(maxUncompressedSize)
    try {
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      const sourceSize = new TextEncoder().encodeInto(s, HEAPU8.subarray(uncompressedOffset, uncompressedOffset + maxUncompressedSize)).written!

      const maxCompressedSize = compressBound(sourceSize)
      const compressedOffset = malloc(maxCompressedSize)
      try {
        // compress - https://zstd.docsforge.com/dev/api/ZSTD_compress_usingCDict/
        // size_t ZSTD_compress(void *dst, size_t dstCapacity, const void *src, size_t srcSize, int compressionLevel)
        // console.time("zstd")
        const sizeOrError = ZSTD_compress_usingCDict(this.context, compressedOffset, maxCompressedSize, uncompressedOffset, sourceSize, this.dict)
        if (isError(sizeOrError)) {
          // noinspection ExceptionCaughtLocallyJS
          throw new Error(`Failed to compress with code ${sizeOrError}`)
        }
        // console.timeEnd("zstd")
        return bytesToBase64(HEAPU8, compressedOffset, sizeOrError)
      } finally {
        free(compressedOffset, maxCompressedSize)
      }
    } finally {
      free(uncompressedOffset, maxUncompressedSize)
    }
  }

  dispose() {
    ZSTD_freeCDict(this.dict)
    ZSTD_freeCCtx(this.context)
    free(this.zstdDictionaryPointer, this.zstdDictionarySize)
  }
}

const base64UrlSafe = [
  "A",
  "B",
  "C",
  "D",
  "E",
  "F",
  "G",
  "H",
  "I",
  "J",
  "K",
  "L",
  "M",
  "N",
  "O",
  "P",
  "Q",
  "R",
  "S",
  "T",
  "U",
  "V",
  "W",
  "X",
  "Y",
  "Z",
  "a",
  "b",
  "c",
  "d",
  "e",
  "f",
  "g",
  "h",
  "i",
  "j",
  "k",
  "l",
  "m",
  "n",
  "o",
  "p",
  "q",
  "r",
  "s",
  "t",
  "u",
  "v",
  "w",
  "x",
  "y",
  "z",
  "0",
  "1",
  "2",
  "3",
  "4",
  "5",
  "6",
  "7",
  "8",
  "9",
  "-",
  "_",
]

function bytesToBase64(bytes: Uint8Array, offset: number, size: number) {
  let result = ""
  let i = offset + 2
  const end = offset + size
  for (; i < end; i += 3) {
    result += base64UrlSafe[bytes[i - 2] >> 2]
    result += base64UrlSafe[((bytes[i - 2] & 0x03) << 4) | (bytes[i - 1] >> 4)]
    result += base64UrlSafe[((bytes[i - 1] & 0x0f) << 2) | (bytes[i] >> 6)]
    result += base64UrlSafe[bytes[i] & 0x3f]
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
    result += base64UrlSafe[(bytes[i - 1] & 0x0f) << 2]
  }
  return result
}
