<template>
  <DashboardPage
    db-name="perfintDev"
    table="kotlin"
    persistent-id="kotlin_dashboard_dev"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  >
    <template #configurator>
      <MeasureSelect
        :configurator="measureConfigurator"
        title="Measure"
        :label-formatter="metricsSelectLabelFormat"
      >
        <template #icon>
          <ChartBarIcon class="w-4 h-4" />
        </template>
      </MeasureSelect>
    </template>
    <Divider title="Completion" />
    <MemoryK1K2DashboardGroupCharts
      :metrics="metrics"
      :definitions="completionCharts"
    />
    <Divider title="Code analysis" />
    <MemoryK1K2DashboardGroupCharts
      :metrics="metrics"
      :definitions="codeAnalysisCharts"
    />
    <MemoryK1K2DashboardGroupCharts
      :metrics="metrics"
      :definitions="highlightingCharts"
    />
    <Divider title="Find usages" />
    <MemoryK1K2DashboardGroupCharts
      :metrics="metrics"
      :definitions="findUsagesCharts"
    />
    <Divider title="Debugger" />
    <MemoryK1K2DashboardGroupCharts
      :metrics="metrics"
      :definitions="evaluateExpressionChars"
    />
    <Divider title="Refactoring" />
    <MemoryK1K2DashboardGroupCharts
      :metrics="metrics"
      :definitions="refactoringCharts"
    />
    <Divider title="Script" />
    <MemoryK1K2DashboardGroupCharts
      :metrics="metrics"
      :definitions="scriptCompletionCharts"
    />
    <MemoryK1K2DashboardGroupCharts
      :metrics="metrics"
      :definitions="scriptFindUsagesCharts"
    />
    <MemoryK1K2DashboardGroupCharts
      :metrics="metrics"
      :definitions="highlightingScriptCharts"
    />
    <MemoryK1K2DashboardGroupCharts
      :metrics="metrics"
      :definitions="codeAnalysisScriptCharts"
    />
  </DashboardPage>
</template>

<script setup lang="ts">
import { computed, Ref, ref } from "vue"
import { SimpleMeasureConfigurator } from "../../../configurators/SimpleMeasureConfigurator"
import { metricsSelectLabelFormat } from "../../../shared/labels"
import MeasureSelect from "../../charts/MeasureSelect.vue"
import DashboardPage from "../../common/DashboardPage.vue"
import Divider from "../../common/Divider.vue"
import MemoryK1K2DashboardGroupCharts from "../MemoryK1K2DashboardGroupCharts.vue"
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
  scriptFindUsagesCharts,
} from "../projects"

const measureConfigurator = new SimpleMeasureConfigurator("metrics", null)
measureConfigurator.initData(["freedMemory"])
const metrics = computed(() => {
  const reference = measureConfigurator.selected
  if (reference.value === null) {
    return ref([])
  }
  return reference as Ref<string[]>
})
</script>
