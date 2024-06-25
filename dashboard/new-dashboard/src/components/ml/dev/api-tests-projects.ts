import { ChartDefinition } from "../../charts/DashboardCharts"
import { computed, ComputedRef } from "vue"
const machine = "Linux EC2 c5.xlarge (4 vCPU, 8 GB)"
const tests = ["inEditorCodeGeneration_function", "inEditorCodeGeneration_all", "getAvailableProfiles"]
const defaultMetrics = ["success_count", "failure_count"]
const timingMetrics = ["response_time#mean_value"]
const nameSuggestTests = [
  "java_name_suggest_PsiAnnotation",
  "java_name_suggest_PsiClass",
  "java_name_suggest_PsiVariable",
  "java_name_suggest_PsiField",
  "java_name_suggest_PsiMethod",
]

export function toSeparateCharts(tests: string[]): ComputedRef<ChartDefinition[]> {
  return computed(() => tests.map((value) => toChartDefinition(value)))
}

export function toOneCharts(label: string, tests: string[], metrics: string[]): ComputedRef<ChartDefinition[]> {
  return computed(() => [toCombineChartDefinition(label, tests, metrics)])
}

export const testProjects = toSeparateCharts(tests)
export const nameSuggestProjects = toOneCharts("Name suggest", nameSuggestTests, defaultMetrics)
export const timingProjects = toOneCharts("Mean response time", [...nameSuggestTests, ...tests], timingMetrics)
function toChartDefinition(test: string): ChartDefinition {
  return {
    labels: [test],
    machines: [machine],
    measures: defaultMetrics,
    projects: [test],
  }
}

function toCombineChartDefinition(label: string, tests: string[], metrics: string[]): ChartDefinition {
  return {
    labels: [label],
    machines: [machine],
    measures: metrics,
    projects: tests,
  }
}
