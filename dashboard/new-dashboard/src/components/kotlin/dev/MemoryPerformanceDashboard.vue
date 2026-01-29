<template>
  <DashboardPage
    db-name="perfintDev"
    table="kotlin"
    persistent-id="kotlin_memory_dashboard_dev"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  >
    <ConfiguratorRegistration
      :configurator="projectConfigurator"
      :data="Object.values(PROJECT_CATEGORIES).flatMap((c) => c.label)"
    />
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
import { SimpleMeasureConfigurator } from "../../../configurators/SimpleMeasureConfigurator"
import { metricsSelectLabelFormat } from "../../../shared/labels"
import MeasureSelect from "../../charts/MeasureSelect.vue"
import DashboardPage from "../../common/DashboardPage.vue"
import Divider from "../../common/Divider.vue"
import MemoryK1K2DashboardGroupCharts from "../MemoryK1K2DashboardGroupCharts.vue"
import {
  createKotlinCharts,
  PROJECT_CATEGORIES,
  completionChartsDescription,
  evaluateExpressionChartsDescription,
  findUsagesChartsDescription,
  codeAnalysisChartsDescription,
  refactoringChartsDescription,
  scriptChartsDescription,
  codeTypingChartsDescription,
} from "../projects"
import ConfiguratorRegistration from "../ConfiguratorRegistration.vue"

const projectConfigurator = new SimpleMeasureConfigurator("project", null)
const {
  completionCharts,
  codeAnalysisCharts,
  refactoringCharts,
  codeTypingCharts,
  findUsagesCharts,
  evaluateExpressionCharts,
  scriptCompletionCharts,
  codeAnalysisScriptCharts,
  scriptFindUsagesCharts,
} = createKotlinCharts(projectConfigurator)

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
