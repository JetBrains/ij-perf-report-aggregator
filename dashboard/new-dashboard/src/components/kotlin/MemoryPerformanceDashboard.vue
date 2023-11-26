<template>
  <DashboardPage
    db-name="perfint"
    table="kotlin"
    persistent-id="kotlin_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <MeasureSelect
      title="Metrics"
      :selected-label="metricsSelectLabelFormat"
      :configurator="measureConfigurator"
    ></MeasureSelect>
    <Divider title="Completion" />
    <MemoryDashboardGroupCharts
      :metrics="measureConfigurator.selectedSafe"
      :definitions="completionCharts"
    />
    <Divider title="Code analysis" />
    <MemoryDashboardGroupCharts
      :metrics="measureConfigurator.selectedSafe"
      :definitions="codeAnalysisCharts"
    />
    <MemoryDashboardGroupCharts
      :metrics="measureConfigurator.selectedSafe"
      :definitions="highlightingCharts"
    />
    <Divider title="Find usages" />
    <MemoryDashboardGroupCharts
      :metrics="measureConfigurator.selectedSafe"
      :definitions="findUsagesCharts"
    />
    <Divider title="Debugger" />
    <MemoryDashboardGroupCharts
      :metrics="measureConfigurator.selectedSafe"
      :definitions="evaluateExpressionChars"
    />
    <Divider title="Refactoring" />
    <MemoryDashboardGroupCharts
      :metrics="measureConfigurator.selectedSafe"
      :definitions="refactoringCharts"
    />
    <Divider title="Script" />
    <MemoryDashboardGroupCharts
      :metrics="measureConfigurator.selectedSafe"
      :definitions="scriptCompletionCharts"
    />
    <MemoryDashboardGroupCharts
      :metrics="measureConfigurator.selectedSafe"
      :definitions="highlightingScriptCharts"
    />
    <MemoryDashboardGroupCharts
      :metrics="measureConfigurator.selectedSafe"
      :definitions="codeAnalysisScriptCharts"
    />
  </DashboardPage>
</template>

<script setup lang="ts">
import { provide } from "vue"
import { SimpleMeasureConfigurator } from "../../configurators/SimpleMeasureConfigurator"
import { simpleMeasureConfiguratorKey } from "../../shared/keys"
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

const measureConfigurator = new SimpleMeasureConfigurator("metrics", null)
measureConfigurator.initData(["freedMemoryByGC", "JVM.diffHeapUsageMb_afterGc"])
provide(simpleMeasureConfiguratorKey, measureConfigurator)
</script>
