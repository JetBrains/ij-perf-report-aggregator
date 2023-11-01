<template>
  <DashboardPage
    db-name="perfint"
    table="kotlin"
    persistent-id="kotlin_dashboard"
    initial-machine="linux-blade-hetzner"
  >
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
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
import MemoryDashboardGroupCharts from "./MemoryDashboardGroupCharts.vue"
import { completionCharts, evaluateExpressionChars, findUsagesCharts, highlightingCharts, codeAnalysisCharts, refactoringCharts } from "./projects"
</script>
