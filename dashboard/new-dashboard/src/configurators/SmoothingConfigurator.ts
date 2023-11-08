import { Observable } from "rxjs"
import { ref, watch } from "vue"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../components/common/dataQuery"
import { useSettingsStore } from "../components/settings/settingsStore"
import { FilterConfigurator } from "./filter"
import { refToObservable } from "./rxjs"

export class SmoothingConfigurator implements DataQueryConfigurator, FilterConfigurator {
  private settingsStore = useSettingsStore()
  readonly value = ref(this.settingsStore.smoothing)

  constructor() {
    watch(
      () => this.settingsStore.smoothing,
      (newValue) => {
        this.value.value = newValue
      }
    )
  }

  createObservable(): Observable<unknown> {
    return refToObservable(this.value)
  }

  configureFilter(_: DataQuery): boolean {
    return true
  }

  configureQuery(_: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
    return true
  }
}

// function exponentialSmoothing(data: number[], alpha: number) {
//   if (data.length === 0) return []
//
//   const smoothed = [data[0]] // initialize with the first data point
//   for (let t = 1; t < data.length; t++) {
//     smoothed[t] = alpha * data[t] + (1 - alpha) * smoothed[t - 1]
//   }
//   return smoothed
// }

// Calculate moving window variability
function movingWindowVariability(data: number[], windowSize: number = 5): number {
  let totalVariability = 0

  for (let i = windowSize; i < data.length; i++) {
    let sumDiff = 0
    for (let j = 0; j < windowSize; j++) {
      sumDiff += Math.abs(data[i - j] - data[i - j - 1])
    }
    totalVariability += sumDiff / windowSize
  }

  return totalVariability / (data.length - windowSize)
}

export function exponentialSmoothingWithAlphaInference(data: number[]): number[] {
  if (data.length === 0) return []

  // Compute variability
  const variability = movingWindowVariability(data)

  // Choose alpha based on variability relative to mean value of data
  const meanValue = data.reduce((acc, val) => acc + val, 0) / data.length
  const relativeVariability = variability / meanValue
  // Choose alpha
  const threshold = 0.1 // You might need to adjust this threshold
  const bestAlpha = relativeVariability > threshold ? 0.1 : 0.5

  // Use the chosen alpha to smooth the dataset
  const smoothed = [data[0]]
  for (let t = 1; t < data.length; t++) {
    smoothed[t] = bestAlpha * data[t] + (1 - bestAlpha) * smoothed[t - 1]
  }

  return smoothed
}
