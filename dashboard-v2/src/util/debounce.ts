export function debounce(func: () => void, wait = 100): () => void {
  let timeout = 0
  let isScheduled = false
  let timestamp = 0
  function later() {
    const last = Date.now() - timestamp
    if (last < wait && last >= 0) {
      timeout = setTimeout(later, wait - last)
    }
    else {
      isScheduled = false
      func()
    }
  }

  const debounced = function () {
    timestamp = Date.now()
    if (!isScheduled) {
      isScheduled = true
      timeout = setTimeout(later, wait)
    }
  }
  debounced.clear = () => {
    if (isScheduled) {
      isScheduled = false
      clearTimeout(timeout)
    }
  }
  return debounced
}
