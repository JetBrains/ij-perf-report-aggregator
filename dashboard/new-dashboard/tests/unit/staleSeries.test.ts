import { describe, expect, it } from "vitest"
import { getChartLastTimestamp, getStaleSeriesMarkLine } from "../../src/components/charts/staleSeries"
import { removeOutliers } from "../../src/components/settings/transforms/outliers"

const DAY_MS = 24 * 60 * 60 * 1000
const START = Date.UTC(2026, 0, 1)

function series(pointCount: number, intervalMs: number, endsAgoMs = 0): number[] {
  const end = START + 100 * DAY_MS - endsAgoMs
  return Array.from({ length: pointCount }, (_, index) => end - (pointCount - 1 - index) * intervalMs)
}

// The chart-wide view the configurator builds one series at a time.
function findStaleSeries(timestampsPerSeries: number[][]): Set<number> {
  const chartLastTimestamp = getChartLastTimestamp(timestampsPerSeries)
  const staleSeries = new Set<number>()
  for (const [index, timestamps] of timestampsPerSeries.entries()) {
    if (getStaleSeriesMarkLine(timestamps, chartLastTimestamp) != undefined) {
      staleSeries.add(index)
    }
  }
  return staleSeries
}

// mergeSeries keeps every measurement, so a build that reported twice shows up as a repeated timestamp.
const withRepeatedMeasurements = (timestamps: number[]) => timestamps.flatMap((timestamp) => [timestamp, timestamp])

const nightly = (endsAgoMs = 0) => series(30, DAY_MS, endsAgoMs)
const weekly = (endsAgoMs = 0) => series(10, 7 * DAY_MS, endsAgoMs)
const perCommit = (endsAgoMs = 0) => series(200, 30 * 60 * 1000, endsAgoMs)

describe("stale series detection", () => {
  it("marks a nightly series that stopped while its neighbour kept reporting", () => {
    expect(findStaleSeries([nightly(), nightly(10 * DAY_MS)])).toStrictEqual(new Set([1]))
  })

  it("marks every series that stopped, and none of those still reporting", () => {
    expect(findStaleSeries([nightly(20 * DAY_MS), nightly(), nightly(10 * DAY_MS), nightly()])).toStrictEqual(new Set([0, 2]))
  })

  it("marks nothing when all the series stop together, however long ago that was", () => {
    expect(findStaleSeries([nightly(90 * DAY_MS), nightly(90 * DAY_MS)])).toStrictEqual(new Set())
  })

  it("marks nothing on a chart of a single series", () => {
    expect(findStaleSeries([nightly(90 * DAY_MS)])).toStrictEqual(new Set())
  })

  it("keeps a weekend-sized gap in a per-commit series unmarked", () => {
    expect(findStaleSeries([perCommit(), perCommit(2 * DAY_MS)])).toStrictEqual(new Set())
  })

  it("marks a per-commit series that has been quiet for a working week", () => {
    expect(findStaleSeries([perCommit(), perCommit(5 * DAY_MS)])).toStrictEqual(new Set([1]))
  })

  it("marks a nightly series that missed a few runs in a row", () => {
    expect(findStaleSeries([nightly(), nightly(4 * DAY_MS)])).toStrictEqual(new Set([1]))
  })

  it("leaves a nightly series that is a run behind unmarked", () => {
    expect(findStaleSeries([nightly(), nightly(2 * DAY_MS)])).toStrictEqual(new Set())
  })

  it("leaves a weekly series unmarked while the gap is within its own cadence", () => {
    expect(findStaleSeries([nightly(), weekly(10 * DAY_MS)])).toStrictEqual(new Set())
  })

  it("marks a weekly series once the gap outgrows its cadence", () => {
    expect(findStaleSeries([nightly(), weekly(30 * DAY_MS)])).toStrictEqual(new Set([1]))
  })

  it("reads the cadence of a weekly series past its repeated timestamps", () => {
    expect(findStaleSeries([nightly(), withRepeatedMeasurements(weekly(4 * DAY_MS))])).toStrictEqual(new Set())
  })

  it("still marks a weekly series with repeated timestamps once it really stops", () => {
    expect(findStaleSeries([nightly(), withRepeatedMeasurements(weekly(30 * DAY_MS))])).toStrictEqual(new Set([1]))
  })

  it("falls back to the floor for a series whose timestamps are all the same", () => {
    expect(findStaleSeries([nightly(), [START, START, START]])).toStrictEqual(new Set([1]))
  })

  it("marks a single-point series that sits far behind the chart", () => {
    expect(findStaleSeries([nightly(), [START]])).toStrictEqual(new Set([1]))
  })

  it("ignores empty series", () => {
    expect(findStaleSeries([[], nightly(), nightly(10 * DAY_MS)])).toStrictEqual(new Set([2]))
  })

  it("marks nothing on an empty chart", () => {
    expect(findStaleSeries([])).toStrictEqual(new Set())
  })

  // The configurator judges staleness before removing outliers - on the filtered timestamps the series below
  // looks like it stopped 11 days before the chart did, though those 11 days are its own last measurement.
  it("stays unmarked when outlier removal drops its latest measurement", () => {
    const reported = nightly()
    const timestamps = [...reported, START + 111 * DAY_MS]
    const values = [...reported.map(() => 100), 1000]
    const remaining = removeOutliers([timestamps, values])[0] as number[]
    expect(remaining).not.toContain(timestamps.at(-1))

    const chartLastTimestamp = getChartLastTimestamp([timestamps])
    expect(getStaleSeriesMarkLine(timestamps, chartLastTimestamp)).toBeUndefined()
    expect(getStaleSeriesMarkLine(remaining, chartLastTimestamp)).toBeDefined()
  })

  it("puts the mark line at the last point of the stale series", () => {
    const stopped = nightly(10 * DAY_MS)
    const markLine = getStaleSeriesMarkLine(stopped, getChartLastTimestamp([nightly(), stopped]))
    expect(markLine?.data).toStrictEqual([{ xAxis: stopped.at(-1) }])
  })
})
