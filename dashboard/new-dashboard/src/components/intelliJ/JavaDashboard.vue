<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="idea_java_dashboard_devserver"
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
    labels: ["First Code Analysis"],
    measures: ["firstCodeAnalysis"],
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
      "intellij_commit/localInspection/java_file-ContentManagerImpl",
      "intellij_commit/localInspection/java_file",
    ],
  },
  {
    labels: ["Alfio project - Total duration of code analysis (firstCodeAnalysis)"],
    measures: ["firstCodeAnalysis#totalDuration"],
    projects: ["multiple-files-highlighting-integration-test/test-multiple-files-highlighting"],
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
      "intellij_commit/localInspection/java_file-ContentManagerImpl",
      "intellij_commit/localInspection/java_file",
    ],
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
    ],
  },
  {
    labels: ["Highlighting - remove symbol", "Highlighting - remove symbol warmup", "Highlighting - type symbol", "Highlighting - type symbol warmup"],
    measures: ["typing_EditorBackSpace_duration", "typing_EditorBackSpace_warmup_duration", "typing_}_duration", "typing_}_warmup_duration"],
    projects: ["intellij_commit/editor-highlighting"],
  },
  {
    labels: ["Inspection", "Inspection (Full GC Pause)", "Inspection (JVM GC collection times)"],
    measures: ["globalInspections", "fullGCPause", "JVM.GC.collectionTimesMs"],
    projects: ["java/inspection", "grails/inspection", "spring_boot_maven/inspection", "spring_boot/inspection", "intellij_commit/jvm-inspection"],
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
      "grails/completion/groovy_file",
      "grails/completion/java_file",
      "intellij_commit/completion/java_file",
    ],
  },
  {
    labels: ["FindUsages (all usages)", "FindUsages (first usage)", "FindUsages (Full GC Pause)", "FindUsages (JVM GC collection times)"],
    measures: [["findUsages", "fus_find_usages_all"], ["findUsages_firstUsage", "fus_find_usages_first"], ["fullGCPause"], ["JVM.GC.collectionTimesMs"]],
    projects: ["intellij_commit/findUsages/Library_getName", "intellij_commit/findUsages/Application_runReadAction", "intellij_commit/findUsages/String_toString"],
  },
  {
    labels: ["Show Intentions (average awt delay)", "Show Intentions (awt dispatch time)", "Show quick fixes"],
    measures: ["test#average_awt_delay", "AWTEventQueue.dispatchTimeTotal", "showQuickFixes"],
    projects: ["grails/showIntentions/Find cause", "spring_boot/showIntentions"],
  },
  {
    labels: ["Creating a new JAVA file"],
    measures: ["createJavaFile"],
    projects: ["intellij_commit/createJavaClass"],
  },
  {
    labels: ["Rename method, rename class, change signature, move"],
    measures: [["performInlineRename", "changeJavaSignature: add parameter", "moveClassToPackage"]],
    projects: ["hadoop_commit/rename-method", "hadoop_commit/rename-class", "hadoop_commit/change-signature", "hadoop_commit/move-class"],
  },
  {
    labels: ["Inline method"],
    measures: ["inlineJavaMethod"],
    projects: ["hadoop_commit/inline-method"],
  },
  {
    labels: ["Rename package"],
    measures: ["renameDirectoryAsPackage"],
    projects: ["hadoop_commit/rename-package"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
