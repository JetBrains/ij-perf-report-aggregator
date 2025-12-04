<template>
  <DashboardPage
    db-name="perfintDev"
    table="kmt"
    persistent-id="kmt_performance_dashboard"
    initial-machine="Mac Cidr Performance"
    :initial-mode="MODES"
    :with-installer="false"
  >
    <template #configurator>
      <MeasureSelect
        :configurator="testConfigurator"
        title="Project"
      >
        <template #icon>
          <ChartBarIcon class="w-4 h-4" />
        </template>
      </MeasureSelect>
    </template>
    <Divider title="KMP IDE Setup" />
    <section>
      <!-- :key is used to uniquely identify charts to re-render when definition changes -->
      <GroupProjectsWithClientChart
        v-for="chart in chartsSetup"
        :key="`${chart.definition.label}-${chart.projects.join(',')}`"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
        :aliases="chart.aliases"
        :legend-formatter="legendFormatter"
      />
    </section>
    <Divider title="Swift & Cross-language Support" />
    <section>
      <GroupProjectsWithClientChart
        v-for="chart in chartsCrossLang"
        :key="`${chart.definition.label}-${chart.projects.join(',')}`"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
        :aliases="chart.aliases"
        :legend-formatter="legendFormatter"
      />
    </section>
    <Divider title="Compose Multiplatform Support" />
    <section>
      <GroupProjectsWithClientChart
        v-for="chart in chartsCompose"
        :key="`${chart.definition.label}-${chart.projects.join(',')}`"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
        :aliases="chart.aliases"
        :legend-formatter="legendFormatter"
      />
    </section>
    <Divider title="Run/Debug configurations" />
    <section>
      <GroupProjectsWithClientChart
        v-for="chart in chartsRunConfigurations"
        :key="`${chart.definition.label}-${chart.projects.join(',')}`"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
        :aliases="chart.aliases"
        :legend-formatter="legendFormatter"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        v-for="chart in chartsRunConfigurationsDebug"
        :key="`${chart.definition.label}-${chart.projects.join(',')}`"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
        :aliases="chart.aliases"
        :legend-formatter="legendFormatter"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import DashboardPage from "../common/DashboardPage.vue"
import GroupProjectsWithClientChart from "../charts/GroupProjectsWithClientChart.vue"
import { legendFormatter, MODES, filterChartsByProjects } from "./KmtMeasurements"
import Divider from "../common/Divider.vue"
import { SimpleMeasureConfigurator } from "../../configurators/SimpleMeasureConfigurator"
import MeasureSelect from "../charts/MeasureSelect.vue"
import { computed } from "vue"

const chartsDeclarationSetup: ChartDefinition[] = [
  {
    labels: ["Indexing"],
    measures: [["indexingTimeWithoutPauses", "scanningTimeWithoutPauses"]],
    projects: ["Wizard/indexing", "KotlinConf/indexing"],
  },
  {
    labels: ["KMP Setup"],
    measures: [
      [
        "Progress: Setting up run configurations...",
        "Progress: Generating Xcode files…",
        "Create KMP Run Configurations",
        "Create KMP JS Run Configuration",
        "Create KMP Wasm Run Configuration",
        "Create KMP JVM Run Configuration",
        "Create Xcode Run Configurations",
        "Load XcodeMetaData",
        "Load SwiftPackageManagerWorkspace",
      ],
    ],
    projects: ["Wizard/setup", "KotlinConf/setup"],
  },
]

const chartsDeclarationCrossLang: ChartDefinition[] = [
  {
    labels: ["Swift Inspections (SourceKit)", "Highlighting"],
    measures: [
      ["globalInspections", "firstCodeAnalysis", "localInspections"],
      ["SourceKitDiagnosticsPass", "SourceKitSemanticHighlightingPass", "SwiftDocCommentHighlighter", "SwiftSoftKeywordHighlighter"],
    ],
    projects: ["Wizard/inspection", "KotlinConf/inspection", "Wizard/localInspection/swift", "KotlinConf/localInspection/swift"],
  },
  {
    labels: ["Completion in Swift"],
    measures: ["completion"],
    projects: ["Wizard/completion/kotlinSymbols", "KotlinConf/completion/kotlinSymbols", "Wizard/completion/swiftSymbols", "KotlinConf/completion/swiftSymbols"],
  },
  {
    labels: ["Navigation (Go To Declaration)", "Find Usages"],
    measures: ["execute_editor_gotodeclaration", ["FindUsagesTotal", "execute_editor_findusages"]],
    projects: [
      "Wizard/gotodeclaration/kotlinSymbols",
      "KotlinConf/gotodeclaration/kotlinSymbols",
      "Wizard/gotodeclaration/swiftSymbols",
      "KotlinConf/gotodeclaration/swiftSymbols",
      "Wizard/findUsages/swiftSymbols",
      "KotlinConf/findUsages/swiftSymbols",
    ],
  },
]

const chartsDeclarationCompose: ChartDefinition[] = [
  {
    labels: ["Navigation (Go To Declaration)", "Find Usages"],
    measures: ["execute_editor_gotodeclaration", ["FindUsagesTotal", "execute_editor_findusages"]],
    projects: ["Wizard/gotodeclaration/composeResourceImage", "KotlinConf/findUsages/composeResourceString", "KotlinConf/gotodeclaration/composeResourceString"],
  },
]

const chartsDeclarationRunConfigurations: ChartDefinition[] = [
  {
    labels: ["iOS app - Run with Simulator"],
    measures: [["GradleBuild", "XCodeBuild", "IosAppStartup", "KmpIosConfigurationRun"]],
    projects: ["Wizard/runConfigurationIos/run", "KotlinConf/runConfigurationIos/run"],
  },
]
const chartsDeclarationRunConfigurationsDebug: ChartDefinition[] = [
  {
    labels: ["iOS app - Debug with Simulator"],
    measures: [["GradleBuild", "XCodeBuild", "IosAppStartupDebug", "KmpIosConfigurationRun"]],
    projects: ["Wizard/runConfigurationIos/debug", "KotlinConf/runConfigurationIos/debug"],
  },
]

const chartsDeclaration = chartsDeclarationSetup.concat(chartsDeclarationCrossLang).concat(chartsDeclarationCompose).concat(chartsDeclarationRunConfigurations)
const uniqueProjects: string[] = [...new Set(chartsDeclaration.flatMap((chart) => chart.projects.map((project) => project.split("/")[0])))]
const testConfigurator = new SimpleMeasureConfigurator("project", null)
testConfigurator.initData(uniqueProjects)

const chartsSetup = computed(() => combineCharts(filterChartsByProjects(chartsDeclarationSetup, testConfigurator.selected.value ?? [])))

const chartsCrossLang = computed(() => combineCharts(filterChartsByProjects(chartsDeclarationCrossLang, testConfigurator.selected.value ?? [])))

const chartsCompose = computed(() => combineCharts(filterChartsByProjects(chartsDeclarationCompose, testConfigurator.selected.value ?? [])))

const chartsRunConfigurations = computed(() => combineCharts(filterChartsByProjects(chartsDeclarationRunConfigurations, testConfigurator.selected.value ?? [])))
const chartsRunConfigurationsDebug = computed(() => combineCharts(filterChartsByProjects(chartsDeclarationRunConfigurationsDebug, testConfigurator.selected.value ?? [])))
</script>
