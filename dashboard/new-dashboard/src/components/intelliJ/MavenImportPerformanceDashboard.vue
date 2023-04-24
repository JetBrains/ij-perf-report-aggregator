<template>
  <div class="flex flex-col gap-5">
    <Toolbar class="customToolbar">
      <template #start>
        <TimeRangeSelect
          :ranges="TimeRangeConfigurator.timeRanges"
          :value="timeRangeConfigurator.value.value"
          :on-change="onChangeRange"
        >
          <template #icon>
            <CalendarIcon class="w-4 h-4 text-gray-500" />
          </template>
        </TimeRangeSelect>
        <BranchSelect
          :branch-configurator="branchConfigurator"
          :release-configurator="releaseConfigurator"
          :triggered-by-configurator="triggeredByConfigurator"
        />
        <DimensionHierarchicalSelect
          label="Machine"
          :dimension="machineConfigurator"
        >
          <template #icon>
            <ComputerDesktopIcon class="w-4 h-4 text-gray-500" />
          </template>
        </DimensionHierarchicalSelect>
      </template>
    </Toolbar>

    <main class="flex">
      <div
        ref="container"
        class="flex flex-1 flex-col gap-6 overflow-hidden"
      >
        <section>
          <GroupProjectsChart
            v-for="chart in charts"
            :key="chart.definition.label"
            :label="chart.definition.label"
            :measure="chart.definition.measure"
            :projects="chart.projects"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
      </div>
      <InfoSidebar />
    </main>
  </div>
</template>

<script setup lang="ts">
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import DimensionHierarchicalSelect from "shared/src/components/DimensionHierarchicalSelect.vue"
import { createBranchConfigurator } from "shared/src/configurators/BranchConfigurator"
import { dimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { privateBuildConfigurator } from "shared/src/configurators/PrivateBuildConfigurator"
import { ReleaseNightlyConfigurator } from "shared/src/configurators/ReleaseNightlyConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRange, TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { refToObservable } from "shared/src/configurators/rxjs"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
import { Accident, getAccidentsFromMetaDb } from "shared/src/meta"
import { provide, ref } from "vue"
import { useRouter } from "vue-router"
import { containerKey, sidebarVmKey } from "../../shared/keys"
import InfoSidebar from "../InfoSidebar.vue"
import { InfoSidebarVmImpl } from "../InfoSidebarVm"
import { ChartDefinition, combineCharts, extractUniqueProjects } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import BranchSelect from "../common/BranchSelect.vue"
import TimeRangeSelect from "../common/TimeRangeSelect.vue"

provideReportUrlProvider()

const dbName = "perfint"
const dbTable = "idea"
const initialMachine = "Linux EC2 C6i.8xlarge (32 vCPU Xeon, 64 GB)"
const container = ref<HTMLElement>()
const router = useRouter()
const sidebarVm = new InfoSidebarVmImpl()

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const metricsName = ["after_apply_duration_ms", "before_apply_duration_ms", "collect_folders_duration_ms", "config_modules_duration_ms", "total_duration_ms"]
const metricsDeclaration = ["CPU | Load |Total % 95th pctl", "Memory | IDE | RESIDENT SIZE (MB) 95th pctl", "Memory | IDE | VIRTUAL SIZE (MB) 95th pctl",
  "gcPause", "gcPauseCount", "fullGCPause", "freedMemoryByGC", "totalHeapUsedMax", "maven.sync.duration", "maven.import.after.import.configuration",
  "maven.import.stats.applying.model.task", "maven.import.stats.importing.task", "maven.import.stats.importing.task.old", "maven.project.importer.base.refreshing.files.task",
  "maven.projects.processor.plugin.resolving.task", "maven.projects.processor.reading.task", "maven.projects.processor.resolving.task",
  "maven.projects.processor.wait.for.completion.task", "quarkus.maven.importer.base.task", "quarkus.maven.post.processor.task", "jpa.facet.importer.lambdas",
  "jpa.facet.importer.reimport.facet", "importer_run.com.intellij.jpa.importer.maven.JpaFacetImporter.total_duration_ms",
  "importer_run.com.intellij.quarkus.run.maven.QsMavenImporter.total_duration_ms", "importer_run.com.intellij.spring.facet.importer.maven.SpringFacetImporter.total_duration_ms",
  "importer_run.com.intellij.spring.mvc.importer.boot.SpringBootWebFacetImporter.total_duration_ms",
  "importer_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporter.total_duration_ms",
  "importer_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporter.total_duration_ms",
  "importer_run.org.jetbrains.idea.maven.importing.MavenAnnotationProcessorConfigurator.total_duration_ms",
  "importer_run.org.jetbrains.idea.maven.importing.MavenCompilerConfigurator.total_duration_ms",
  "importer_run.org.jetbrains.idea.maven.importing.MavenEncodingConfigurator.total_duration_ms",
  "importer_run.org.jetbrains.idea.maven.importing.MavenExternalAnnotationsConfigurator.total_duration_ms",
  "importer_run.org.jetbrains.idea.maven.importing.MavenRemoteRepositoriesConfigurator.total_duration_ms",
  "importer_run.org.jetbrains.kotlin.idea.maven.KotlinMavenImporter.total_duration_ms", "legacy_import.create_modules.duration_ms", "legacy_import.delete_obsolete.duration_ms",
  "legacy_import.duration_ms", "legacy_import.importers.duration_ms", "workspace_commit.attempts", "workspace_commit.duration_in_background_ms",
  "workspace_commit.duration_in_write_action_ms", "workspace_commit.duration_of_workspace_update_call_ms", "workspace_import.commit.duration_ms",
  ...metricsName.map(metric => "workspace_import.configurator_run.com.intellij.spring.facet.importer.maven.SpringFacetImporter." + metric),
  ...metricsName.map(metric => "workspace_import.configurator_run.com.intellij.spring.mvc.importer.boot.SpringBootWebFacetImporter." + metric),
  ...metricsName.map(metric => "workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenAnnotationProcessorConfigurator." + metric),
  ...metricsName.map(metric => "workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenCompilerConfigurator." + metric),
  ...metricsName.map(metric => "workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenEncodingConfigurator." + metric),
  ...metricsName.map(metric => "workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenExternalAnnotationsConfigurator." + metric),
  ...metricsName.map(metric => "workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenRemoteRepositoriesConfigurator." + metric),
  ...metricsName.map(metric => "workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenWslTargetConfigurator." + metric),
  ...metricsName.map(metric => "workspace_import.configurator_run.org.jetbrains.idea.maven.plugins.groovy.GroovyPluginConfigurator." + metric),
  ...metricsName.map(metric => "workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporter" + metric),
  ...metricsName.map(metric => "workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporter." + metric),
  ...metricsName.map(metric => "workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporterEx." + metric),
  ...metricsName.map(metric => "workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporterEx." + metric),
  "workspace_import.duration_ms", "workspace_import.legacy_importers.duration_ms", "workspace_import.legacy_importers.stats.duration_of_bridges_creation_ms",
  "workspace_import.legacy_importers.stats.duration_of_bridges_commit_ms", "workspace_import.populate.duration_ms", "maven.project.importer.post.importing.task.marker",
  "post_import_tasks_run.total_duration_ms",
]


const chartsDeclaration: Array<ChartDefinition> = metricsDeclaration.map(metric => {
  return {
    labels: [metric],
    measures: [metric],
    projects: [
      "project-import-maven-quarkus/measureStartup",
      "project-import-maven-500-modules/measureStartup", "project-import-maven-1000-modules/measureStartup",
      "project-import-maven-keycloak/measureStartup", "project-import-maven-javaee7/measureStartup",
      "project-import-maven-javaee8/measureStartup", "project-import-maven-jersey/measureStartup",
      "project-import-maven-flink/measureStartup", "project-import-maven-drill/measureStartup",
      "project-import-maven-azure-sdk-java/measureStartup", "project-import-maven-hive/measureStartup",
      "project-import-maven-quarkus-to-legacy-model/measureStartup", "project-import-maven-500-modules-to-legacy-model/measureStartup",
      "project-import-maven-1000-modules-to-legacy-model/measureStartup",
    ],
  }
})
const charts = combineCharts(chartsDeclaration)


const serverConfigurator = new ServerConfigurator(dbName, dbTable)

const persistenceForDashboard = new PersistentStateManager("idea_gradle_dashboard", {
  machine: initialMachine,
  project: [],
  branch: "master",
}, router)

const timeRangeConfigurator = new TimeRangeConfigurator(persistenceForDashboard)

const scenarioConfigurator = dimensionConfigurator(
  "project",
  serverConfigurator,
  null,
  true,
  [timeRangeConfigurator]
)
scenarioConfigurator.selected.value = extractUniqueProjects(chartsDeclaration)

const branchConfigurator = createBranchConfigurator(serverConfigurator, persistenceForDashboard, [timeRangeConfigurator, scenarioConfigurator])
const machineConfigurator = new MachineConfigurator(
  serverConfigurator,
  persistenceForDashboard,
  [timeRangeConfigurator, branchConfigurator, scenarioConfigurator],
)
const releaseConfigurator = new ReleaseNightlyConfigurator(persistenceForDashboard)
const triggeredByConfigurator = privateBuildConfigurator(
  serverConfigurator,
  persistenceForDashboard,
  [branchConfigurator, timeRangeConfigurator, scenarioConfigurator],
)
const dashboardConfigurators = [
  branchConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  releaseConfigurator,
  triggeredByConfigurator,
]

function onChangeRange(value: string) {
  timeRangeConfigurator.value.value = value
}
const warnings = ref<Array<Accident>>()
refToObservable(timeRangeConfigurator.value).subscribe(data => {
  getAccidentsFromMetaDb(warnings, null, data as TimeRange)
})
</script>

<style>
.customToolbar {
  background-color: transparent;
  border: none;
  padding: 0;
}
</style>