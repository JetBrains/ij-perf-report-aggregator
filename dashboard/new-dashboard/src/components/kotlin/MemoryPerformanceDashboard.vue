<template>
  <DashboardPage
    db-name="perfint"
    table="kotlin"
    persistent-id="kotlin_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <template #configurator>
      <MeasureSelect
        :configurator="measureConfigurator"
        title="Measure"
        :label-formatter="metricsSelectLabelFormat"
      >
        <template #icon>
          <ChartBarIcon class="w-4 h-4 text-gray-500" />
        </template>
      </MeasureSelect>
    </template>
    <Divider title="Completion" />
    <MemoryDashboardGroupCharts
      :metrics="metrics"
      :definitions="completionCharts"
    />
    <Divider title="Code analysis" />
    <MemoryDashboardGroupCharts
      :metrics="metrics"
      :definitions="codeAnalysisCharts"
    />
    <MemoryDashboardGroupCharts
      :metrics="metrics"
      :definitions="highlightingCharts"
    />
    <Divider title="Find usages" />
    <MemoryDashboardGroupCharts
      :metrics="metrics"
      :definitions="findUsagesCharts"
    />
    <Divider title="Debugger" />
    <MemoryDashboardGroupCharts
      :metrics="metrics"
      :definitions="evaluateExpressionChars"
    />
    <Divider title="Refactoring" />
    <MemoryDashboardGroupCharts
      :metrics="metrics"
      :definitions="refactoringCharts"
    />
    <Divider title="Script" />
    <MemoryDashboardGroupCharts
      :metrics="metrics"
      :definitions="scriptCompletionCharts"
    />
    <MemoryDashboardGroupCharts
      :metrics="metrics"
      :definitions="highlightingScriptCharts"
    />
    <MemoryDashboardGroupCharts
      :metrics="metrics"
      :definitions="codeAnalysisScriptCharts"
    />
  </DashboardPage>
</template>

<script setup lang="ts">
import { SimpleMeasureConfigurator } from "../../configurators/SimpleMeasureConfigurator"
import { metricsSelectLabelFormat } from "../../shared/labels"
import MeasureSelect from "../charts/MeasureSelect.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
import MemoryDashboardGroupCharts from "./MemoryDashboardGroupCharts.vue"
import {
  completionCharts,
  evaluateExpressionChars,
  findUsagesCharts,
  highlightingCharts,
  codeAnalysisCharts,
  refactoringCharts,
  scriptCompletionCharts,
  highlightingScriptCharts,
  codeAnalysisScriptCharts,
} from "./projects"
import { computed, Ref, ref } from "vue"

const measureConfigurator = new SimpleMeasureConfigurator("metrics", null)
measureConfigurator.initData(["freedMemoryByGC", "JVM.diffHeapUsageMb_afterGc"])

const metrics = computed(() => {
  const reference = measureConfigurator.selected
  if (reference.value === null) {
    return ref([])
  }
  return reference as Ref<string[]>
})
</script>
