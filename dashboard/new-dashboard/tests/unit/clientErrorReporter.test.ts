import { describe, expect, it } from "vitest"
import { ClientErrorReporter, installClientErrorReporting } from "../../src/shared/clientErrorReporter"

interface CapturedRequest {
  input: RequestInfo | URL
  init?: RequestInit
}

interface ClientErrorRequestBody {
  errors: {
    source: string
    version: string
    message: string
    count: number
  }[]
}

function parseBody(request: CapturedRequest): ClientErrorRequestBody {
  return JSON.parse(String(request.init?.body)) as ClientErrorRequestBody
}

describe("client error reporter", () => {
  it("deduplicates errors and posts a bounded batch", async () => {
    const requests: CapturedRequest[] = []
    const reporter = new ClientErrorReporter({
      endpoint: "/api/client-errors",
      fetchImpl: (input, init) => {
        requests.push({ input, init })
        return Promise.resolve(new Response("", { status: 202 }))
      },
      now: () => Date.parse("2026-01-01T00:00:00.000Z"),
      flushIntervalMs: 60_000,
    })

    reporter.report({ source: "console_error", message: "boom" })
    reporter.report({ source: "console_error", message: "boom" })

    await reporter.flush()
    reporter.dispose()

    expect(requests).toHaveLength(1)
    expect(requests[0].input).toBe("/api/client-errors")
    expect(parseBody(requests[0]).errors).toStrictEqual([
      expect.objectContaining({
        source: "console_error",
        message: "boom",
        count: 2,
      }),
    ])
  })

  it("stamps each reported event with the configured build version", async () => {
    const requests: CapturedRequest[] = []
    const reporter = new ClientErrorReporter({
      endpoint: "/api/client-errors",
      fetchImpl: (input, init) => {
        requests.push({ input, init })
        return Promise.resolve(new Response("", { status: 202 }))
      },
      now: () => Date.parse("2026-01-01T00:00:00.000Z"),
      flushIntervalMs: 60_000,
      appVersion: "abc1234",
    })

    reporter.report({ source: "window_error", message: "boom" })
    await reporter.flush()
    reporter.dispose()

    expect(requests).toHaveLength(1)
    expect(parseBody(requests[0]).errors[0]).toStrictEqual(
      expect.objectContaining({
        source: "window_error",
        version: "abc1234",
      })
    )
  })

  it("always stamps a non-empty build version even when none is configured", async () => {
    const requests: CapturedRequest[] = []
    const reporter = new ClientErrorReporter({
      endpoint: "/api/client-errors",
      fetchImpl: (input, init) => {
        requests.push({ input, init })
        return Promise.resolve(new Response("", { status: 202 }))
      },
      now: () => Date.parse("2026-01-01T00:00:00.000Z"),
      flushIntervalMs: 60_000,
    })

    reporter.report({ source: "window_error", message: "boom" })
    await reporter.flush()
    reporter.dispose()

    // No appVersion configured: the reporter falls back to the build-time __APP_VERSION__ global
    // (injected by Vite's `define`), or the literal "unknown" in contexts where it is absent.
    const { version } = parseBody(requests[0]).errors[0]
    expect(typeof version).toBe("string")
    expect(version.length).toBeGreaterThan(0)
  })

  it("drops reports after the per-minute client limit is reached", async () => {
    const requests: CapturedRequest[] = []
    let now = Date.parse("2026-01-01T00:00:00.000Z")
    const reporter = new ClientErrorReporter({
      endpoint: "/api/client-errors",
      fetchImpl: (input, init) => {
        requests.push({ input, init })
        return Promise.resolve(new Response("", { status: 202 }))
      },
      now: () => now,
      maxReportsPerMinute: 2,
      flushIntervalMs: 60_000,
    })

    reporter.report({ source: "window_error", message: "first" })
    reporter.report({ source: "window_error", message: "second" })
    reporter.report({ source: "window_error", message: "third" })
    await reporter.flush()

    now += 60_000
    reporter.report({ source: "window_error", message: "fourth" })
    await reporter.flush()
    reporter.dispose()

    expect(requests).toHaveLength(2)
    expect(parseBody(requests[0]).errors.map((event) => event.message)).toStrictEqual(["first", "second"])
    expect(parseBody(requests[1]).errors.map((event) => event.message)).toStrictEqual(["fourth"])
  })

  it("captures console.error without blocking the original call", async () => {
    const requests: CapturedRequest[] = []
    const previousConsoleError = console.error
    const consoleCalls: unknown[][] = []
    console.error = (...args: unknown[]) => {
      consoleCalls.push(args)
    }

    const dispose = installClientErrorReporting({
      endpoint: "/api/client-errors",
      fetchImpl: (input, init) => {
        requests.push({ input, init })
        return Promise.resolve(new Response("", { status: 202 }))
      },
      flushIntervalMs: 0,
    })

    try {
      console.error("installed boom")
      await new Promise((resolve) => {
        setTimeout(resolve, 10)
      })
    } finally {
      dispose()
      console.error = previousConsoleError
    }

    expect(consoleCalls).toStrictEqual([["installed boom"]])
    expect(requests).toHaveLength(1)
    expect(parseBody(requests[0]).errors[0]).toStrictEqual(
      expect.objectContaining({
        source: "console_error",
        message: "installed boom",
        count: 1,
      })
    )
  })
})
