class LocalStorageMock implements Storage {
  private readonly data = new Map<string, string>()

  getItem(key: string): string | null {
    return this.data.get(key) ?? null
  }

  setItem(key: string, value: string): void {
    this.data.set(key, value)
  }

  removeItem(key: string): void {
    this.data.delete(key)
  }

  clear(): void {
    this.data.clear()
  }

  get length(): number {
    return this.data.size
  }

  key(index: number): string | null {
    const keys = [...this.data.keys()]
    return keys[index] ?? null
  }
}

Object.defineProperty(globalThis, "localStorage", {
  value: new LocalStorageMock(),
  writable: true,
  configurable: true,
})

globalThis.fetch = () => Promise.resolve(new Response("", { status: 200 }))

// happy-dom exposes AbortController as a per-window subclass, so `abort` lives on the parent prototype.
// @rxjs/observable-polyfill refuses to initialize unless it is an own property of AbortController.prototype.
if (Object.getOwnPropertyDescriptor(AbortController.prototype, "abort") == null) {
  const inherited = Object.getOwnPropertyDescriptor(Object.getPrototypeOf(AbortController.prototype) as object, "abort")
  if (inherited == null) {
    throw new Error("AbortController.prototype.abort not found on the prototype chain")
  }
  Object.defineProperty(AbortController.prototype, "abort", inherited)
}
