<template>
  <StartupDashboardPage
    db-name="perfintDev"
    :table="table"
    :persistent-id="persistentId"
    initial-machine="Linux Munich i7-13700, 64 Gb"
    :default-project="defaultProject"
  >
    <template #default="{ projectConfigurator }">
      <Divider label="Main Metrics" />
      <section>
        <GroupProjectsWithClientChart
          v-for="chart in getMainCharts(projectConfigurator)"
          :key="generateChartKey(chart, projectConfigurator)"
          :label="chart.definition.label"
          :measure="chart.definition.measure"
          :projects="chart.projects"
        />

        <GroupProjectsWithClientChart
          v-for="chart in getCustomCharts(projectConfigurator)"
          :key="generateChartKey(chart, projectConfigurator)"
          :label="chart.definition.label"
          :measure="chart.definition.measure"
          :projects="chart.projects"
        />
      </section>
      <Accordion :lazy="true">
        <AccordionPanel value="0">
          <AccordionHeader>Additional metrics</AccordionHeader>
          <AccordionContent>
            <Divider label="Low level startup metrics" />
            <section>
              <GroupProjectsWithClientChart
                v-for="chart in getLowLevelCharts(projectConfigurator)"
                :key="generateChartKey(chart, projectConfigurator)"
                :label="chart.definition.label"
                :measure="chart.definition.measure"
                :projects="chart.projects"
              />
            </section>

            <Divider label="Highlighting" />
            <section>
              <GroupProjectsWithClientChart
                v-for="chart in getHighlightingCharts(projectConfigurator)"
                :key="generateChartKey(chart, projectConfigurator)"
                :label="chart.definition.label"
                :measure="chart.definition.measure"
                :projects="chart.projects"
              />
            </section>

            <Divider label="Notifications" />
            <section>
              <GroupProjectsWithClientChart
                v-for="chart in getNotificationsCharts(projectConfigurator)"
                :key="generateChartKey(chart, projectConfigurator)"
                :label="chart.definition.label"
                :measure="chart.definition.measure"
                :projects="chart.projects"
              />
            </section>

            <Divider label="GC metrics" />
            <section>
              <GroupProjectsWithClientChart
                v-for="chart in getGCCharts(projectConfigurator)"
                :key="generateChartKey(chart, projectConfigurator)"
                :label="chart.definition.label"
                :measure="chart.definition.measure"
                :projects="chart.projects"
              />
            </section>

            <Divider label="Exit metrics" />
            <section>
              <GroupProjectsWithClientChart
                v-for="chart in getExitCharts(projectConfigurator)"
                :key="generateChartKey(chart, projectConfigurator)"
                :label="chart.definition.label"
                :measure="chart.definition.measure"
                :projects="chart.projects"
              />
            </section>
          </AccordionContent>
        </AccordionPanel>
      </Accordion>
    </template>
  </StartupDashboardPage>
</template>

<script setup lang="ts">
import { Chart, ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import Divider from "./Divider.vue"
import StartupDashboardPage from "./StartupDashboardPage.vue"
import GroupProjectsWithClientChart from "../charts/GroupProjectsWithClientChart.vue"
import type { DimensionConfigurator } from "../../configurators/DimensionConfigurator"
import { computed } from "vue"

interface CustomChart {
  label: string
  measure: string | string[]
}

const { table, customCharts = [] } = defineProps<{
  table: string
  defaultProject: string
  customCharts?: CustomChart[]
}>()

const persistentId = computed(() => `${table}-startup-dashboard`)

const generateChartKey = (chart: Chart, configurator: DimensionConfigurator) => {
  const projectId = configurator.selected.value ?? "default"
  const projectKey = Array.isArray(projectId) ? projectId.join(",") : projectId
  return `${chart.definition.label}-${projectKey}`
}

const getProject = (configurator: DimensionConfigurator): string[] => {
  const selected = configurator.selected.value
  if (!selected) {
    return []
  }
  if (Array.isArray(selected)) {
    return selected
  }
  return [selected]
}

const getChartsForCategory = (charts: ChartDefinition[], configurator: DimensionConfigurator) => {
  const projects = getProject(configurator)
  if (projects.length === 0) {
    return []
  }
  const transformed = charts.map((chart) => ({
    ...chart,
    projects: chart.projects.flatMap((projectValue) => projects.map((project) => `${project}/${projectValue}`)),
  }))
  return combineCharts(transformed)
}

const getMainCharts = (configurator: DimensionConfigurator) => getChartsForCategory(mainCharts, configurator)
const getLowLevelCharts = (configurator: DimensionConfigurator) => getChartsForCategory(lowLevelCharts, configurator)
const getHighlightingCharts = (configurator: DimensionConfigurator) => getChartsForCategory(highlightingCharts, configurator)
const getNotificationsCharts = (configurator: DimensionConfigurator) => getChartsForCategory(notificationsCharts, configurator)
const getGCCharts = (configurator: DimensionConfigurator) => getChartsForCategory(gcCharts, configurator)
const getExitCharts = (configurator: DimensionConfigurator) => getChartsForCategory(exitCharts, configurator)

const getCustomCharts = (configurator: DimensionConfigurator) => {
  const customChartsAsChartDefinition: ChartDefinition[] = customCharts.map((chart) => ({
    labels: [chart.label],
    measures: [chart.measure],
    projects: ["measureStartup"],
  }))
  return getChartsForCategory(customChartsAsChartDefinition, configurator)
}

const mainCharts: ChartDefinition[] = [
  {
    labels: ["Reopen Project FUS code visible in editor duration"],
    measures: [["reopenProjectPerformance/fusCodeVisibleInEditorDurationMs", "fus_reopen_startup_code_loaded_and_visible_in_editor"]],
    projects: ["measureStartup"],
  },
  {
    labels: ["Reopen Project FUS first UI shown"],
    measures: [["reopenProjectPerformance/fusFirstUIShowsMs", "fus_reopen_startup_first_ui_shown"]],
    projects: ["measureStartup"],
  },
  {
    labels: ["Reopen Project FUS frame became interactive"],
    measures: [["reopenProjectPerformance/fusFrameBecameInteractiveMs", "fus_reopen_startup_frame_became_interactive"]],
    projects: ["measureStartup"],
  },
  {
    labels: ["Reopen Project FUS frame became visible"],
    measures: [["reopenProjectPerformance/fusFrameBecameVisibleMs", "fus_reopen_startup_frame_became_visible"]],
    projects: ["measureStartup"],
  },
  {
    labels: ["Startup FUS total duration"],
    measures: [["startup/fusTotalDuration", "fus_startup_totalDuration"]],
    projects: ["measureStartup"],
  },
  {
    labels: ["Total Opening Time"],
    measures: ["totalOpeningTime/timeFromAppStartTillAnalysisFinished"],
    projects: ["measureStartup"],
  },
  {
    labels: ["Code Analysis Daemon FUS execution time"],
    measures: [["codeAnalysisDaemon/fusExecutionTime", "fus_daemon_finished_full_duration_since_started_ms"]],
    projects: ["measureStartup"],
  },
  {
    labels: ["Open File"],
    measures: ["openFile"],
    projects: ["warmup"],
  },
  {
    labels: ["Reopen File After IDE Restart"],
    measures: ["reopenFileAfterIdeRestart"],
    projects: ["measureStartup"],
  },
]

const lowLevelCharts: ChartDefinition[] = [
  {
    labels: ["App Initialization"],
    measures: ["app initialization"],
    projects: ["measureStartup"],
  },
  {
    labels: ["Bootstrap"],
    measures: ["bootstrap"],
    projects: ["measureStartup"],
  },
  {
    labels: ["RunManager Initialization"],
    measures: ["RunManager initialization"],
    projects: ["measureStartup"],
  },
  {
    labels: ["Connect FSRecords"],
    measures: ["connect FSRecords"],
    projects: ["measureStartup"],
  },
  {
    labels: ["Editor Restoring"],
    measures: [["editor restoring", '"editor restoring till paint"']],
    projects: ["measureStartup"],
  },
  {
    labels: ["File Opening in EDT"],
    measures: ["file opening in EDT"],
    projects: ["measureStartup"],
  },
  {
    labels: ["Plugin Descriptor Loading"],
    measures: ["plugin descriptor loading"],
    projects: ["measureStartup"],
  },
  {
    labels: ["Class Loading Time"],
    measures: [["classLoadingTime", "classLoadingSearchTime", "classLoadingDefineTime"]],
    projects: ["measureStartup"],
  },
  {
    labels: ["Class Loading - Loaded Count"],
    measures: ["classLoadingLoadedCount"],
    projects: ["measureStartup"],
  },
  {
    labels: ["Class Loading"],
    measures: [["classLoadingMetrics/companionCount", "classLoadingMetrics/inlineCount", "classLoadingMetrics/lambdaCount", "classLoadingMetrics/methodHandleCount"]],
    projects: ["measureStartup"],
  },
  {
    labels: ["Class Loading - Prepared Count"],
    measures: ["classLoadingPreparedCount"],
    projects: ["measureStartup"],
  },
]

const highlightingCharts: ChartDefinition[] = [
  {
    labels: ["Highlighing Passes"],
    measures: [["CodeVisionPass", "GeneralHighlightingPass", "InlayHintsPass", "LocalInspectionsPass", "ShowIntentionsPass", "StickyLinesPass"]],
    projects: ["measureStartup"],
  },
]

const notificationsCharts: ChartDefinition[] = [
  {
    labels: ["Notifications Number"],
    measures: ["notificationsNumber"],
    projects: ["warmup"],
  },
]

const gcCharts: ChartDefinition[] = [
  {
    labels: ["GC - Freed Memory"],
    measures: ["gc/freedMemory"],
    projects: ["measureStartup"],
  },
  {
    labels: ["GC - Freed Memory by GC"],
    measures: ["gc/freedMemoryByGC"],
    projects: ["measureStartup"],
  },
  {
    labels: ["GC - Full GC Pause"],
    measures: ["gc/fullGCPause"],
    projects: ["measureStartup"],
  },
  {
    labels: ["GC - G1GC Concurrent Mark Cycles"],
    measures: ["gc/g1gcConcurrentMarkCycles"],
    projects: ["measureStartup"],
  },
  {
    labels: ["GC - G1GC Concurrent Mark Time (ms)"],
    measures: ["gc/g1gcConcurrentMarkTimeMs"],
    projects: ["measureStartup"],
  },
  {
    labels: ["GC - G1GC Heap Shrinkage Count"],
    measures: ["gc/g1gcHeapShrinkageCount"],
    projects: ["measureStartup"],
  },
  {
    labels: ["GC - G1GC Heap Shrinkage (MB)"],
    measures: ["gc/g1gcHeapShrinkageMegabytes"],
    projects: ["measureStartup"],
  },
  {
    labels: ["GC - GC Pause"],
    measures: ["gc/gcPause"],
    projects: ["measureStartup"],
  },
  {
    labels: ["GC - GC Pause Count"],
    measures: ["gc/gcPauseCount"],
    projects: ["measureStartup"],
  },
  {
    labels: ["GC - Total Heap Used Max"],
    measures: ["gc/totalHeapUsedMax"],
    projects: ["measureStartup"],
  },
]

const exitCharts: ChartDefinition[] = [
  {
    labels: ["Application Exit"],
    measures: ["application.exit"],
    projects: ["measureStartup"],
  },
  {
    labels: ["Save Settings on Exit"],
    measures: ["saveSettingsOnExit"],
    projects: ["measureStartup"],
  },
  {
    labels: ["Dispose Projects"],
    measures: ["disposeProjects"],
    projects: ["measureStartup"],
  },
]
</script>
