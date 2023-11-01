<template>
  <DashboardPage
    v-slot="{ averagesConfigurators }"
    db-name="perfint"
    table="kotlin"
    persistent-id="kotlin_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'completion\_%'"
          :aggregated-project="'%\_k1'"
          :is-like="true"
          :title="'completion K1'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'completion\_%'"
          :aggregated-project="'%\_k2'"
          :is-like="true"
          :title="'completion K2'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'semanticHighlighting\_%'"
          :aggregated-project="'%\_k1'"
          :is-like="true"
          :title="'semantic highlighting K1'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'semanticHighlighting\_%'"
          :aggregated-project="'%\_k2'"
          :is-like="true"
          :title="'semantic highlighting K2'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'findUsage\_%'"
          :aggregated-project="'%\_k1'"
          :is-like="true"
          :title="'findUsage K1'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'findUsage\_%'"
          :aggregated-project="'%\_k2'"
          :is-like="true"
          :title="'findUsage K2'"
        />
      </div>
    </section>
    <Divider title="Completion" />
    <MemoryDashboardGroupCharts :definitions="completionCharts" />
    <Divider title="Code analysis" />
    <MemoryDashboardGroupCharts :definitions="codeAnalysisCharts" />
    <MemoryDashboardGroupCharts :definitions="highlightingCharts" />
    <Divider title="Find usages" />
    <MemoryDashboardGroupCharts :definitions="findUsagesCharts" />
    <Divider title="Debugger" />
    <MemoryDashboardGroupCharts :definitions="evaluateExpressionChars" />
    <Divider title="Refactoring" />
    <MemoryDashboardGroupCharts :definitions="refactoringCharts" />
  </DashboardPage>
</template>

<script setup lang="ts">
import AggregationChart from "../charts/AggregationChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
import MemoryDashboardGroupCharts from "./MemoryDashboardGroupCharts.vue"
import { completionCharts, evaluateExpressionChars, findUsagesCharts, highlightingCharts, codeAnalysisCharts, refactoringCharts } from "./projects"
</script>
