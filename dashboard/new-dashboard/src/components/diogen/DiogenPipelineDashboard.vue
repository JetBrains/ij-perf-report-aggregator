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
        <span class="font-semibold text-gray-700 dark:text-gray-300">User-group pair</span> — a unique combination of (user, group). Each pair counts once regardless of how many times
        the user reported. This avoids noise from crash loops and reflects breadth of impact.
      </p>
      <p>
        <span class="font-semibold text-gray-700 dark:text-gray-300">Distinct group</span> — a group created by a human-written mapping, as opposed to auto-generated groups.
        Auto-generated groups are created automatically when a report doesn't match any existing mapping — they are unstable and disappear when reports are re-processed. Distinct groups
        persist across re-processing and represent issues that someone has explicitly identified and classified.
      </p>
    </div>

    <p class="text-sm text-gray-500 dark:text-gray-400 max-w-2xl mb-4 mt-6">
      <span class="font-semibold text-gray-700 dark:text-gray-300">Incoming</span> — How much pain are our users experiencing, and how diverse are the problems? Chart going up means
      more users or more problems; chart going down means fewer.
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
      <span class="font-semibold text-gray-700 dark:text-gray-300">Sorting</span> — Do we know about the pain? Percent of user-group pairs in groups that have a YT issue attached.
      Chart going up — we know about more of the pain. Chart going down — more pain is going unnoticed.
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
      <span class="font-semibold text-gray-700 dark:text-gray-300">Fixing</span> — Did we resolve the pain? Percent of user-group pairs and groups with a fixed YT issue. Compare
      the two charts: high fixed user-group pairs + low fixed groups = good prioritization (few fixes, many users helped). Low fixed user-group pairs + high fixed groups = poor
      prioritization (many fixes, low impact).
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
