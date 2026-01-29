<template>
  <DashboardPage
    ref="dashboardPage"
    db-name="perfint"
    table="kotlin"
    persistent-id="kotlin_dashboard"
  >
    <template #configurator>
      <MeasureSelect
        :configurator="projectConfigurator"
        title="Project"
        :selected-label="projectSelectedLabel"
      >
        <template #icon>
          <ChartBarIcon class="w-4 h-4" />
        </template>
      </MeasureSelect>
    </template>
    <ConfiguratorRegistration
      :configurator="projectConfigurator"
      :data="Object.values(PROJECT_CATEGORIES).flatMap((c) => c.label)"
    />
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
      :description="evaluateExpressionChartsDescription"
    />
    <K1K2DashboardGroupCharts :definitions="evaluateExpressionCharts" />
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
      :description="convertJavaToKotlinProjectsChartsDescription"
    />
    <K1K2DashboardGroupCharts :definitions="convertJavaToKotlinProjectsCharts" />
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
import K1K2DashboardGroupCharts from "./K1K2DashboardGroupCharts.vue"
import {
  codeAnalysisChartsDescription,
  codeTypingChartsDescription,
  completionChartsDescription,
  convertJavaToKotlinProjectsChartsDescription,
  createKotlinCharts,
  evaluateExpressionChartsDescription,
  findUsagesChartsDescription,
  PROJECT_CATEGORIES,
  refactoringChartsDescription,
  scriptChartsDescription,
} from "./projects"
import SlackLink from "./SlackLink.vue"
import MeasureSelect from "../charts/MeasureSelect.vue"
import { projectSelectedLabel } from "./label-formatter"
import { SimpleMeasureConfigurator } from "../../configurators/SimpleMeasureConfigurator"
import ConfiguratorRegistration from "./ConfiguratorRegistration.vue"

const projectConfigurator = new SimpleMeasureConfigurator("project", null)

const {
  completionCharts,
  codeAnalysisCharts,
  refactoringCharts,
  codeTypingCharts,
  findUsagesCharts,
  evaluateExpressionCharts,
  convertJavaToKotlinProjectsCharts,
  codeAnalysisScriptCharts,
  scriptCompletionCharts,
  scriptFindUsagesCharts,
} = createKotlinCharts(projectConfigurator)
</script>
