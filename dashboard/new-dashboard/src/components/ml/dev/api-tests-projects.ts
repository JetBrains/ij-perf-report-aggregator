import { ChartDefinition } from "../../charts/DashboardCharts"
import { computed, ComputedRef } from "vue"
const machine = "Linux EC2 c5.xlarge (4 vCPU, 8 GB)"
const tests = ["inEditorCodeGeneration_function", "inEditorCodeGeneration_all", "getAvailableProfiles", "documentation_generation_calculator_sum"]
const defaultMetric = "success_rate"
const timingMetric = "response_time#mean_value"
const nameSuggestTests = [
  "java_name_suggest_PsiAnnotation",
  "java_name_suggest_PsiClass",
  "java_name_suggest_PsiVariable",
  "java_name_suggest_PsiField",
  "java_name_suggest_PsiMethod",
]
const testGenerationTests = [
  "test_generation_add_tests: openai-gpt-4",
  "test_generation_new_tests: openai-gpt-4",
  "test_generation: default",
  "test_generation_add_tests: gpt-4o",
  "test_generation_new_tests: gpt-4o",
]
export function toSeparateCharts(tests: string[]): ComputedRef<ChartDefinition[]> {
  return computed(() => tests.map((value) => toChartDefinition(value)))
}

export function toOneCharts(label: string, tests: string[], metric: string): ComputedRef<ChartDefinition[]> {
  return computed(() => [toCombineChartDefinition(label, tests, metric)])
}

export const testProjects = toSeparateCharts(tests)
export const nameSuggestProjects = toOneCharts("Name suggest", nameSuggestTests, defaultMetric)
export const timingProjects = toOneCharts("Responses", [...nameSuggestTests, ...tests], timingMetric)
export const testGenerationPassedProjects = toOneCharts("Test generation", testGenerationTests, "Passed Rate")
export const testGenerationCompilationProjects = toOneCharts("Test generation", testGenerationTests, "Compilation Rate")
export const testGenerationPsiValidProjects = toOneCharts("Test generation", testGenerationTests, "Psi Valid Rate")
function toChartDefinition(test: string): ChartDefinition {
  return {
    labels: [test + " (" + defaultMetric + ")"],
    machines: [machine],
    measures: [defaultMetric],
    projects: [test],
  }
}

function toCombineChartDefinition(label: string, tests: string[], metric: string): ChartDefinition {
  return {
    labels: [label + " (" + metric + ")"],
    machines: [machine],
    measures: [metric],
    projects: tests,
  }
}
