export interface DebouncedFunction {
  (): void

  clear: () => void
}

export function debounceSync(func: () => void, wait: number = 100): DebouncedFunction {
  let timeout = 0
  let isScheduled = false
  let timestamp = 0
  function later() {
    const last = Date.now() - timestamp
    if (last < wait && last >= 0) {
      timeout = window.setTimeout(later, wait - last)
    }
    else {
      isScheduled = false
      func()
    }
  }

  const debounced: DebouncedFunction = function () {
    timestamp = Date.now()
    if (!isScheduled) {
      isScheduled = true
      timeout = window.setTimeout(later, wait)
    }
  }
  debounced.clear = function () {
    if (isScheduled) {
      isScheduled = false
      clearTimeout(timeout)
    }
  }
  return debounced
}