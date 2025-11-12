<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="idea_ultimate_dashboard_devserver"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :charts="charts"
    :with-installer="false"
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
  </DashboardPage>
</template>

<script setup lang="ts">
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["Indexing", "Scanning"],
    measures: ["indexingTimeWithoutPauses", "scanningTimeWithoutPauses"],
    projects: ["swagger_indexing/indexing"],
  },
  {
    labels: ["First Code Analysis", "File Openings: code loaded", "File Openings: tab shown"],
    measures: ["firstCodeAnalysis", "fus_file_types_usage_duration_ms", "fus_file_types_usage_time_to_show_ms"],
    projects: [
      "keycloak_release_20/ultimateCase/AuthenticationManagementResource",
      "keycloak_release_20/ultimateCase/IdentityBrokerService",
      "keycloak_release_20/ultimateCase/JpaUserProvider",
      "keycloak_release_20/ultimateCase/RealmAdminResource",
      "keycloak_release_20/localInspection/QuarkusRuntimePomXml",
      "keycloak_release_20/localInspection/RootPomXml",
      "keycloak_release_20/localInspection/CorePomXml",
      "keycloak_release_20/typing/ClientEntity",
      "keycloak_release_20/typing/PolicyEntity",
      "train-ticket/ultimateCase/AdminBasicInfoController",
      "train-ticket/ultimateCase/ExecuteServiceImpl",
      "train-ticket/ultimateCase/InsidePaymentServiceImpl",
      "train-ticket/ultimateCase/OrderController",
      "train-ticket/ultimateCase/OrderServiceImpl",
      "toolbox_enterprise/ultimateCase/SecurityTests",
      "toolbox_enterprise/ultimateCase/ToolController",
      "toolbox_enterprise/ultimateCase/ToolService",
      "toolbox_enterprise/ultimateCase/UserController",
      "toolbox_enterprise/ultimateCase/UserRepository",
      "json_azure/typing/openAndType",
      "json_azure/completion",
      "json_schema_modes_comparison/localInspection/NewSchema",
      "json_schema_modes_comparison/localInspection/OldSchema",
    ],
  },
  {
    labels: ["Local Inspection"],
    measures: ["localInspections"],
    projects: [
      "keycloak_release_20/ultimateCase/AuthenticationManagementResource",
      "keycloak_release_20/ultimateCase/IdentityBrokerService",
      "keycloak_release_20/ultimateCase/JpaUserProvider",
      "keycloak_release_20/ultimateCase/RealmAdminResource",
      "keycloak_release_20/localInspection/QuarkusRuntimePomXml",
      "keycloak_release_20/localInspection/RootPomXml",
      "keycloak_release_20/localInspection/CorePomXml",
      "train-ticket/ultimateCase/AdminBasicInfoController",
      "train-ticket/ultimateCase/ExecuteServiceImpl",
      "train-ticket/ultimateCase/InsidePaymentServiceImpl",
      "train-ticket/ultimateCase/OrderController",
      "train-ticket/ultimateCase/OrderServiceImpl",
      "toolbox_enterprise/ultimateCase/SecurityTests",
      "toolbox_enterprise/ultimateCase/ToolController",
      "toolbox_enterprise/ultimateCase/ToolService",
      "toolbox_enterprise/ultimateCase/UserController",
      "toolbox_enterprise/ultimateCase/UserRepository",
      "json_schema_modes_comparison/localInspection/NewSchema",
      "json_schema_modes_comparison/localInspection/OldSchema",
    ],
  },
  {
    labels: ["Total duration of code analysis (firstCodeAnalysis)"],
    measures: ["firstCodeAnalysis#totalDuration"],
    projects: ["multiple-files-highlighting-integration-test/test-multiple-files-highlighting"],
  },
  {
    labels: ["Typing (typingCodeAnalyzing)", "Typing (average_awt_delay)", "Typing (max_awt_delay)"],
    measures: ["typingCodeAnalyzing", "test#average_awt_delay", "test#max_awt_delay"],
    projects: [
      "keycloak_release_20/typing/ClientEntity",
      "keycloak_release_20/typing/PolicyEntity",
      "keycloak_release_20/ultimateCase/AuthenticationManagementResource",
      "keycloak_release_20/ultimateCase/IdentityBrokerService",
      "keycloak_release_20/ultimateCase/JpaUserProvider",
      "keycloak_release_20/ultimateCase/RealmAdminResource",
      "train-ticket/ultimateCase/AdminBasicInfoController",
      "train-ticket/ultimateCase/ExecuteServiceImpl",
      "train-ticket/ultimateCase/InsidePaymentServiceImpl",
      "train-ticket/ultimateCase/OrderController",
      "train-ticket/ultimateCase/OrderServiceImpl",
      "toolbox_enterprise/ultimateCase/SecurityTests",
      "toolbox_enterprise/ultimateCase/ToolController",
      "toolbox_enterprise/ultimateCase/ToolService",
      "toolbox_enterprise/ultimateCase/UserController",
      "toolbox_enterprise/ultimateCase/UserRepository",
      "json_azure/typing/openAndType",
    ],
  },
  {
    labels: ["Completion", "Completion 90p", "Completion time to show 90p"],
    measures: ["completion", "fus_completion_duration_90p", "fus_time_to_show_90p"],
    projects: [
      "keycloak_release_20/ultimateCase/AuthenticationManagementResource",
      "keycloak_release_20/ultimateCase/IdentityBrokerService",
      "keycloak_release_20/ultimateCase/JpaUserProvider",
      "keycloak_release_20/ultimateCase/RealmAdminResource",
      "keycloak_release_20/completion/QuarkusRuntimePomXml",
      "keycloak_release_20/completion/RootPomXml",
      "keycloak_release_20/completion/CorePomXml",
      "train-ticket/ultimateCase/AdminBasicInfoController",
      "train-ticket/ultimateCase/ExecuteServiceImpl",
      "train-ticket/ultimateCase/InsidePaymentServiceImpl",
      "train-ticket/ultimateCase/OrderController",
      "train-ticket/ultimateCase/OrderServiceImpl",
      "toolbox_enterprise/ultimateCase/SecurityTests",
      "toolbox_enterprise/ultimateCase/ToolController",
      "toolbox_enterprise/ultimateCase/ToolService",
      "toolbox_enterprise/ultimateCase/UserController",
      "toolbox_enterprise/ultimateCase/UserRepository",
      "json_schema_modes_comparison/completion/OldSchema",
      "json_schema_modes_comparison/completion/NewSchema",
      "json_azure/completion",
    ],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
