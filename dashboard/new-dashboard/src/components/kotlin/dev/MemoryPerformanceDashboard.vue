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
    <Divider
      title="Completion"
      :description="completionChartsDescription"
    />
    <MemoryKotlinDashboardGroupCharts
      :metrics="metrics"
      :definitions="completionCharts"
    />
    <Divider
      title="Code analysis"
      :description="codeAnalysisChartsDescription"
    />
    <MemoryKotlinDashboardGroupCharts
      :metrics="metrics"
      :definitions="codeAnalysisCharts"
    />
    <Divider
      title="Find usages"
      :description="findUsagesChartsDescription"
    />
    <MemoryKotlinDashboardGroupCharts
      :metrics="metrics"
      :definitions="findUsagesCharts"
    />
    <Divider
      title="Debugger"
      :description="evaluateExpressionChartsDescription"
    />
    <MemoryKotlinDashboardGroupCharts
      :metrics="metrics"
      :definitions="evaluateExpressionCharts"
    />
    <Divider
      title="Refactoring"
      :description="refactoringChartsDescription"
    />
    <MemoryKotlinDashboardGroupCharts
      :metrics="metrics"
      :definitions="refactoringCharts"
    />
    <Divider
      title="Code Typing"
      :description="codeTypingChartsDescription"
    />
    <MemoryKotlinDashboardGroupCharts
      :metrics="metrics"
      :definitions="codeTypingCharts"
    />
    <Divider
      title="Script"
      :description="scriptChartsDescription"
    />
    <MemoryKotlinDashboardGroupCharts
      :metrics="metrics"
      :definitions="scriptCompletionCharts"
    />
    <MemoryKotlinDashboardGroupCharts
      :metrics="metrics"
      :definitions="scriptFindUsagesCharts"
    />
    <MemoryKotlinDashboardGroupCharts
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
import MemoryKotlinDashboardGroupCharts from "../MemoryKotlinDashboardGroupCharts.vue"
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
