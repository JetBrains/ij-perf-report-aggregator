// The lagging metrics are inherently unstable: the same scenario legitimately reports 100 ms on one run
// and 500 ms on the next. An autoscaled axis re-frames every such jump to the full chart height, so
// ordinary noise reads as a dramatic regression. A fixed frame keeps those jumps in proportion.
const LAGGING_DURATION_TOP_MS = 3_000
const LAGGING_SHARE_TOP_PERCENT = 10

// Lagging measures that are durations in milliseconds. "#count" (a number of lags) and
// "#percentage_share" are on completely different scales, so they get their own frame or none at all -
// including on a chart that mixes them with the durations.
const LAGGING_DURATION_MEASURES: ReadonlySet<string> = new Set(["ui.lagging#average", "ui.lagging#max", "ui.lagging#max_value", "ui.lagging#sum"])

const LAGGING_SHARE_MEASURE = "ui.lagging#percentage_share"

// The lower bound for the chart's y-axis maximum, or undefined to leave the axis fully autoscaled. It is
// a floor rather than a limit: the caller still grows the axis past it, so a genuine multi-second lag or
// an unusually high share stays visible instead of being clipped off the top.
export function getMinYAxisTop(measures: string[]): number | undefined {
  if (measures.length === 0) {
    return undefined
  }
  if (measures.every((measure) => LAGGING_DURATION_MEASURES.has(measure))) {
    return LAGGING_DURATION_TOP_MS
  }
  return measures.every((measure) => measure === LAGGING_SHARE_MEASURE) ? LAGGING_SHARE_TOP_PERCENT : undefined
}
