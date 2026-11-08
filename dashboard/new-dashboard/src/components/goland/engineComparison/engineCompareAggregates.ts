import { DIFF_PERCENT_WARN } from "../../charts/compareStats"

// Pure aggregation of the LEGACY-vs-NEW engine comparison, following SPEC/Fleming-Wallace methodology:
// the summary statistic is the geometric mean of per-cell ratios (NEW/LEGACY). A single grand geomean
// hides a Simpson's-paradox split (NEW can win overall while losing every bucket), so per-bucket and
// per-phase geomeans are reported alongside it. "Faster/slower" is three-way (improved / regressed /
// neutral), gated by the same significance test the compare table uses.

// The minimal per-cell input the aggregation reads. The composable's richer row (with the raw stats)
// is assignable to this, and unit tests can build plain objects.
export interface EngineComparePoint {
  base: string
  title: string
  metricType: string
  phase: string
  bucket: string
  // NEW.center / LEGACY.center; NaN when the ratio is not computable (missing or non-positive baseline).
  ratio: number
  // (NEW.center − LEGACY.center) / LEGACY.center * 100; NaN when not computable.
  diffPercent: number
  // Both statistically and practically different from LEGACY (see isSignificantBranch).
  significant: boolean
}

// Geomean of the ratios in one group (fast/medium/slow bucket, or cold/warm/typing phase).
interface GeomeanGroup {
  key: string
  geomean: number
  count: number
}

// A single worst/best cell, surfaced in the verdict.
interface EngineExtreme {
  base: string
  title: string
  metricType: string
  diffPercent: number
}

// A cell that one engine did not publish in some builds of the compared scope (see countMissing).
export interface EngineGap {
  base: string
  title: string
  metricType: string
  engine: "LEGACY" | "NEW"
  builds: number
}

export interface EngineAggregates {
  // Cells with a usable ratio (finite, positive) — the denominator of grandGeomean.
  count: number
  // Cells the comparison covers: every project × every selected metric type, measured or not.
  cells: number
  gaps: EngineGap[]
  // Geomean of NEW/LEGACY across all usable cells; NaN when none. < 1 means NEW is faster overall.
  grandGeomean: number
  // Per-bucket geomeans, in the order buckets first appear (fast -> medium -> slow).
  perBucket: GeomeanGroup[]
  // Per-phase geomeans, in the order phases first appear (cold -> warm -> typing).
  perPhase: GeomeanGroup[]
  // Three-way counts over every cell: significant slowdowns, significant speedups, and the rest.
  regressed: number
  improved: number
  neutral: number
  // The significant slowdown with the largest Δ%, and the significant speedup with the most negative Δ%.
  worstRegression: EngineExtreme | null
  bestImprovement: EngineExtreme | null
}

// Geomean over the finite positive ratios only. Zero/negative/NaN ratios are excluded (a ratio needs a
// positive baseline to be meaningful), and an all-excluded group yields NaN with count 0.
function geomeanOf(ratios: readonly number[]): { geomean: number; count: number } {
  const valid = ratios.filter((r) => Number.isFinite(r) && r > 0)
  if (valid.length === 0) return { geomean: Number.NaN, count: 0 }
  const sumLn = valid.reduce((acc, r) => acc + Math.log(r), 0)
  return { geomean: Math.exp(sumLn / valid.length), count: valid.length }
}

// Groups rows by `keyFn` preserving first-appearance order, then computes each group's geomean. Rows
// arrive in phase × bucket order, so first-appearance already yields fast/medium/slow and cold/warm/typing.
function groupGeomean(rows: readonly EngineComparePoint[], keyFn: (row: EngineComparePoint) => string): GeomeanGroup[] {
  const order: string[] = []
  const buckets = new Map<string, number[]>()
  for (const row of rows) {
    const key = keyFn(row)
    if (key === "") continue
    let bucket = buckets.get(key)
    if (bucket == null) {
      bucket = []
      buckets.set(key, bucket)
      order.push(key)
    }
    bucket.push(row.ratio)
  }
  return order.map((key) => {
    const { geomean, count } = geomeanOf(buckets.get(key) ?? [])
    return { key, geomean, count }
  })
}

export function computeEngineAggregates(rows: readonly EngineComparePoint[], gaps: readonly EngineGap[], cells: number): EngineAggregates {
  const { geomean: grandGeomean, count } = geomeanOf(rows.map((row) => row.ratio))

  let regressed = 0
  let improved = 0
  let neutral = 0
  let worstRegression: EngineExtreme | null = null
  let bestImprovement: EngineExtreme | null = null

  for (const row of rows) {
    if (!row.significant || !Number.isFinite(row.diffPercent) || row.diffPercent === 0) {
      neutral++
      continue
    }
    const extreme: EngineExtreme = { base: row.base, title: row.title, metricType: row.metricType, diffPercent: row.diffPercent }
    if (row.diffPercent > 0) {
      regressed++
      if (worstRegression == null || row.diffPercent > worstRegression.diffPercent) worstRegression = extreme
    } else {
      improved++
      if (bestImprovement == null || row.diffPercent < bestImprovement.diffPercent) bestImprovement = extreme
    }
  }

  return {
    count,
    cells,
    gaps: [...gaps],
    grandGeomean,
    perBucket: groupGeomean(rows, (row) => row.bucket),
    perPhase: groupGeomean(rows, (row) => row.phase),
    regressed,
    improved,
    neutral,
    worstRegression,
    bestImprovement,
  }
}

// -1 when NEW is lower (better) by a practically relevant margin, 1 when higher, 0 otherwise.
export function geomeanDirection(geomean: number): -1 | 0 | 1 {
  if (!Number.isFinite(geomean) || geomean <= 0) return 0
  const percent = (geomean - 1) * 100
  if (Math.abs(percent) < DIFF_PERCENT_WARN) return 0
  return percent < 0 ? -1 : 1
}

export type EngineVerdict =
  | { kind: "noData" }
  | { kind: "incomplete"; newGaps: number; legacyGaps: number }
  | { kind: "same" }
  | { kind: "mixed" }
  | { kind: "better" }
  | { kind: "worse" }

// The headline over the aggregates. NEW is better or worse only when the grand geomean moves by a
// practically relevant margin, no bucket or phase moves the other way, and the significant cells against
// the move are fewer than those with it. In a single run, a cell without both results also blocks it;
// over a window, one failed build would block it until the window moves past, so the gaps stay a caveat.
export function resolveVerdict(aggregates: EngineAggregates, single: boolean): EngineVerdict {
  if (aggregates.count === 0) return { kind: "noData" }
  if (single && aggregates.gaps.length > 0) {
    const newGaps = aggregates.gaps.filter((gap) => gap.engine === "NEW").length
    return { kind: "incomplete", newGaps, legacyGaps: aggregates.gaps.length - newGaps }
  }
  const groupDirections = new Set<number>([...aggregates.perBucket, ...aggregates.perPhase].map((group) => geomeanDirection(group.geomean)))
  const direction = geomeanDirection(aggregates.grandGeomean)
  if (groupDirections.has(-1) && groupDirections.has(1)) return { kind: "mixed" }
  if (direction === 0) return { kind: "same" }
  if (groupDirections.has(-direction)) return { kind: "mixed" }
  const against = direction < 0 ? aggregates.regressed : aggregates.improved
  const along = direction < 0 ? aggregates.improved : aggregates.regressed
  if (against > 0 && against >= along) return { kind: "mixed" }
  return { kind: direction < 0 ? "better" : "worse" }
}
