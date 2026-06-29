export type ClientErrorSource = "window_error" | "unhandled_rejection" | "console_error" | "resource_error" | "vue_error"

// Build version injected by Vite (`define`) at build time. It is absent under unit tests and any
// non-Vite context, where `typeof` on an undeclared global is safe and yields "undefined".
declare const __APP_VERSION__: string
const buildVersion = typeof __APP_VERSION__ === "string" ? __APP_VERSION__ : "unknown"

export interface ClientErrorReporterOptions {
  endpoint: string | (() => string)
  fetchImpl?: typeof fetch
  now?: () => number
  flushIntervalMs?: number
  maxBatchSize?: number
  maxBufferedEvents?: number
  maxReportsPerMinute?: number
  appVersion?: string
}

export interface ClientErrorReportInput {
  source: ClientErrorSource
  message?: unknown
  error?: unknown
  url?: string
  line?: number
  column?: number
}

interface ClientErrorPayload {
  source: ClientErrorSource
  version: string
  message: string
  stack?: string
  url?: string
  pageUrl: string
  userAgent: string
  line?: number
  column?: number
  count: number
  firstSeen: string
  lastSeen: string
}

const maxMessageLength = 1000
const maxStackLength = 4000
const maxURLLength = 1000
const maxUserAgentLength = 300
const defaultFlushIntervalMs = 2000
const defaultMaxBatchSize = 10
const defaultMaxBufferedEvents = 50
const defaultMaxReportsPerMinute = 20
const rateLimitWindowMs = 60_000

let installedDispose: (() => void) | null = null
let activeReporter: ClientErrorReporter | null = null

/**
 * Reports an error through the currently installed reporter, if any. Used to forward errors from
 * sources that aren't global browser events (e.g. Vue's `app.config.errorHandler`). No-op when the
 * reporter hasn't been installed.
 */
export function reportClientError(input: ClientErrorReportInput): void {
  activeReporter?.report(input)
}

export class ClientErrorReporter {
  private readonly endpointProvider: () => string
  private readonly fetchImpl: typeof fetch
  private readonly now: () => number
  private readonly flushIntervalMs: number
  private readonly maxBatchSize: number
  private readonly maxBufferedEvents: number
  private readonly maxReportsPerMinute: number
  private readonly appVersion: string
  private readonly events = new Map<string, ClientErrorPayload>()
  private flushTimer: ReturnType<typeof setTimeout> | null = null
  private inFlight = false
  private rateLimitWindowStartedAt: number
  private acceptedReportsInWindow = 0

  constructor(options: ClientErrorReporterOptions) {
    if (typeof options.endpoint === "function") {
      this.endpointProvider = options.endpoint
    } else {
      const endpoint = options.endpoint
      this.endpointProvider = () => endpoint
    }
    this.fetchImpl = options.fetchImpl ?? ((input, init) => fetch(input, init))
    this.now = options.now ?? Date.now
    this.flushIntervalMs = options.flushIntervalMs ?? defaultFlushIntervalMs
    this.maxBatchSize = options.maxBatchSize ?? defaultMaxBatchSize
    this.maxBufferedEvents = options.maxBufferedEvents ?? defaultMaxBufferedEvents
    this.maxReportsPerMinute = options.maxReportsPerMinute ?? defaultMaxReportsPerMinute
    this.appVersion = options.appVersion ?? buildVersion
    this.rateLimitWindowStartedAt = this.now()
  }

  report(input: ClientErrorReportInput): void {
    if (this.endpointProvider() === "") {
      return
    }

    const now = this.now()
    const event = createClientErrorPayload(input, now, this.appVersion)
    const fingerprint = createFingerprint(event)
    const existing = this.events.get(fingerprint)
    if (existing != null) {
      if (!this.consumeRateLimitToken(now)) {
        return
      }
      existing.count += 1
      existing.lastSeen = event.lastSeen
    } else {
      if (this.events.size >= this.maxBufferedEvents) {
        return
      }
      if (!this.consumeRateLimitToken(now)) {
        return
      }
      this.events.set(fingerprint, event)
    }

    if (this.events.size >= this.maxBatchSize) {
      void this.flush()
    } else {
      this.scheduleFlush()
    }
  }

  async flush(): Promise<void> {
    if (this.flushTimer != null) {
      clearTimeout(this.flushTimer)
      this.flushTimer = null
    }
    if (this.inFlight || this.events.size === 0) {
      return
    }

    const endpoint = this.endpointProvider()
    if (endpoint === "") {
      return
    }

    const entries = [...this.events.entries()].slice(0, this.maxBatchSize)
    for (const [fingerprint] of entries) {
      this.events.delete(fingerprint)
    }

    const body = JSON.stringify({ errors: entries.map(([, event]) => event) })
    // keepalive lets the request survive a page unload, but browsers cap the total keepalive
    // body at 64 KiB measured in bytes, not characters — so size it with the encoded byte length.
    const bodyByteLength = new TextEncoder().encode(body).length
    this.inFlight = true
    try {
      await this.fetchImpl(endpoint, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body,
        keepalive: bodyByteLength <= 60 * 1024,
      })
    } catch {
      // Reporting is best-effort; never create more console noise from the reporter itself.
    } finally {
      this.inFlight = false
      if (this.events.size > 0) {
        this.scheduleFlush(0)
      }
    }
  }

  dispose(): void {
    if (this.flushTimer != null) {
      clearTimeout(this.flushTimer)
      this.flushTimer = null
    }
    this.events.clear()
  }

  private consumeRateLimitToken(now: number): boolean {
    if (now - this.rateLimitWindowStartedAt >= rateLimitWindowMs) {
      this.rateLimitWindowStartedAt = now
      this.acceptedReportsInWindow = 0
    }
    if (this.acceptedReportsInWindow >= this.maxReportsPerMinute) {
      return false
    }
    this.acceptedReportsInWindow += 1
    return true
  }

  private scheduleFlush(delayMs = this.flushIntervalMs): void {
    if (this.flushTimer != null) {
      return
    }
    this.flushTimer = setTimeout(() => {
      this.flushTimer = null
      void this.flush()
    }, delayMs)
  }
}

export function installClientErrorReporting(options: ClientErrorReporterOptions): () => void {
  installedDispose?.()

  const reporter = new ClientErrorReporter(options)
  activeReporter = reporter
  const onError = (event: Event) => {
    if (isErrorEvent(event)) {
      reporter.report({
        source: "window_error",
        message: event.message,
        error: event.error,
        url: event.filename,
        line: event.lineno,
        column: event.colno,
      })
    } else {
      const target = event.target
      reporter.report({
        source: "resource_error",
        message: resourceErrorMessage(target),
        url: resourceURL(target),
      })
    }
  }

  const onUnhandledRejection = (event: PromiseRejectionEvent) => {
    reporter.report({
      source: "unhandled_rejection",
      message: messageFromUnknown(event.reason),
      error: event.reason,
    })
  }

  const originalConsoleError = console.error
  console.error = (...args: unknown[]) => {
    try {
      reporter.report({
        source: "console_error",
        message: formatConsoleArguments(args),
        error: args.find((arg) => arg instanceof Error),
      })
    } catch {
      // Reporting is best-effort and must never interfere with the original console.error call.
    }
    originalConsoleError(...args)
  }

  const flush = () => {
    void reporter.flush()
  }
  const onVisibilityChange = () => {
    if (document.visibilityState === "hidden") {
      flush()
    }
  }

  window.addEventListener("error", onError, true)
  window.addEventListener("unhandledrejection", onUnhandledRejection)
  window.addEventListener("pagehide", flush)
  document.addEventListener("visibilitychange", onVisibilityChange)

  const dispose = () => {
    window.removeEventListener("error", onError, true)
    window.removeEventListener("unhandledrejection", onUnhandledRejection)
    window.removeEventListener("pagehide", flush)
    document.removeEventListener("visibilitychange", onVisibilityChange)
    console.error = originalConsoleError
    reporter.dispose()
    if (activeReporter === reporter) {
      activeReporter = null
    }
    if (installedDispose === dispose) {
      installedDispose = null
    }
  }
  installedDispose = dispose
  return dispose
}

function createClientErrorPayload(input: ClientErrorReportInput, now: number, version: string): ClientErrorPayload {
  const stack = stackFromUnknown(input.error)
  const timestamp = new Date(now).toISOString()
  return {
    source: input.source,
    version,
    message: truncate(messageFromUnknown(input.message ?? input.error), maxMessageLength),
    stack: stack == null ? undefined : truncate(stack, maxStackLength),
    url: input.url == null ? undefined : truncate(input.url, maxURLLength),
    pageUrl: truncate(window.location.href, maxURLLength),
    userAgent: truncate(navigator.userAgent, maxUserAgentLength),
    line: normalizePositiveInteger(input.line),
    column: normalizePositiveInteger(input.column),
    count: 1,
    firstSeen: timestamp,
    lastSeen: timestamp,
  }
}

function createFingerprint(event: ClientErrorPayload): string {
  return [event.source, event.message, event.stack?.split("\n", 1)[0] ?? "", event.url ?? "", event.line ?? "", event.column ?? ""].join("\n")
}

function isErrorEvent(event: Event): event is ErrorEvent {
  return "message" in event && "filename" in event
}

export function messageFromUnknown(value: unknown): string {
  if (value instanceof Error) {
    return value.message === "" ? value.name : value.message
  }
  if (typeof value === "string") {
    return value
  }
  if (value == null) {
    return "Unknown client error"
  }
  if (typeof value === "number" || typeof value === "boolean" || typeof value === "bigint") {
    return String(value)
  }
  try {
    const json = JSON.stringify(value)
    return json == null ? String(value) : json
  } catch {
    return String(value)
  }
}

function stackFromUnknown(value: unknown): string | undefined {
  if (value instanceof Error && typeof value.stack === "string" && value.stack !== "") {
    return value.stack
  }
  return undefined
}

function formatConsoleArguments(args: unknown[]): string {
  return args.map((arg) => messageFromUnknown(arg)).join(" ")
}

function resourceErrorMessage(target: EventTarget | null): string {
  if (target instanceof HTMLElement) {
    return `Resource failed to load: ${target.tagName.toLowerCase()}`
  }
  return "Resource failed to load"
}

function resourceURL(target: EventTarget | null): string | undefined {
  if (target instanceof HTMLScriptElement) {
    return target.src
  }
  if (target instanceof HTMLLinkElement) {
    return target.href
  }
  if (target instanceof HTMLImageElement) {
    return target.src
  }
  return undefined
}

function normalizePositiveInteger(value: number | undefined): number | undefined {
  if (value == null || !Number.isFinite(value) || value <= 0) {
    return undefined
  }
  return Math.floor(value)
}

function truncate(value: string, maxLength: number): string {
  if (value.length <= maxLength) {
    return value
  }
  let end = maxLength
  // Don't cut in the middle of a surrogate pair, which would leave a dangling high surrogate.
  // charCodeAt (not codePointAt) is intentional: we need the single UTF-16 code unit at the
  // boundary, whereas codePointAt would merge the pair and hide the split.
  // oxlint-disable-next-line unicorn/prefer-code-point
  const lastCode = value.charCodeAt(end - 1)
  if (lastCode >= 0xd800 && lastCode <= 0xdbff) {
    end -= 1
  }
  return value.slice(0, end)
}
