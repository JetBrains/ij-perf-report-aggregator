import { BetterDirection, ChangeCriteria, getChangePointIndexes, relativeChangeOfLogs } from "../../shared/changeDetector/algorithm"
import { median } from "../../shared/changeDetector/statistic"
import { MeasureUnit } from "../common/formatter"
import { compareConfigurations } from "./platforms"

/** One row of the summary: a benchmark, or one mode of it. */
export interface SummaryRow {
  label: string
  /** The project name prefix: projects are named "<benchmark>_<configuration>". */
  benchmark: string
  /** The benchmark's dashboard tab. */
  route: string
  unit: MeasureUnit
  betterDirection: Extract<BetterDirection, "lower" | "higher">
  /** The suffix of the row's measure names (Render's "_<mode>"); without it, every measure of the benchmark. */
  measureSuffix?: string
}

/** The value one CI run reported for one scenario on one configuration. */
export interface RunValue {
  project: string
  /** When the run reported, in seconds since the epoch; all values of one run share it. */
  time: number
  measure: string
  value: number
}

/**
 * How the runs after a change point compare to the runs before it: the Hodges–Lehmann estimate of the change with its
 * bounds at the misrate. Changes are performance changes, oriented by the benchmark's direction: 0.05 is 5 % better,
 * -0.05 is 5 % worse.
 */
export interface Comparison {
  change: number
  /** The bounds of the change; null with too few runs for the misrate. */
  low: number | null
  high: number | null
  runsBefore: number
  runsAfter: number
  /** The bounds lie wholly beyond a change of MIN_EFFECT: the change is confidently at least that large. */
  significant: boolean
}

interface ChangePoint extends Comparison {
  /** When the first run after the change reported, in seconds since the epoch. */
  time: number
}

interface ScenarioComparison extends Comparison {
  scenario: string
  medianBefore: number
  medianAfter: number
}

export interface SummaryCell {
  /** Every change point the detector found, oldest first. */
  changePoints: ChangePoint[]
  /** The latest significant change point with enough runs after it; null if there is none. */
  latest: ChangePoint | null
  /**
   * The last change point, if fewer than MIN_RUNS_AFTER runs followed it: too few to judge, so it is never `latest`.
   * Its comparison is provisional.
   */
  pending: ChangePoint | null
  /** The scenarios before and after the latest significant change point, worst first. */
  scenarios: ScenarioComparison[]
  /** The misrate of the scenarios' bounds: ALPHA, Bonferroni-corrected over the scenarios. */
  scenarioMisrate: number
  /** The worst scenario, if its regression at that point stays significant after the Bonferroni correction. */
  worst: ScenarioComparison | null
  runs: number
  firstRun: number
}

/** The misrate of each comparison's bounds, the significance level of its test. */
export const ALPHA = 0.05
/** A change must confidently be at least this large to be flagged: CI hardware does not keep smaller ones. */
export const MIN_EFFECT = 0.01

/** The rule by which the JBR charts mark change points, the same one that the summary applies. */
export const JBR_CHANGE_CRITERIA: ChangeCriteria = { kind: "bounds", misrate: ALPHA, minRelativeChange: MIN_EFFECT }

/** The fewest runs between two change points, as the charts' change detection uses. */
const MIN_SEGMENT = 5

/**
 * The fewest runs after the last change point to judge it. With fewer, the Bonferroni-corrected bounds of the
 * scenarios span nearly every pairwise difference, so that a single run back at the old level hides any change.
 */
export const MIN_RUNS_AFTER = 8

/**
 * The cells of one row, keyed by configuration.
 *
 * Every value becomes a log performance, negated where lower is better, so that differences are log ratios with
 * "higher is better" everywhere. Each run gets an index, the mean over its scenarios of their log ratio to the
 * scenario's median over all runs (the geometric mean of the ratios). The change point detection runs on that index,
 * one series per cell, and each change point compares the runs after it with those since the change point before.
 * The last change point is judged only once MIN_RUNS_AFTER runs followed it; until then it is pending.
 * The scenarios are compared at the latest significant change point, their bounds Bonferroni-corrected over the
 * scenarios; the worst one is flagged if it regressed significantly by them.
 */
export function summarizeRow(row: SummaryRow, values: readonly RunValue[]): Map<string, SummaryCell> {
  const performance = (value: number) => (row.betterDirection === "lower" ? -Math.log(value) : Math.log(value))
  const cells = new Map<string, SummaryCell>()
  for (const [configuration, scenarioValues] of groupByConfiguration(row, values)) {
    const baselines = new Map([...scenarioValues].map(([scenario, runs]) => [scenario, median(runs.map((it) => performance(it.value)))]))

    const ratiosByRun = new Map<number, number[]>()
    for (const [scenario, runs] of scenarioValues) {
      const baseline = baselines.get(scenario) ?? 0
      for (const run of runs) {
        let ratios = ratiosByRun.get(run.time)
        if (ratios === undefined) ratiosByRun.set(run.time, (ratios = []))
        ratios.push(performance(run.value) - baseline)
      }
    }
    // Runs that miss most scenarios (a failed or partial run) would skew the index.
    const indexed = [...ratiosByRun]
      .filter(([, ratios]) => ratios.length * 2 >= scenarioValues.size)
      .map(([time, ratios]) => ({ time, index: ratios.reduce((sum, it) => sum + it, 0) / ratios.length }))
      .toSorted((a, b) => a.time - b.time)
    if (indexed.length === 0) continue

    const indexes = indexed.map((it) => it.index)
    const boundaries = getChangePointIndexes(indexes, MIN_SEGMENT)
    const changePoints = boundaries.map((boundary, i) => {
      const start = i === 0 ? 0 : boundaries[i - 1]
      const end = i === boundaries.length - 1 ? indexes.length : boundaries[i + 1]
      return { time: indexed[boundary].time, ...compare(indexes.slice(boundary, end), indexes.slice(start, boundary), ALPHA) }
    })

    const last = changePoints.at(-1)
    const pending = last != null && last.runsAfter < MIN_RUNS_AFTER ? last : null
    const judged = pending == null ? changePoints.length : changePoints.length - 1
    const latestIndex = changePoints.slice(0, judged).findLastIndex((it) => it.significant)
    const latest = latestIndex === -1 ? null : changePoints[latestIndex]
    let scenarios: ScenarioComparison[] = []
    let worst: ScenarioComparison | null = null
    const scenarioMisrate = ALPHA / scenarioValues.size
    if (latest != null) {
      const boundary = boundaries[latestIndex]
      const from = latestIndex === 0 ? -Infinity : indexed[boundaries[latestIndex - 1]].time
      const to = latestIndex === boundaries.length - 1 ? Infinity : indexed[boundaries[latestIndex + 1]].time
      const split = indexed[boundary].time
      for (const [scenario, runs] of scenarioValues) {
        const before = runs.filter((it) => it.time >= from && it.time < split).map((it) => it.value)
        const after = runs.filter((it) => it.time >= split && it.time < to).map((it) => it.value)
        if (before.length === 0 || after.length === 0) continue
        const comparison = compare(
          after.map((value) => performance(value)),
          before.map((value) => performance(value)),
          scenarioMisrate
        )
        scenarios.push({ scenario, medianBefore: median(before), medianAfter: median(after), ...comparison })
      }
      scenarios = scenarios.toSorted((a, b) => a.change - b.change)
      const candidate = scenarios.at(0)
      if (candidate?.significant === true && candidate.change < 0) worst = candidate
    }
    cells.set(configuration, { changePoints, latest, pending, scenarios, scenarioMisrate, worst, runs: indexed.length, firstRun: indexed[0].time })
  }
  return cells
}

/** The row's values by configuration and scenario, without other benchmarks, other modes, statistics and non-values. */
function groupByConfiguration(row: SummaryRow, values: readonly RunValue[]): Map<string, Map<string, RunValue[]>> {
  const byConfiguration = new Map<string, Map<string, RunValue[]>>()
  for (const value of values) {
    const configuration = configurationOfProject(value.project, row.benchmark)
    if (configuration == null || value.measure.includes("#")) continue
    if (row.measureSuffix != null && !value.measure.endsWith(row.measureSuffix)) continue
    if (!(Number.isFinite(value.value) && value.value > 0)) continue
    const scenario = row.measureSuffix == null ? value.measure : value.measure.slice(0, -row.measureSuffix.length)
    let scenarios = byConfiguration.get(configuration)
    if (scenarios === undefined) byConfiguration.set(configuration, (scenarios = new Map()))
    let runs = scenarios.get(scenario)
    if (runs === undefined) scenarios.set(scenario, (runs = []))
    runs.push(value)
  }
  return byConfiguration
}

/** Compares two samples of log performances: the shift estimate and its bounds at the misrate, as changes. */
function compare(after: number[], before: number[], misrate: number): Comparison {
  const { change, lower, upper, significant } = relativeChangeOfLogs(after, before, misrate, MIN_EFFECT)
  return { change, low: lower, high: upper, runsBefore: before.length, runsAfter: after.length, significant }
}

/** The configuration part of a project name of the given benchmark. */
function configurationOfProject(project: string, benchmark: string): string | null {
  return project.startsWith(benchmark + "_") ? project.slice(benchmark.length + 1) : null
}

/** Orders configuration names by platform (macOS, Linux, Windows, then the rest), then by name. */
export function compareConfigurationNames(a: string, b: string): number {
  return compareConfigurations("_" + a, "_" + b)
}
