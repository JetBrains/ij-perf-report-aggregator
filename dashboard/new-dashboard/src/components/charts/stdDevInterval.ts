import { LineSeriesOption } from "echarts/charts"

const MEAN_SUFFIX = "#mean_value"
const STANDARD_DEVIATION_SUFFIX = "#standard_deviation"
export const STD_DEV_BAND_OPACITY = 0.15

export function getStdDevMeasureName(measure: string): string | null {
  return measure.endsWith(MEAN_SUFFIX) ? measure.slice(0, -MEAN_SUFFIX.length) + STANDARD_DEVIATION_SUFFIX : null
}

export interface StdDevBand {
  readonly timestamps: number[]
  readonly lowerEdge: number[]
  readonly width: number[]
}

interface StdDevBandRequest {
  readonly timestamps: readonly number[]
  readonly values: readonly number[]
  readonly stdDevs: readonly number[]
  readonly valueScale: number
}

// Values and deviations are paired before filtering; valueScale puts the spread on the same axis.
export function buildStdDevBand({ timestamps, values, stdDevs, valueScale }: StdDevBandRequest): StdDevBand | null {
  const band: StdDevBand = { timestamps: [], lowerEdge: [], width: [] }
  let hasWidth = false
  for (const [index, timestamp] of timestamps.entries()) {
    const value = values[index]
    if (!Number.isFinite(value)) continue
    const stdDev = (stdDevs[index] ?? 0) * valueScale
    // Missing deviations pinch to the mean; duration bands cannot extend below zero.
    const lower = Math.max(0, value - stdDev)
    band.timestamps.push(timestamp)
    band.lowerEdge.push(lower)
    band.width.push(value + stdDev - lower)
    hasWidth ||= stdDev > 0
  }
  return hasWidth ? band : null
}

// An invisible lower line lifts the filled width series off the axis.
export function getStdDevBandSeries(band: StdDevBand, ownerSeriesId: string, seriesName: string): LineSeriesOption[] {
  const common: LineSeriesOption = {
    name: seriesName, // Share the owner's palette color and legend entry.
    type: "line",
    stack: "stdDevBand:" + ownerSeriesId,
    stackStrategy: "all", // Also stack onto a lower edge of zero.
    symbol: "none",
    silent: true,
    tooltip: { show: false },
    legendHoverLink: false,
    animation: false,
    z: 1,
    lineStyle: { opacity: 0 },
    encode: { x: 0, y: 1 },
  }
  return [
    {
      ...common,
      id: ownerSeriesId + "stdDevBandLowerEdge",
      data: band.timestamps.map((time, i) => [time, band.lowerEdge[i]]),
    },
    {
      ...common,
      id: ownerSeriesId + "stdDevBandFill",
      data: band.timestamps.map((time, i) => [time, band.width[i]]),
      areaStyle: { opacity: STD_DEV_BAND_OPACITY },
    },
  ]
}
