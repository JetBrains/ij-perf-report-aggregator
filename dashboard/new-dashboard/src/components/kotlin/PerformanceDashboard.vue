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
        :selected-label="projectSelectedLabel"
      >
        <template #icon>
          <ChartBarIcon class="w-4 h-4" />
        </template>
      </MeasureSelect>
    </template>
    <SlackLink></SlackLink>
    <Divider title="Completion" />
    <K1K2DashboardGroupCharts :definitions="completionCharts" />
    <Divider title="Code analysis" />
    <K1K2DashboardGroupCharts :definitions="codeAnalysisCharts" />
    <Divider title="Find usages" />
    <K1K2DashboardGroupCharts :definitions="findUsagesCharts" />
    <Divider title="Debugger" />
    <K1K2DashboardGroupCharts :definitions="evaluateExpressionChars" />
    <Divider title="Refactoring" />
    <K1K2DashboardGroupCharts :definitions="refactoringCharts" />
    <Divider title="Code Typing" />
    <K1K2DashboardGroupCharts :definitions="codeTypingCharts" />
    <Divider title="Script" />
    <K1K2DashboardGroupCharts :definitions="scriptCompletionCharts" />
    <K1K2DashboardGroupCharts :definitions="codeAnalysisScriptCharts" />
    <K1K2DashboardGroupCharts :definitions="scriptFindUsagesCharts" />
    <Divider title="Convert Java to Kotlin" />
    <K1K2DashboardGroupCharts :definitions="convertJavaToKotlinProjectsChars" />
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
import K1K2DashboardGroupCharts from "./K1K2DashboardGroupCharts.vue"
import {
  codeAnalysisCharts,
  codeAnalysisScriptCharts,
  codeTypingCharts,
  completionCharts,
  convertJavaToKotlinProjectsChars,
  evaluateExpressionChars,
  findUsagesCharts,
  KOTLIN_PROJECT_CONFIGURATOR,
  refactoringCharts,
  scriptCompletionCharts,
  scriptFindUsagesCharts,
} from "./projects"
import SlackLink from "./SlackLink.vue"
import MeasureSelect from "../charts/MeasureSelect.vue"
import { projectSelectedLabel } from "./label-formatter"
</script>
