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
          <ChartBarIcon class="w-4 h-4" />
        </template>
      </MeasureSelect>
    </template>
    <Divider
      title="Completion"
      :description="completionChartsDescription"
    />
    <MemoryK1K2DashboardGroupCharts
      :metrics="metrics"
      :definitions="completionCharts"
    />
    <Divider
      title="Code analysis"
      :description="codeAnalysisChartsDescription"
    />
    <MemoryK1K2DashboardGroupCharts
      :metrics="metrics"
      :definitions="codeAnalysisCharts"
    />
    <Divider
      title="Find usages"
      :description="findUsagesChartsDescription"
    />
    <MemoryK1K2DashboardGroupCharts
      :metrics="metrics"
      :definitions="findUsagesCharts"
    />
    <Divider
      title="Debugger"
      :description="evaluateExpressionChartsDescription"
    />
    <MemoryK1K2DashboardGroupCharts
      :metrics="metrics"
      :definitions="evaluateExpressionCharts"
    />
    <Divider
      title="Refactoring"
      :description="refactoringChartsDescription"
    />
    <MemoryK1K2DashboardGroupCharts
      :metrics="metrics"
      :definitions="refactoringCharts"
    />
    <Divider
      title="Code Typing"
      :description="codeTypingChartsDescription"
    />
    <MemoryK1K2DashboardGroupCharts
      :metrics="metrics"
      :definitions="codeTypingCharts"
    />
    <Divider
      title="Script"
      :description="scriptChartsDescription"
    />
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
      :definitions="codeAnalysisScriptCharts"
    />
  </DashboardPage>
</template>

<script setup lang="ts">
import { computed, Ref, ref } from "vue"
import { SimpleMeasureConfigurator } from "../../configurators/SimpleMeasureConfigurator"
import { metricsSelectLabelFormat } from "../../shared/labels"
import MeasureSelect from "../charts/MeasureSelect.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
import MemoryK1K2DashboardGroupCharts from "./MemoryK1K2DashboardGroupCharts.vue"
import {
  completionCharts,
  completionChartsDescription,
  evaluateExpressionCharts,
  evaluateExpressionChartsDescription,
  findUsagesCharts,
  findUsagesChartsDescription,
  codeAnalysisCharts,
  codeAnalysisChartsDescription,
  refactoringCharts,
  refactoringChartsDescription,
  scriptCompletionCharts,
  scriptChartsDescription,
  codeAnalysisScriptCharts,
  scriptFindUsagesCharts,
  codeTypingCharts,
  codeTypingChartsDescription,
} from "./projects"

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
