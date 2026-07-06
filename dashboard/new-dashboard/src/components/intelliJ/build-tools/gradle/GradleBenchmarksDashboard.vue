<template>
  <DashboardPage
    db-name="perfUnitTests"
    table="report"
    persistent-id="gradleBenchmarksDashboard"
    initial-machine="linux-blade-hetzner"
    :charts="charts"
    :with-installer="false"
  >
    <section>
      <GroupProjectsChart
        v-for="chart in charts"
        :key="chart.definition.label"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../../../common/DashboardPage.vue"
import GroupProjectsChart from "../../../charts/GroupProjectsChart.vue"
import { ChartDefinition, combineCharts } from "../../../charts/DashboardCharts"

const MEASURE = "attempt.mean.ms"

// Only the current test names (com.intellij.gradle.*.tests.benchmark.*) are charted. The tests were renamed ~2026-07 from
// org.jetbrains.plugins.gradle.* when they were extracted into a dedicated GradleBenchmarkTests build; the pre-rename baselines
// are not comparable (different classpath/plugins/host), so the retired project names are deliberately excluded.
const IMPORTING_PREFIX = "com.intellij.gradle.java.tests.benchmark.importing."
const SERVICE_PROJECT_PREFIX = "com.intellij.gradle.tests.benchmark.service.project."
const DSL_PREFIX = "com.intellij.gradle.java.tests.benchmark.dsl.groovy."
const KOTLIN_COMPLETION_PREFIX = "com.intellij.gradle.completion.kotlin.tests.benchmark."

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["Gradle sync"],
    measures: [MEASURE],
    projects: [
      IMPORTING_PREFIX + "GradleJavaSyncPerformanceTest$ReSync.test - Gradle sync (Gradle 9.6.0, SIMPLE_PROJECT)",
      IMPORTING_PREFIX + "GradleJavaSyncPerformanceTest$ReSync.test - Gradle sync (Gradle 9.6.0, SOURCE_SET_DEPENDENCY_PROJECT)",
    ],
  },
  {
    labels: ["Java project resolver"],
    measures: [MEASURE],
    projects: [
      IMPORTING_PREFIX + "GradleJavaProjectResolverPerformanceTest$SdkData.test JavaGradleProjectResolver#populateModuleExtraModels - 1000 x 100 modules",
      IMPORTING_PREFIX + "GradleJavaProjectResolverPerformanceTest$SdkData.test JavaGradleProjectResolver#populateModuleExtraModels - 10000 x 10 modules",
      IMPORTING_PREFIX + "GradleJavaProjectResolverPerformanceTest$SdkData.test JavaGradleProjectResolver#populateProjectExtraModels - 1000 x 100 modules",
      IMPORTING_PREFIX + "GradleJavaProjectResolverPerformanceTest$SdkData.test JavaGradleProjectResolver#populateProjectExtraModels - 10000 x 10 modules",
    ],
  },
  {
    labels: ["Module data index"],
    measures: [MEASURE],
    projects: [
      SERVICE_PROJECT_PREFIX + "GradleModuleDataIndexPerformanceTest.test performance of GradleModuleDataIndex#findModuleData - 1000 x 100 modules",
      SERVICE_PROJECT_PREFIX + "GradleModuleDataIndexPerformanceTest.test performance of GradleModuleDataIndex#findModuleData - 10000 x 10 modules",
    ],
  },
  {
    labels: ["Groovy DSL highlighting"],
    measures: [MEASURE],
    projects: [
      DSL_PREFIX + "GradleHighlightingPerformanceTest$testPerformance$1$1$2.testPerformance",
      DSL_PREFIX + "GradleHighlightingPerformanceTest$testCompletionPerformance$1$1$3.testCompletionPerformance",
    ],
  },
  {
    labels: ["Kotlin dependencies completion"],
    measures: [MEASURE],
    projects: [
      KOTLIN_COMPLETION_PREFIX + "KotlinGradleDependenciesCompletionPerformanceTest.test completing a scope",
      KOTLIN_COMPLETION_PREFIX + "KotlinGradleDependenciesCompletionPerformanceTest.test completing a scope argument without quotes",
      KOTLIN_COMPLETION_PREFIX + "KotlinGradleDependenciesCompletionPerformanceTest.test completing gav coordinates inside a scope argument",
      KOTLIN_COMPLETION_PREFIX + "KotlinGradleDependenciesCompletionPerformanceTest.test completing a Dependency-returning method argument without quotes",
    ],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
