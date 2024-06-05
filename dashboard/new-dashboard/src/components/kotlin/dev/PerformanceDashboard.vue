<template>
  <DashboardPage
    db-name="perfintDev"
    table="kotlin"
    persistent-id="kotlin_dashboard_dev"
    :with-installer="false"
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
    <Divider title="Completion" />
    <K1K2DashboardGroupCharts :definitions="completionCharts" />
    <Divider title="Code analysis" />
    <K1K2DashboardGroupCharts :definitions="codeAnalysisCharts" />
    <K1K2DashboardGroupCharts :definitions="highlightingCharts" />
    <Divider title="Find usages" />
    <K1K2DashboardGroupCharts :definitions="findUsagesCharts" />
    <Divider title="Debugger" />
    <K1K2DashboardGroupCharts :definitions="evaluateExpressionChars" />
    <Divider title="Refactoring" />
    <K1K2DashboardGroupCharts :definitions="refactoringCharts" />
    <Divider title="Script" />
    <K1K2DashboardGroupCharts :definitions="scriptCompletionCharts" />
    <K1K2DashboardGroupCharts :definitions="highlightingScriptCharts" />
    <K1K2DashboardGroupCharts :definitions="codeAnalysisScriptCharts" />
    <K1K2DashboardGroupCharts :definitions="scriptFindUsagesCharts" />
    <Divider title="Convert Java to Kotlin" />
    <K1K2DashboardGroupCharts :definitions="convertJavaToKotlinProjectsChars" />
  </DashboardPage>
</template>

<script setup lang="ts">
import K1K2DashboardGroupCharts from "../K1K2DashboardGroupCharts.vue"
import DashboardPage from "../../common/DashboardPage.vue"
import Divider from "../../common/Divider.vue"
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
} from "../projects"
import SlackLink from "../SlackLink.vue"
import MeasureSelect from "../../charts/MeasureSelect.vue"
</script>
