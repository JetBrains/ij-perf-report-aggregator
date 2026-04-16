<template>
  <DashboardPage
    db-name="perfintDev"
    table="clion"
    persistent-id="clion_lagging_latency_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :with-installer="false"
  >
    <section>
      <Divider title="Lagging during indexing" />
      <Accordion value="0">
        <AccordionPanel value="0">
          <AccordionHeader>Do not skip heavy metric collection</AccordionHeader>
          <AccordionContent>
            <GroupProjectsChart
              v-for="chart in laggingIndexingChartsCombined"
              :key="chart.definition.label"
              :label="chart.definition.label"
              :measure="chart.definition.measure"
              :projects="chart.projects"
            />
          </AccordionContent>
        </AccordionPanel>
        <AccordionPanel value="1">
          <AccordionHeader>Skip heavy metric collection</AccordionHeader>
          <AccordionContent>
            <GroupProjectsChart
              v-for="chart in pureLaggingIndexingChartsCombined"
              :key="chart.definition.label"
              :label="chart.definition.label"
              :measure="chart.definition.measure"
              :projects="chart.projects"
            />
          </AccordionContent>
        </AccordionPanel>
      </Accordion>

      <Divider title="Lagging during zephyr project indexing" />
      <Accordion value="0">
        <AccordionPanel value="0">
          <AccordionHeader>Do not skip heavy metric collection</AccordionHeader>
          <AccordionContent>
            <GroupProjectsChart
              v-for="chart in laggingZephyrIndexingChartsCombined"
              :key="chart.definition.label"
              :label="chart.definition.label"
              :measure="chart.definition.measure"
              :projects="chart.projects"
            />
          </AccordionContent>
        </AccordionPanel>
        <AccordionPanel value="1">
          <AccordionHeader>Skip heavy metric collection</AccordionHeader>
          <AccordionContent>
            <GroupProjectsChart
              v-for="chart in pureLaggingZephyrIndexingChartsCombined"
              :key="chart.definition.label"
              :label="chart.definition.label"
              :measure="chart.definition.measure"
              :projects="chart.projects"
            />
          </AccordionContent>
        </AccordionPanel>
      </Accordion>

      <Divider title="Lagging during completion" />
      <Accordion value="0">
        <AccordionPanel value="0">
          <AccordionHeader>Do not skip heavy metric collection</AccordionHeader>
          <AccordionContent>
            <GroupProjectsChart
              v-for="chart in laggingCompletionChartsCombined"
              :key="chart.definition.label"
              :label="chart.definition.label"
              :measure="chart.definition.measure"
              :projects="chart.projects"
            />
          </AccordionContent>
        </AccordionPanel>
        <AccordionPanel value="1">
          <AccordionHeader>Skip heavy metric collection</AccordionHeader>
          <AccordionContent>
            <GroupProjectsChart
              v-for="chart in pureLaggingCompletionChartsCombined"
              :key="chart.definition.label"
              :label="chart.definition.label"
              :measure="chart.definition.measure"
              :projects="chart.projects"
            />
          </AccordionContent>
        </AccordionPanel>
      </Accordion>

      <Divider title="Lagging during navigation" />
      <Accordion value="0">
        <AccordionPanel value="0">
          <AccordionHeader>Do not skip heavy metric collection</AccordionHeader>
          <AccordionContent>
            <GroupProjectsChart
              v-for="chart in laggingNavigationChartsCombined"
              :key="chart.definition.label"
              :label="chart.definition.label"
              :measure="chart.definition.measure"
              :projects="chart.projects"
            />
          </AccordionContent>
        </AccordionPanel>
        <AccordionPanel value="1">
          <AccordionHeader>Skip heavy metric collection</AccordionHeader>
          <AccordionContent>
            <GroupProjectsChart
              v-for="chart in pureLaggingNavigationChartsCombined"
              :key="chart.definition.label"
              :label="chart.definition.label"
              :measure="chart.definition.measure"
              :projects="chart.projects"
            />
          </AccordionContent>
        </AccordionPanel>
      </Accordion>

      <Divider title="Lagging during browsing files" />
      <GroupProjectsChart
        v-for="chart in laggingHighlightingChartsCombined"
        :key="chart.definition.label"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
      />
      <Divider title="Lagging during debugging" />
      <GroupProjectsChart
        v-for="chart in laggingDebuggingChartsCombined"
        :key="chart.definition.label"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"

function createLaggingCharts(label: string, projects: string[], aliases: string[]): ChartDefinition[] {
  return [
    { labels: [`${label} - average, max`], measures: [["ui.lagging#average", "ui.lagging#max"]], projects, aliases },
    { labels: [`${label} - sum`], measures: ["ui.lagging#sum"], projects, aliases },
    { labels: [`${label} - count`], measures: ["ui.lagging#count"], projects, aliases },
    { labels: [`${label} - percentage share`], measures: [["ui.lagging#percentage_share"]], projects, aliases },
  ]
}

const indexingProjects = ["radler/llvm/indexing", "radler/opencv/indexing", "radler/big_project_50k_10k_many_symbols/indexing"]
const laggingIndexingProjects = ["radler/llvm/lagging/indexing", "radler/opencv/lagging/indexing", "radler/big_project_50k_10k_many_symbols/lagging/indexing"]
const indexingAliases = ["LLVM", "OpenCV", "Big Project Many Symbols"]

const zephyrIndexingProjects = ["radler/zephyr_bap_broadcast_sink/indexing"]
const zephyrLaggingIndexingProjects = ["radler/zephyr_bap_broadcast_sink/lagging/indexing"]
const zephyrIndexingAliases = ["Zephyr Bap Broadcast Sink"]

const completionProjects = [
  "radler/fmtlib/completion/fmt.join_view (dep) (hot)",
  "radler/fmtlib/completion/std.shared_ptr (dep) (hot)",
  "radler/fmtlib/completion/std.string (hot)",
]
const completionLaggingProjects = ["radler/fmtlib/lagging/completion/fmt.join_view (dep) (hot)"]
const completionAliases = ["fmt.join_view (dep) (hot)", "std.shared_ptr (dep) (hot)", "std.string (hot)"]

const navigationProjects = ["radler/luau/findUsages/class template (DenseHashTable)", "radler/luau/gotoDeclaration/time.h", "radler/luau/gotoDeclaration/TypeChecker.getScopes"]
const navigationLaggingProjects = [
  "radler/luau/lagging/findUsages/class template (DenseHashTable)",
  "radler/luau/lagging/gotoDeclaration/time.h",
  "radler/luau/lagging/gotoDeclaration/TypeChecker.getScopes",
]
const navigationAliases = ["class template (DenseHashTable)", "time.h", "TypeChecker.getScopes"]

const syntaxHighlightingProjects = ["radler/opencv/syntaxHighlighting/opencv"]
const syntaxHighlightingAliases = ["syntaxHighlighting opencv"]

const debugProjects = ["radler/fmtlib/debug/args-test/basic"]
const debugAliases = ["fmtlib"]

const laggingHighlightingCharts: ChartDefinition[] = [
  {
    labels: ["Lagging during browsing - average, max"],
    measures: [["ui.lagging#average", "ui.lagging#max", "ui.lagging#percentage_share"]],
    projects: syntaxHighlightingProjects,
    aliases: syntaxHighlightingAliases,
  },
  {
    labels: ["Lagging during browsing - lagging percentage share"],
    measures: [["ui.lagging#percentage_share"]],
    projects: syntaxHighlightingProjects,
    aliases: syntaxHighlightingAliases,
  },
]

const laggingIndexingChartsCombined = combineCharts(createLaggingCharts("Lagging during indexing", indexingProjects, indexingAliases))
const pureLaggingIndexingChartsCombined = combineCharts(createLaggingCharts("Lagging during indexing", laggingIndexingProjects, indexingAliases))

const laggingZephyrIndexingChartsCombined = combineCharts(createLaggingCharts("Lagging during indexing", zephyrIndexingProjects, zephyrIndexingAliases))
const pureLaggingZephyrIndexingChartsCombined = combineCharts(createLaggingCharts("Lagging during indexing", zephyrLaggingIndexingProjects, zephyrIndexingAliases))

const laggingCompletionChartsCombined = combineCharts(createLaggingCharts("Lagging during completion", completionProjects, completionAliases))
const pureLaggingCompletionChartsCombined = combineCharts(createLaggingCharts("Lagging during completion", completionLaggingProjects, completionAliases))

const laggingNavigationChartsCombined = combineCharts(createLaggingCharts("Lagging during navigation", navigationProjects, navigationAliases))
const pureLaggingNavigationChartsCombined = combineCharts(createLaggingCharts("Lagging during navigation", navigationLaggingProjects, navigationAliases))

const laggingHighlightingChartsCombined = combineCharts(laggingHighlightingCharts)
const laggingDebuggingChartsCombined = combineCharts(createLaggingCharts("Lagging during debugging", debugProjects, debugAliases))
</script>
