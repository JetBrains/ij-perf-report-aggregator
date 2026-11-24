<template>
  <DashboardPage
    db-name="jbr"
    table="report"
    persistent-id="jbr_team_dashboard"
    branch="dev_main"
    :with-installer="false"
    :is-build-number-exists="true"
  >
    <template #configurator>
      <JbrConfigurationSelect
        v-model="selectedConfigurations"
        :benchmark="benchmark"
      />
      <SelectButton
        v-if="modes.length > 0"
        v-model="selectedMode"
        :options="modes"
        option-label="label"
        option-value="id"
        :allow-empty="false"
        aria-label="The mode that the charts show"
      />
      <SelectButton
        v-if="statistics.length > 1"
        v-model="selectedStatistics"
        :options="statistics"
        option-label="label"
        option-value="id"
        multiple
        :allow-empty="false"
        aria-label="The statistics that the charts show"
      />
    </template>
    <template
      v-for="group in groups"
      :key="group.title ?? ''"
    >
      <Divider
        v-if="group.title != null"
        :title="group.title"
      />
      <section :class="['grid grid-cols-1 gap-6', group.charts.length > 1 ? 'xl:grid-cols-2' : '']">
        <GroupProjectsChart
          v-for="chart in group.charts"
          :key="chart.scenario"
          :label="chart.scenario"
          :measure="chart.measures"
          :projects="shownConfigurations"
          :aliases="shownConfigurations.map(configurationOf)"
          :description="description"
          :legend-formatter="chart.legendFormatter"
          :value-unit="valueUnit"
          :better-direction="betterDirection"
          :change-criteria="JBR_CHANGE_CRITERIA"
        />
      </section>
    </template>
  </DashboardPage>
</template>

<script setup lang="ts">
import { useStorage } from "@vueuse/core"
import { computed } from "vue"
import { BetterDirection } from "../../shared/changeDetector/algorithm"
import { ValueUnit } from "../common/chart"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import { SERIES_NAME_SEPARATOR } from "../common/DataQueryExecutor"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
import JbrConfigurationSelect from "./JbrConfigurationSelect.vue"
import { compareConfigurations, configurationOf } from "./platforms"
import { JBR_CHANGE_CRITERIA } from "./summary"

export interface ScenarioGroup {
  /** The heading above the group's charts; null for a benchmark whose scenarios form one group. */
  title: string | null
  scenarios: string[]
}

export interface Choice {
  id: string
  label: string
}

/** A statistic of a scenario, reported as the measure "<scenario>[_<mode>]<suffix>"; the empty suffix is the median. */
export interface Statistic extends Choice {
  suffix: string
}

/**
 * The dashboard of one JBR benchmark: one chart per scenario, in which the selected configurations share the time and
 * value axes, so platforms and pipelines compare by eye. Configurations, mode and statistics are chosen once for the
 * whole page in the toolbar and stay for the next visit of this browser.
 */
const {
  benchmark,
  groups: scenarioGroups,
  description,
  defaultConfigurations,
  valueUnit = "ms",
  betterDirection = "lower",
  modes = [],
  statistics = [{ id: "median", label: "Median", suffix: "" }],
} = defineProps<{
  /** The project name prefix: projects are named "<benchmark>_<configuration>". */
  benchmark: string
  groups: ScenarioGroup[]
  description: string
  /** The unit of the values; milliseconds unless the benchmark reports a rate. */
  valueUnit?: ValueUnit
  betterDirection?: BetterDirection
  /** One configuration per platform, as far as the benchmark runs on it. */
  defaultConfigurations: string[]
  /** The modes the scenarios run in, the suffix "_<mode>" of their measure names, the first the default. */
  modes?: Choice[]
  /** The statistics the benchmark reports for each scenario, the first the default. */
  statistics?: Statistic[]
}>()

const storagePrefix = "jbr_" + benchmark.toLowerCase()
const selectedConfigurations = useStorage<string[]>(
  storagePrefix + "_configurations",
  defaultConfigurations.map((configuration) => benchmark + "_" + configuration)
)
const selectedMode = useStorage<string>(storagePrefix + "_mode", modes.length > 0 ? modes[0].id : "")
const selectedStatistics = useStorage<string[]>(storagePrefix + "_statistics", [statistics[0].id])

// In platform order, so that a configuration keeps its place (and its color) whatever else is selected before it.
const shownConfigurations = computed(() => selectedConfigurations.value.toSorted(compareConfigurations))

/** Names each line by its configuration and, where a chart shows several statistics, by the statistic. */
function statisticLegend(baseMeasure: string): (name: string) => string {
  return (name: string) => {
    const parts = name.split(SERIES_NAME_SEPARATOR)
    const measureIndex = parts.findIndex((part) => part === baseMeasure || part.startsWith(baseMeasure + "#"))
    if (measureIndex === -1) return name
    const statistic = statistics.find((it) => baseMeasure + it.suffix === parts[measureIndex])
    parts[measureIndex] = statistic?.label ?? parts[measureIndex]
    return parts.join(SERIES_NAME_SEPARATOR)
  }
}

const groups = computed(() => {
  const shownStatistics = statistics.filter((statistic) => selectedStatistics.value.includes(statistic.id))
  const shown = shownStatistics.length > 0 ? shownStatistics : [statistics[0]]
  const modeSuffix = modes.length > 0 ? "_" + selectedMode.value : ""
  return scenarioGroups.map((group) => ({
    title: group.title,
    charts: group.scenarios.map((scenario) => {
      const baseMeasure = scenario + modeSuffix
      return { scenario, measures: shown.map((statistic) => baseMeasure + statistic.suffix), legendFormatter: statisticLegend(baseMeasure) }
    }),
  }))
})
</script>
