export type YAxisScale = "linear" | "log"

// 1, 2 and 5 in every decade, 0.001 … 5e12. ECharts' own log ticks sit only on whole powers of the base and it
// pads the axis out to them: 3.3–10.3 s would get a 1 s–100 s axis with three labels.
const LOG_TICKS = Array.from({ length: 16 }, (_, index) => index - 3).flatMap((exponent) => [1, 2, 5].map((step) => step * 10 ** exponent))

// Keeps the outermost lines off the chart edges, and the bounds apart when every value is the same.
const PADDING = 1.1

// undefined = no usable bound (no data yet, or nothing positive): leave the axis autoscaled.
export function logTickAtOrBelow(value: number): number | undefined {
  return Number.isFinite(value) && value > 0 ? LOG_TICKS.findLast((tick) => tick <= value) : undefined
}

export function logTickAtOrAbove(value: number): number | undefined {
  return Number.isFinite(value) && value > 0 ? LOG_TICKS.find((tick) => tick >= value) : undefined
}

// A log axis has no zero, so the zero baseline ("Flexible zero") and the lagging floor do not apply to it.
// Not for #mean_value measures: their std-dev band is stacked in linear space.
export function logYAxis() {
  return {
    type: "log" as const,
    min: (extent: { min: number }) => logTickAtOrBelow(extent.min / PADDING) ?? null,
    max: (extent: { max: number }) => logTickAtOrAbove(extent.max * PADDING) ?? null,
    axisLabel: { customValues: LOG_TICKS },
    axisTick: { customValues: LOG_TICKS },
  }
}
