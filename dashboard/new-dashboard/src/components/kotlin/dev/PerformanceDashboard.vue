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
    <KotlinDashboardGroupCharts :definitions="completionCharts" />
    <Divider
      title="Code analysis"
      :description="codeAnalysisChartsDescription"
    />
    <KotlinDashboardGroupCharts :definitions="codeAnalysisCharts" />
    <Divider
      title="Find usages"
      :description="findUsagesChartsDescription"
    />
    <KotlinDashboardGroupCharts :definitions="findUsagesCharts" />
    <Divider
      title="Debugger"
      :description="evaluateExpressionChartsDescription"
    />
    <KotlinDashboardGroupCharts :definitions="evaluateExpressionCharts" />
    <Divider
      title="Refactoring"
      :description="refactoringChartsDescription"
    />
    <KotlinDashboardGroupCharts :definitions="refactoringCharts" />
    <Divider
      title="Code Typing"
      :description="codeTypingChartsDescription"
    />
    <KotlinDashboardGroupCharts :definitions="codeTypingCharts" />
    <Divider
      title="Script"
      :description="scriptChartsDescription"
    />
    <KotlinDashboardGroupCharts :definitions="scriptCompletionCharts" />
    <KotlinDashboardGroupCharts :definitions="codeAnalysisScriptCharts" />
    <KotlinDashboardGroupCharts :definitions="scriptFindUsagesCharts" />
    <Divider
      title="Convert Java to Kotlin"
      :description="convertJavaToKotlinProjectsChartsDescription"
    />
    <KotlinDashboardGroupCharts :definitions="convertJavaToKotlinProjectsCharts" />
  </DashboardPage>
</template>

<script setup lang="ts">
import KotlinDashboardGroupCharts from "../KotlinDashboardGroupCharts.vue"
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
  convertJavaToKotlinProjectsCharts,
  convertJavaToKotlinProjectsChartsDescription,
  evaluateExpressionCharts,
  evaluateExpressionChartsDescription,
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
