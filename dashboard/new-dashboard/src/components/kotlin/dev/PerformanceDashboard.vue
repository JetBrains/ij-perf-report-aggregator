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
        :selected-label="projectSelectedLabel"
      >
        <template #icon>
          <ChartBarIcon class="w-4 h-4" />
        </template>
      </MeasureSelect>
    </template>
    <SlackLink></SlackLink>
    <Divider
      title="Completion"
      :description="completionChartsDescription"
    />
    <K1K2DashboardGroupCharts :definitions="completionCharts" />
    <Divider
      title="Code analysis"
      :description="codeAnalysisChartsDescription"
    />
    <K1K2DashboardGroupCharts :definitions="codeAnalysisCharts" />
    <Divider
      title="Find usages"
      :description="findUsagesChartsDescription"
    />
    <K1K2DashboardGroupCharts :definitions="findUsagesCharts" />
    <Divider
      title="Debugger"
      :description="evaluateExpressionCharsDescription"
    />
    <K1K2DashboardGroupCharts :definitions="evaluateExpressionChars" />
    <Divider
      title="Refactoring"
      :description="refactoringChartsDescription"
    />
    <K1K2DashboardGroupCharts :definitions="refactoringCharts" />
    <Divider
      title="Code Typing"
      :description="codeTypingChartsDescription"
    />
    <K1K2DashboardGroupCharts :definitions="codeTypingCharts" />
    <Divider
      title="Script"
      :description="scriptChartsDescription"
    />
    <K1K2DashboardGroupCharts :definitions="scriptCompletionCharts" />
    <K1K2DashboardGroupCharts :definitions="codeAnalysisScriptCharts" />
    <K1K2DashboardGroupCharts :definitions="scriptFindUsagesCharts" />
    <Divider
      title="Convert Java to Kotlin"
      :description="convertJavaToKotlinProjectsCharsDescription"
    />
    <K1K2DashboardGroupCharts :definitions="convertJavaToKotlinProjectsChars" />
  </DashboardPage>
</template>

<script setup lang="ts">
import K1K2DashboardGroupCharts from "../K1K2DashboardGroupCharts.vue"
import DashboardPage from "../../common/DashboardPage.vue"
import Divider from "../../common/Divider.vue"
import {
  codeAnalysisCharts,
  codeAnalysisChartsDescription,
  codeAnalysisScriptCharts,
  codeTypingCharts,
  codeTypingChartsDescription,
  completionCharts,
  completionChartsDescription,
  convertJavaToKotlinProjectsChars,
  convertJavaToKotlinProjectsCharsDescription,
  evaluateExpressionChars,
  evaluateExpressionCharsDescription,
  findUsagesCharts,
  findUsagesChartsDescription,
  KOTLIN_PROJECT_CONFIGURATOR,
  refactoringCharts,
  refactoringChartsDescription,
  scriptCompletionCharts,
  scriptChartsDescription,
  scriptFindUsagesCharts,
} from "../projects"
import SlackLink from "../SlackLink.vue"
import MeasureSelect from "../../charts/MeasureSelect.vue"
import { projectSelectedLabel } from "../label-formatter"
</script>
