<template>
  <Divider label="Bootstrap" />
  <section class="grid grid-cols-2 gap-x-6">
    <LineChart
      :measures="['appInit_d', 'app initialization.end']"
      title="App Initialization"
      :configurators="configurators"
      tooltip-trigger="axis"
    />
    <LineChart
      :measures="['bootstrap_d']"
      title="Bootstrap"
      :configurators="configurators"
      :with-measure-name="true"
    />
  </section>

  <section class="grid grid-cols-2 gap-x-6">
    <LineChart
      :measures="['classLoadingPreparedCount', 'classLoadingLoadedCount']"
      title="Class Loading (Count)"
      :configurators="configurators"
      tooltip-trigger="axis"
    />
    <LineChart
      :configurators="configurators"
      :measures="['editorRestoring']"
      title="Editor restoring"
      :with-measure-name="true"
    />
  </section>

  <span v-if="highlightingPasses">
    <Divider label="Highlighting Passes" />
    <span v-if="showAllPasses">
      <LineChart
        title="Highlighting Passes"
        :measures="highlightingPasses"
        :configurators="configurators"
      />
    </span>
    <LineChart
      title="Code Analysis"
      :measures="['metrics.codeAnalysisDaemon/fusExecutionTime', 'metrics.runDaemon/executionTime']"
      :configurators="configurators"
      tooltip-trigger="axis"
    />
  </span>
  <slot :configurators="configurators"></slot>
  <Divider label="Notifications" />
  <LineChart
    title="Notifications"
    :measures="['metrics.notifications/number']"
    :skip-zero-values="false"
    :configurators="configurators"
    :with-measure-name="true"
  />

  <Divider label="Exit" />
  <LineChart
    title="Exit Metrics"
    :measures="['metrics.exitMetrics/application.exit', 'metrics.exitMetrics/saveSettingsOnExit', 'metrics.exitMetrics/disposeProjects']"
    :configurators="configurators"
    tooltip-trigger="axis"
  />
</template>
<script setup lang="ts">
import Divider from "./Divider.vue"
import LineChart from "../charts/LineChart.vue"
import { fetchHighlightingPasses } from "../startup/utils"
import { computed } from "vue"
import { DataQueryConfigurator } from "./dataQuery"
import { FilterConfigurator } from "../../configurators/filter"
import { DimensionConfigurator } from "../../configurators/DimensionConfigurator"

const highlightingPasses = fetchHighlightingPasses()
interface AdditionalMetricsProps {
  configurators: (DataQueryConfigurator | FilterConfigurator)[]
  projectConfigurator: DimensionConfigurator
}
const props = defineProps<AdditionalMetricsProps>()
const showAllPasses = computed(() => {
  return props.projectConfigurator.selected.value == null || props.projectConfigurator.selected.value.length == 1 || typeof props.projectConfigurator.selected.value == "string"
})
</script>
