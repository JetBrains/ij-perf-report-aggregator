/**
 * Sometimes test results only exist for the baseline or current test, which leads to invalid table entries. A valid entry has existing
 * baseline and current results (if both exist, the difference will also exist). An invalid entry will still be shown in the table for
 * transparency, but its missing values will be replaced with "N/A".
 */
export interface TestComparisonTableEntry {
  test: string
  baselineValue: number | undefined
  currentValue: number | undefined
  difference: number | undefined
  branch: string | undefined
  machineName: string | undefined
  measureName: string | undefined
}

export function isValidTestComparisonTableEntry(entry: TestComparisonTableEntry) {
  return entry.baselineValue !== undefined && entry.currentValue !== undefined && entry.difference !== undefined
}
