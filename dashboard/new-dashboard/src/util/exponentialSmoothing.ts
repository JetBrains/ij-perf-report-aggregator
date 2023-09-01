import { DataTransformOption, ExternalDataTransformResultItem, ExternalSource } from "echarts/types/src/data/helper/transform"

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

export interface ExternalDataTransformParam<TO extends DataTransformOption = DataTransformOption> {
  upstream: ExternalSource
  upstreamList: ExternalSource[]
  config: TO["config"] | undefined
}

export const exponentialSmoothingTransform = {
  type: "ijPerf:exponentialSmoothing",

  transform: function transform(params: ExternalDataTransformParam): ExternalDataTransformResultItem | ExternalDataTransformResultItem[] {
    const upstream = params.upstream

    const rawData = upstream.cloneRawData() as number[][] | string[][]
    const yValues = rawData[1] as number[]

    // const smoothedValues = exponentialSmoothing(yValues, alpha)
    const smoothedValues = exponentialSmoothingWithAlphaInference(yValues)
    // const smoothedValues = doubleExponentialSmoothing(yValues, alpha, 0.2)

    const points = [rawData[0], smoothedValues]
    return [
      {
        data: points,
      },
    ]
  },
}
