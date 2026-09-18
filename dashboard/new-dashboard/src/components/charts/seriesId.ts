export function getSeriesId(measureName: string, seriesName: string): string {
  return measureName === seriesName ? seriesName : `${measureName}@${seriesName}`
}

const DERIVED_SERIES_ROLES = ["smoothed", "stdDevBandLowerEdge", "stdDevBandFill"] as const

// Derived series share their owner's hover group and are excluded from exports.
export function parseSeriesId(seriesId: string): { ownerSeriesId: string; role: (typeof DERIVED_SERIES_ROLES)[number] | null } {
  const role = DERIVED_SERIES_ROLES.find((role) => seriesId.endsWith(role)) ?? null
  return { ownerSeriesId: role == null ? seriesId : seriesId.slice(0, -role.length), role }
}
