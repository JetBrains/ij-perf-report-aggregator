<template>
  <DashboardPage
    db-name="diogen"
    table="report"
    persistent-id="diogen_pipeline_dashboard"
    :with-installer="false"
    :branch="null"
    :initial-machine="null"
  >
    <template #configurator>
      <MeasureSelect
        :configurator="projectConfigurator"
        title="Domain"
      >
        <template #icon>
          <ChartBarIcon class="w-4 h-4" />
        </template>
      </MeasureSelect>
    </template>

    <div class="text-sm text-gray-500 dark:text-gray-400 max-w-2xl mb-4 mt-2">
      <p class="mb-2">
        <span class="font-semibold text-gray-700 dark:text-gray-300">User-group pair</span> — one user affected by one issue. If the same user hits the same issue 50 times, it still
        counts as one pair. This gives a fair measure of how widespread the pain is, without noise from repeated reports.
      </p>
      <p>
        <span class="font-semibold text-gray-700 dark:text-gray-300">Distinct group</span> — a meaningful issue identified by a human. When reports come in, some get matched to known
        issues (via mappings written by engineers), while others fall into auto-generated buckets. Only the known, mapped issues count as distinct groups. More distinct groups means we're
        tracking more individual problems.
      </p>
    </div>

    <p class="text-sm text-gray-500 dark:text-gray-400 max-w-2xl mb-4 mt-6">
      <span class="font-semibold text-gray-700 dark:text-gray-300">Incoming</span> — How much pain are our users experiencing, and how diverse are the problems? Charts going up means
      more users or more problems; going down means fewer.
    </p>
    <section>
      <GroupProjectsChart
        v-for="chart in chartsIncoming"
        :key="chart.definition.label"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
        :value-unit="'counter'"
      />
    </section>

    <p class="text-sm text-gray-500 dark:text-gray-400 max-w-2xl mb-4 mt-6">
      <span class="font-semibold text-gray-700 dark:text-gray-300">Sorting</span> — Do we know about the pain? Shows what percent of reported problems has a tracked issue in YouTrack. Going
      up — we're aware of more problems. Going down — more pain is slipping through unnoticed.
    </p>
    <section>
      <GroupProjectsChart
        v-for="chart in chartsSorting"
        :key="chart.definition.label"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
        :value-unit="'counter'"
      />
    </section>

    <p class="text-sm text-gray-500 dark:text-gray-400 max-w-2xl mb-4 mt-6">
      <span class="font-semibold text-gray-700 dark:text-gray-300">Fixing</span> — Are we resolving the pain? Shows what percent of reported problems and known issues have been fixed. Compare
      the two charts: high fixed pairs + low fixed groups means good prioritization — few fixes but many users helped. The opposite means we're fixing many small issues without much
      impact.
    </p>
    <section>
      <GroupProjectsChart
        v-for="chart in chartsFixing"
        :key="chart.definition.label"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
        :value-unit="'counter'"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import { computed } from "vue"
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import DashboardPage from "../common/DashboardPage.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import { SimpleMeasureConfigurator } from "../../configurators/SimpleMeasureConfigurator"
import MeasureSelect from "../charts/MeasureSelect.vue"

const DOMAINS = ["ijAllCur", "freezeFlat", "oom"]

const projectConfigurator = new SimpleMeasureConfigurator("project", null)
projectConfigurator.initData(DOMAINS)

function filterBySelectedDomains(charts: ChartDefinition[]): ChartDefinition[] {
  const selected = projectConfigurator.selected.value ?? []
  return charts
    .map((chart) => ({
      ...chart,
      projects: chart.projects.filter((p) => selected.includes(p)),
    }))
    .filter((chart) => chart.projects.length > 0)
}

const chartsIncomingDeclaration: ChartDefinition[] = [
  {
    labels: ["Total User-Group Pairs"],
    measures: ["count-user-group-pairs"],
    projects: DOMAINS,
  },
  {
    labels: ["Number of Distinct Groups"],
    measures: ["count-distinct-groups"],
    projects: DOMAINS,
  },
]

const chartsSortingDeclaration: ChartDefinition[] = [
  {
    labels: ["Covered User-Group Pairs, %"],
    measures: ["percent-covered-user-group-pairs"],
    projects: DOMAINS,
  },
]

const chartsFixingDeclaration: ChartDefinition[] = [
  {
    labels: ["Fixed User-Group Pairs, %"],
    measures: ["percent-fixed-user-group-pairs"],
    projects: DOMAINS,
  },
  {
    labels: ["Fixed Distinct Groups, %"],
    measures: ["percent-fixed-distinct-groups"],
    projects: DOMAINS,
  },
]

const chartsIncoming = computed(() => combineCharts(filterBySelectedDomains(chartsIncomingDeclaration)))
const chartsSorting = computed(() => combineCharts(filterBySelectedDomains(chartsSortingDeclaration)))
const chartsFixing = computed(() => combineCharts(filterBySelectedDomains(chartsFixingDeclaration)))
</script>
