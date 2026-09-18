import { LineSeriesOption } from "echarts/charts"
import { median } from "../../shared/changeDetector/statistic"

// A series that stops reporting mid-chart is otherwise indistinguishable from one that merely runs behind its
// neighbours, so it gets a line in the series color through the whole plot at the point its data ends: inside a
// dense band of lines a marker on the data itself competes with the data for the same pixels.

const DAY_MS = 24 * 60 * 60 * 1000

// Whatever the cadence, a gap has to outlast a long weekend before it means the data stopped.
const MIN_STALE_GAP_MS = 3 * DAY_MS

// ...and it has to be long by the standards of the series itself: a nightly series is stale after a
// few days, a weekly one is not.
const STALE_GAP_TO_INTERVAL_RATIO = 3

function getStaleGapThreshold(timestamps: readonly number[]): number {
  // A build that reported a measurement more than once leaves repeated timestamps in the series, and the
  // zero-length intervals between them would pull the median down to a cadence the series never ran at -
  // far enough, with enough repeats, to call a healthy weekly series stale days before its next run.
  const intervals = timestamps
    .slice(1)
    .map((timestamp, index) => timestamp - timestamps[index])
    .filter((interval) => interval > 0)
  return intervals.length === 0 ? MIN_STALE_GAP_MS : Math.max(MIN_STALE_GAP_MS, STALE_GAP_TO_INTERVAL_RATIO * median(intervals))
}

/**
 * The reference a series is judged stale against - the chart's own last point rather than the current time:
 * when every series stops at once (a retired test, or a time range that ends in the past) the whole chart is
 * flat and there is nothing to single out. It also means a chart of a single series never marks it.
 */
export function getChartLastTimestamp(timestampsPerSeries: readonly (readonly number[])[]): number {
  let chartLastTimestamp = Number.NEGATIVE_INFINITY
  for (const timestamps of timestampsPerSeries) {
    const lastTimestamp = timestamps.at(-1)
    if (lastTimestamp != undefined && lastTimestamp > chartLastTimestamp) {
      chartLastTimestamp = lastTimestamp
    }
  }
  return chartLastTimestamp
}

/**
 * The mark line for a series that stopped reporting before the rest of the chart did, or `undefined` while it is
 * keeping up. Timestamps are expected in ascending order, as `mergeSeries` leaves them.
 */
export function getStaleSeriesMarkLine(timestamps: readonly number[], chartLastTimestamp: number): LineSeriesOption["markLine"] {
  const lastTimestamp = timestamps.at(-1)
  if (lastTimestamp == undefined) {
    return undefined
  }
  const gap = chartLastTimestamp - lastTimestamp
  // The threshold can never drop below the floor, so the cadence is only worth computing once the gap clears it.
  if (gap <= MIN_STALE_GAP_MS || gap <= getStaleGapThreshold(timestamps)) {
    return undefined
  }
  return {
    // Unset lineStyle.color makes the line take the color of the series it belongs to; thin and dashed keeps it
    // from reading as data.
    lineStyle: { type: "dashed", width: 1, opacity: 0.6 },
    symbol: "none",
    label: { show: false },
    silent: true,
    animation: false,
    data: [{ xAxis: lastTimestamp }],
  }
}
