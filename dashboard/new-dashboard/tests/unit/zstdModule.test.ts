import { afterEach, describe, expect, it, vi } from "vitest"
import { zstdReady } from "../../src/components/common/zstd-module"

describe("zstd WASM initialization", () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it("forwards a failed WASM load to the subscriber instead of hanging", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn<() => Promise<Response>>(() => Promise.resolve(new Response()))
    )
    vi.stubGlobal("WebAssembly", {
      ...WebAssembly,
      instantiateStreaming: vi.fn<() => Promise<WebAssembly.WebAssemblyInstantiatedSource>>(() => Promise.reject(new Error("wasm unavailable"))),
    })

    const error = await new Promise<unknown>((resolve, reject) => {
      zstdReady().subscribe({
        next: () => {
          reject(new Error("unexpected value"))
        },
        complete: () => {
          reject(new Error("unexpected complete"))
        },
        error: resolve,
      })
    })
    expect(error).toStrictEqual(new Error("wasm unavailable"))
  })
})
