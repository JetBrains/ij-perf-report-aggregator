<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="idea_ultimate_dashboard"
    initial-machine="Linux EC2 C6i.8xlarge (32 vCPU Xeon, 64 GB)"
    :charts="charts"
  >
    <section>
      <GroupProjectsChart
        v-for="chart in charts"
        :key="chart.definition.label"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
      />
    </section>
  </DashboardPage>>
</template>

<script setup lang="ts">
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const chartsDeclaration: ChartDefinition[] = [{
  labels: ["Indexing", "Scanning", "Number of indexed files"],
  measures: ["indexing", "scanning", "numberOfIndexedFiles"],
  projects: ["keycloak_release_20/indexing", "train-ticket/indexing", "toolbox_enterprise/indexing"],
}, {
  labels: ["Local Inspection", "First Code Analysis"],
  measures: ["localInspections", "firstCodeAnalysis"],
  projects: ["keycloak_release_20/ultimateCase/AuthenticationManagementResource", "keycloak_release_20/ultimateCase/IdentityBrokerService",
    "keycloak_release_20/ultimateCase/JpaUserProvider", "keycloak_release_20/ultimateCase/RealmAdminResource",
    "keycloak_release_20/localInspection/QuarkusRuntimePomXml", "keycloak_release_20/localInspection/RootPomXml",
    "keycloak_release_20/localInspection/CorePomXml", "train-ticket/ultimateCase/AdminBasicInfoController",
    "train-ticket/ultimateCase/ExecuteServiceImpl", "train-ticket/ultimateCase/InsidePaymentServiceImpl",
    "train-ticket/ultimateCase/OrderController", "train-ticket/ultimateCase/OrderServiceImpl",
    "toolbox_enterprise/ultimateCase/SecurityTests", "toolbox_enterprise/ultimateCase/ToolController",
    "toolbox_enterprise/ultimateCase/ToolService", "toolbox_enterprise/ultimateCase/UserController",
    "toolbox_enterprise/ultimateCase/UserRepository"],
}, {
  labels: ["Completion"],
  measures: ["completion"],
  projects: ["keycloak_release_20/ultimateCase/AuthenticationManagementResource", "keycloak_release_20/ultimateCase/IdentityBrokerService",
    "keycloak_release_20/ultimateCase/JpaUserProvider", "keycloak_release_20/ultimateCase/RealmAdminResource",
    "keycloak_release_20/completion/QuarkusRuntimePomXml", "keycloak_release_20/completion/RootPomXml",
    "keycloak_release_20/completion/CorePomXml", "train-ticket/ultimateCase/AdminBasicInfoController",
    "train-ticket/ultimateCase/ExecuteServiceImpl", "train-ticket/ultimateCase/InsidePaymentServiceImpl",
    "train-ticket/ultimateCase/OrderController", "train-ticket/ultimateCase/OrderServiceImpl",
    "toolbox_enterprise/ultimateCase/SecurityTests", "toolbox_enterprise/ultimateCase/ToolController",
    "toolbox_enterprise/ultimateCase/ToolService", "toolbox_enterprise/ultimateCase/UserController",
    "toolbox_enterprise/ultimateCase/UserRepository"],
}, {
  labels: ["Typing (firstCodeAnalysis)", "Typing (typingCodeAnalyzing)", "Typing (average_awt_delay)", "Typing (max_awt_delay)"],
  measures: ["firstCodeAnalysis", "typingCodeAnalyzing", "test#average_awt_delay", "test#max_awt_delay"],
  projects: ["keycloak_release_20/typing/ClientEntity", "keycloak_release_20/typing/PolicyEntity",
    "keycloak_release_20/ultimateCase/AuthenticationManagementResource", "keycloak_release_20/ultimateCase/IdentityBrokerService",
    "keycloak_release_20/ultimateCase/JpaUserProvider", "keycloak_release_20/ultimateCase/RealmAdminResource",
    "train-ticket/ultimateCase/AdminBasicInfoController",
    "train-ticket/ultimateCase/ExecuteServiceImpl", "train-ticket/ultimateCase/InsidePaymentServiceImpl",
    "train-ticket/ultimateCase/OrderController", "train-ticket/ultimateCase/OrderServiceImpl",
    "toolbox_enterprise/ultimateCase/SecurityTests", "toolbox_enterprise/ultimateCase/ToolController",
    "toolbox_enterprise/ultimateCase/ToolService", "toolbox_enterprise/ultimateCase/UserController",
    "toolbox_enterprise/ultimateCase/UserRepository"],
}]

const charts = combineCharts(chartsDeclaration)
</script>