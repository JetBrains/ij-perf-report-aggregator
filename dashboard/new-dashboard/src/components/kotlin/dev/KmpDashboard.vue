<template>
  <DashboardPage
    v-slot="{ averagesConfigurators }"
    db-name="perfintDev"
    table="kotlin"
    persistent-id="kmp_kotlin_dashboard_dev"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
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
    <DashboardGroupCharts :definitions="mppCompletionCharts" />
    <Divider title="Code analysis" />
    <DashboardGroupCharts :definitions="mppCodeAnalysisCharts" />
    <DashboardGroupCharts :definitions="mppHighlightingCharts" />
    <Divider title="Find usages" />
    <DashboardGroupCharts :definitions="mppFindUsagesCharts" />
    <Divider title="Debugger" />
    <DashboardGroupCharts :definitions="mppEvaluateExpressionChars" />
    <Divider title="Refactoring" />
    <DashboardGroupCharts :definitions="mppRefactoringCharts" />
  </DashboardPage>
</template>

<script setup lang="ts">
import AggregationChart from "../../charts/AggregationChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"
import Divider from "../../common/Divider.vue"
import DashboardGroupCharts from "../DashboardGroupCharts.vue"
import { mppCompletionCharts, mppFindUsagesCharts, mppHighlightingCharts, mppCodeAnalysisCharts, mppRefactoringCharts, mppEvaluateExpressionChars } from "../projects"
</script>
