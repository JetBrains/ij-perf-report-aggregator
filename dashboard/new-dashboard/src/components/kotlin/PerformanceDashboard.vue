<template>
  <DashboardPage
    db-name="perfint"
    table="kotlin"
    persistent-id="kotlin_dashboard"
  >
    <template #configurator>
      <MeasureSelect
        :configurator="KOTLIN_PROJECT_CONFIGURATOR"
        title="Project"
      >
        <template #icon>
          <ChartBarIcon class="w-4 h-4 text-gray-500" />
        </template>
      </MeasureSelect>
    </template>
    <SlackLink></SlackLink>
    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <AggregationChart
          :aggregated-measure="'completion\_%'"
          :aggregated-project="'%\_k1'"
          :is-like="true"
          :title="'completion K1'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :aggregated-measure="'completion\_%'"
          :aggregated-project="'%\_k2'"
          :is-like="true"
          :title="'completion K2'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :aggregated-measure="'semanticHighlighting\_%'"
          :aggregated-project="'%\_k1'"
          :is-like="true"
          :title="'semantic highlighting K1'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :aggregated-measure="'semanticHighlighting\_%'"
          :aggregated-project="'%\_k2'"
          :is-like="true"
          :title="'semantic highlighting K2'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :aggregated-measure="'findUsage\_%'"
          :aggregated-project="'%\_k1'"
          :is-like="true"
          :title="'findUsage K1'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :aggregated-measure="'findUsage\_%'"
          :aggregated-project="'%\_k2'"
          :is-like="true"
          :title="'findUsage K2'"
        />
      </div>
    </section>
    <Divider title="Completion" />
    <DashboardGroupCharts :definitions="completionCharts" />
    <Divider title="Code analysis" />
    <DashboardGroupCharts :definitions="codeAnalysisCharts" />
    <DashboardGroupCharts :definitions="highlightingCharts" />
    <Divider title="Find usages" />
    <DashboardGroupCharts :definitions="findUsagesCharts" />
    <Divider title="Debugger" />
    <DashboardGroupCharts :definitions="evaluateExpressionChars" />
    <Divider title="Refactoring" />
    <DashboardGroupCharts :definitions="refactoringCharts" />
    <Divider title="Script" />
    <DashboardGroupCharts :definitions="scriptCompletionCharts" />
    <DashboardGroupCharts :definitions="highlightingScriptCharts" />
    <DashboardGroupCharts :definitions="codeAnalysisScriptCharts" />
    <DashboardGroupCharts :definitions="scriptFindUsagesCharts" />
    <Divider title="Convert Java to Kotlin" />
    <DashboardGroupCharts :definitions="convertJavaToKotlinProjectsChars" />
  </DashboardPage>
</template>

<script setup lang="ts">
import AggregationChart from "../charts/AggregationChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
import DashboardGroupCharts from "./DashboardGroupCharts.vue"
import {
  codeAnalysisCharts,
  codeAnalysisScriptCharts,
  completionCharts,
  convertJavaToKotlinProjectsChars,
  evaluateExpressionChars,
  findUsagesCharts,
  highlightingCharts,
  highlightingScriptCharts,
  KOTLIN_PROJECT_CONFIGURATOR,
  refactoringCharts,
  scriptCompletionCharts,
  scriptFindUsagesCharts,
} from "./projects"
import SlackLink from "./SlackLink.vue"
import MeasureSelect from "../charts/MeasureSelect.vue"
</script>
