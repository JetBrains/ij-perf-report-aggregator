/**
 * Charts of the JVM cost of one measured pass: the `<measure>#jvm.*` sub-metrics that the GoLand inspection
 * performance tests publish for that pass alone (GO-21588), unlike the whole-run `JVM.*` metrics.
 *
 * @param key prefix of the chart keys; must differ between the calls of one group, which renders the charts in one list
 * @param label prefix of the chart titles, the title of the duration chart of the pass
 * @param measure duration metric of the pass, for example `globalInspections`; the charts plot its `#jvm.*` sub-metrics
 */
export function jvmCostCharts(key: string, label: string, measure: string): { key: string; label: string; measure: string }[] {
  return [
    { key: `${key}Cpu`, label: `${label}: CPU time`, measure: `${measure}#jvm.cpu.time.ms` },
    { key: `${key}Allocated`, label: `${label}: Allocated`, measure: `${measure}#jvm.alloc.mb` },
    { key: `${key}GcTime`, label: `${label}: GC time`, measure: `${measure}#jvm.gc.time.ms` },
  ]
}
