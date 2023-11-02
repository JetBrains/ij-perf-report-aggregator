<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="idea_jps_dashboard"
    initial-machine="linux-blade-hetzner"
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
  </DashboardPage>
</template>

<script setup lang="ts">
import { ChartDefinition, combineCharts } from "../../../charts/DashboardCharts"
import GroupProjectsChart from "../../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../../common/DashboardPage.vue"

const metricsDeclaration = [
  "jps.aggregate.sync.duration",
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
  "workspaceModel.load.cache.metadata.from.file.ms",
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

  "AWTEventQueue.dispatchTimeTotal",
  "CPU | Load |Total % 95th pctl",
  "Memory | IDE | RESIDENT SIZE (MB) 95th pctl",
  "Memory | IDE | VIRTUAL SIZE (MB) 95th pctl",
  "gcPause",
  "gcPauseCount",
  "fullGCPause",
  "freedMemoryByGC",
  "totalHeapUsedMax",
  "JVM.GC.collectionTimesMs",
  "JVM.GC.collections",
  "JVM.maxHeapMegabytes",
  "JVM.threadCount",
  "JVM.totalCpuTimeMs",
  "JVM.totalMegabytesAllocated",
  "JVM.usedHeapMegabytes",
  "JVM.usedNativeMegabytes",
]

const chartsDeclaration: ChartDefinition[] = metricsDeclaration.map((metric) => {
  return {
    labels: [metric],
    measures: [metric],
    projects: [
      // JPS projects
      "project-import-jps-kotlin-10_000-modules/measureStartup",
      "project-import-jps-kotlin-50_000-modules/measureStartup",
      "project-reimport-jps-kotlin-10_000-modules/measureStartup",
      "project-reimport-jps-kotlin-50_000-modules/measureStartup",
      "project-import-from-cache-jps-kotlin-50_000-modules/measureStartup",
      "project-import-from-cache-jps-kotlin-10_000-modules/measureStartup",
      "project-import-jps-java-1_000-modules/measureStartup",
      "project-reimport-jps-java-1_000-modules/measureStartup",
      "project-import-from-cache-jps-java-1_000-modules/measureStartup",
      "project-import-idea-community-jps/measureStartup",
      "jps_10K-modules-checkout-branch-with-changes/measureStartup",
      "jps-1K-modules-checkout-branch-with-many-dependencies/measureStartup",
      "jps-cyclic-branches-checkout/measureStartup",
    ],
  }
})
const charts = combineCharts(chartsDeclaration)
</script>
