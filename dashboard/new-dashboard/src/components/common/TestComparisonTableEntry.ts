/**
 * Sometimes test results only exist for the baseline or current test, which leads to invalid table entries. A valid entry has existing
 * baseline and current results (if both exist, the difference will also exist). An invalid entry will still be shown in the table for
 * transparency, but its missing values will be replaced with "N/A".
 */
export interface TestComparisonTableEntry {
  test: string
  baselineValue: number | null
  currentValue: number | null
  difference: number | null
}

export function isValidTestComparisonTableEntry(entry: TestComparisonTableEntry) {
  return entry.baselineValue !== null && entry.currentValue !== null && entry.difference !== null
}
