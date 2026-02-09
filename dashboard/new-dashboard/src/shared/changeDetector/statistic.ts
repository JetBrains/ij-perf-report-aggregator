export function median(arr: number[]): number {
  if (arr.length === 0) {
    throw new Error("Data array is empty")
  }
  const sortedArr = arr.toSorted((a, b) => a - b)
  const mid = Math.floor(sortedArr.length / 2)
  return sortedArr.length % 2 === 0 ? (sortedArr[mid - 1] + sortedArr[mid]) / 2 : sortedArr[mid]
}

export function rollingMad(arr: number[], windowSize: number): { medians: number[]; mads: number[] } {
  const medians = []
  const mads = []

  for (let i = 0; i < arr.length; i++) {
    const start = Math.max(0, i - Math.floor(windowSize / 2))
    const end = Math.min(arr.length, i + Math.ceil(windowSize / 2))
    const window = arr.slice(start, end)

    const med = median(window)
    const mad = median(window.map((val) => Math.abs(val - med)))

    medians.push(med)
    mads.push(mad)
  }

  return { medians, mads }
}
