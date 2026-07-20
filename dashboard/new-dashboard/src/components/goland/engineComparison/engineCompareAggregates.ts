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
export interface GeomeanGroup {
  key: string
  geomean: number
  count: number
}

// A single worst/best cell, surfaced in the verdict.
export interface EngineExtreme {
  base: string
  title: string
  metricType: string
  diffPercent: number
}

export interface EngineAggregates {
  // Cells with a usable ratio (finite, positive) — the denominator of grandGeomean.
  count: number
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

export function computeEngineAggregates(rows: readonly EngineComparePoint[]): EngineAggregates {
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
