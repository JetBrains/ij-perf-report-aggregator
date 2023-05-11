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

const metricsDeclaration = [
  "project.opening",
  "jps.app.storage.content.reader.load.component.ms",
  "jps.app.storage.content.writer.save.component.ms",
  "jps.apply.loaded.storage.ms",
  "jps.artifact.entities.serializer.load.entities.ms",
  "jps.artifact.entities.serializer.save.entities.ms",
  "jps.facet.change.listener.before.change.events.ms",
  "jps.facet.change.listener.init.bridge.ms",
  "jps.facet.change.listener.process.change.events.ms",
  "jps.global.get.libraries.ms",
  "jps.global.get.library.by.name.ms",
  "jps.global.get.library.ms",
  "jps.global.handle.before.change.events.ms",
  "jps.global.handle.changed.events.ms",
  "jps.global.initialize.library.bridges.after.loading.ms",
  "jps.global.initialize.library.bridges.ms",
  "jps.library.entities.serializer.load.entities.ms",
  "jps.library.entities.serializer.save.entities.ms",
  "jps.load.initial.state.ms",
  "jps.load.project.to.empty.storage.ms",
  "jps.module.iml.entities.serializer.load.entities.ms",
  "jps.module.iml.entities.serializer.save.entities.ms",
  "jps.project.serializers.load.ms",
  "jps.project.serializers.save.ms",
  "jps.reload.project.entities.ms",
  "jps.save.changed.project.entities.ms",
  "jps.save.global.entities.ms",
  "jps.storage.jps.conf.reader.load.component.ms",
  "workspaceModel.check.recursive.update.ms",
  "workspaceModel.collect.changes.ms",
  "workspaceModel.delayed.project.synchronizer.sync.ms",
  "workspaceModel.global.apply.state.to.project.builder.ms",
  "workspaceModel.global.apply.state.to.project.ms",
  "workspaceModel.global.updates.count",
  "workspaceModel.global.updates.ms",
  "workspaceModel.init.bridges.ms",
  "workspaceModel.initializing.ms",
  "workspaceModel.load.cache.from.file.ms",
  "workspaceModel.loading.from.cache.ms",
  "workspaceModel.loading.total.ms",
  "workspaceModel.moduleBridge.before.changed.ms",
  "workspaceModel.moduleBridge.facet.initialization.ms",
  "workspaceModel.moduleBridge.update.option.ms",
  "workspaceModel.moduleBridgeLoader.loading.modules.ms",
  "workspaceModel.moduleManagerBridge.build.module.graph.ms",
  "workspaceModel.moduleManagerBridge.create.module.instance.ms",
  "workspaceModel.moduleManagerBridge.get.modules.ms",
  "workspaceModel.moduleManagerBridge.load.all.modules.ms",
  "workspaceModel.moduleManagerBridge.load.module.ms",
  "workspaceModel.moduleManagerBridge.new.nonPersistent.module.ms",
  "workspaceModel.moduleManagerBridge.newModule.ms",
  "workspaceModel.moduleManagerBridge.set.unloadedModules.ms",
  "workspaceModel.mutableEntityStorage.add.diff.ms",
  "workspaceModel.mutableEntityStorage.add.entity.ms",
  "workspaceModel.mutableEntityStorage.collect.changes.ms",
  "workspaceModel.mutableEntityStorage.entities.by.source.ms",
  "workspaceModel.mutableEntityStorage.entities.ms",
  "workspaceModel.mutableEntityStorage.has.same.entities.ms",
  "workspaceModel.mutableEntityStorage.modify.entity.ms",
  "workspaceModel.mutableEntityStorage.mutable.ext.mapping.ms",
  "workspaceModel.mutableEntityStorage.mutable.vfurl.index.ms",
  "workspaceModel.mutableEntityStorage.put.entity.ms",
  "workspaceModel.mutableEntityStorage.referrers.ms",
  "workspaceModel.mutableEntityStorage.remove.entity.ms",
  "workspaceModel.mutableEntityStorage.replace.by.source.ms",
  "workspaceModel.mutableEntityStorage.resolve.ms",
  "workspaceModel.mutableEntityStorage.to.snapshot.ms",
  "workspaceModel.orphan.listener.update.ms",
  "workspaceModel.pre.handlers.ms",
  "workspaceModel.replace.project.model.ms",
  "workspaceModel.save.cache.to.file.ms",
  "workspaceModel.sync.entities.ms",
  "workspaceModel.to.snapshot.ms",
  "workspaceModel.update.unloaded.entities.ms",
  "workspaceModel.updates.count",
  "workspaceModel.updates.ms",
  "workspaceModel.updates.precise.ms",

  "CPU | Load |Total % 95th pctl", "Memory | IDE | RESIDENT SIZE (MB) 95th pctl", "Memory | IDE | VIRTUAL SIZE (MB) 95th pctl",
  "gcPause", "gcPauseCount", "fullGCPause", "freedMemoryByGC", "totalHeapUsedMax",
]


const chartsDeclaration: Array<ChartDefinition> = metricsDeclaration.map(metric => {
  return {
    labels: [metric],
    measures: [metric],
    projects: [
      "project-import-jps-kotlin-10_000-modules/measureStartup",
      "project-import-jps-kotlin-50_000-modules/measureStartup",
      "project-import-jps-java-50_000-modules/measureStartup",
      "project-import-idea-community-jps/measureStartup",
    ],
  }
})
const charts = combineCharts(chartsDeclaration)


const serverConfigurator = new ServerConfigurator(dbName, dbTable)

const persistenceForDashboard = new PersistentStateManager("idea_jps_dashboard", {
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
  [timeRangeConfigurator],
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

function onChangeRange(value: TimeRange) {
  timeRangeConfigurator.value.value = value
}

const warnings = ref<Array<Accident>>()
refToObservable(timeRangeConfigurator.value).subscribe(data => {
  getAccidentsFromMetaDb(warnings, null, data)
})
</script>

<style>
.customToolbar {
  background-color: transparent;
  border: none;
  padding: 0;
}
</style>