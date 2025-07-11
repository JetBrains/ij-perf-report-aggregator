<template>
  <DashboardPage
    db-name="perfintDev"
    table="ijent"
    persistent-id="ijent_benchmark_dashboard"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
    :charts="charts"
  >
    <template #configurator>
      <MeasureSelect
        :configurator="testConfigurator"
        title="Test"
      >
        <template #icon>
          <ChartBarIcon class="w-4 h-4" />
        </template>
      </MeasureSelect>
    </template>
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
import { computed } from "vue"
import { SimpleMeasureConfigurator } from "../../configurators/SimpleMeasureConfigurator"
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import MeasureSelect from "../charts/MeasureSelect.vue"
import DashboardPage from "../common/DashboardPage.vue"

const metricsDeclaration = [
  "ijent.events.count",
  "attempt.sum.ms",
  "ijent.directoryStreamClose.sum.ms",
  "ijent.directoryStreamIteratorNext.sum.ms",
  "ijent.fileSystemClose.sum.ms",
  "ijent.providerCheckAccess.sum.ms",
  "ijent.providerCopy.sum.ms",
  "ijent.providerCreateDirectory.sum.ms",
  "ijent.providerDelete.sum.ms",
  "ijent.providerMove.sum.ms",
  "ijent.providerNewByteChannel.sum.ms",
  "ijent.providerNewDirectoryStream.sum.ms",
  "ijent.providerReadAttributes.sum.ms",
  "ijent.seekableByteChannelClose.sum.ms",
  "ijent.seekableByteChannelNewPosition.sum.ms",
  "ijent.seekableByteChannelRead.sum.ms",
  "ijent.seekableByteChannelSize.sum.ms",
  "ijent.seekableByteChannelWrite.sum.ms",
]

const projects = [
  "Docker ubuntu:22.04 - copy file concurrent-1",
  "Docker ubuntu:22.04 - copy file concurrent-12",
  "Docker ubuntu:22.04 - copy file to output stream concurrent-1",
  "Docker ubuntu:22.04 - copy file to output stream concurrent-12",
  "Docker ubuntu:22.04 - copy input stream to file concurrent-1",
  "Docker ubuntu:22.04 - copy input stream to file concurrent-12",
  "Docker ubuntu:22.04 - create directory",
  "Docker ubuntu:22.04 - delete path",
  "Docker ubuntu:22.04 - Files Exists concurrent-1",
  "Docker ubuntu:22.04 - Files Exists concurrent-12",
  "Docker ubuntu:22.04 - move directories with files concurrent-1",
  "Docker ubuntu:22.04 - move directories with files concurrent-12",
  "Docker ubuntu:22.04 - move files concurrent-1",
  "Docker ubuntu:22.04 - move files concurrent-12",
  "Docker ubuntu:22.04 - newByteChannel concurrent-1",
  "Docker ubuntu:22.04 - newByteChannel concurrent-12",
  "Docker ubuntu:22.04 - newFileChannel concurrent-1",
  "Docker ubuntu:22.04 - newFileChannel concurrent-12",
  "Docker ubuntu:22.04 - readAttributes concurrent-1",
  "Docker ubuntu:22.04 - readAttributes concurrent-12",
  "Docker ubuntu:22.04 - Read file concurrent-1",
  "Docker ubuntu:22.04 - Read file concurrent-12",
  "Docker ubuntu:22.04 - Read file single request concurrent-1",
  "Docker ubuntu:22.04 - Read file single request concurrent-12",
  "Docker ubuntu:22.04 - upload file tree",

  "WSL Ubuntu (WSL2) - copy file concurrent-1",
  "WSL Ubuntu (WSL2) - copy file concurrent-24",
  "WSL Ubuntu (WSL2) - copy file to output stream concurrent-1",
  "WSL Ubuntu (WSL2) - copy file to output stream concurrent-24",
  "WSL Ubuntu (WSL2) - copy input stream to file concurrent-1",
  "WSL Ubuntu (WSL2) - copy input stream to file concurrent-24",
  "WSL Ubuntu (WSL2) - create directory",
  "WSL Ubuntu (WSL2) - delete path",
  "WSL Ubuntu (WSL2) - Files Exists concurrent-1",
  "WSL Ubuntu (WSL2) - Files Exists concurrent-24",
  "WSL Ubuntu (WSL2) - move directories with files concurrent-1",
  "WSL Ubuntu (WSL2) - move directories with files concurrent-24",
  "WSL Ubuntu (WSL2) - move files concurrent-1",
  "WSL Ubuntu (WSL2) - move files concurrent-24",
  "WSL Ubuntu (WSL2) - newByteChannel concurrent-1",
  "WSL Ubuntu (WSL2) - newByteChannel concurrent-24",
  "WSL Ubuntu (WSL2) - newFileChannel concurrent-1",
  "WSL Ubuntu (WSL2) - newFileChannel concurrent-24",
  "WSL Ubuntu (WSL2) - readAttributes concurrent-1",
  "WSL Ubuntu (WSL2) - readAttributes concurrent-24",
  "WSL Ubuntu (WSL2) - Read file concurrent-1",
  "WSL Ubuntu (WSL2) - Read file concurrent-24",
  "WSL Ubuntu (WSL2) - Read file single request concurrent-1",
  "WSL Ubuntu (WSL2) - Read file single request concurrent-24",
  "WSL Ubuntu (WSL2) - upload file tree",

  "Standard WSL - Default FS - copy file concurrent-1",
  "Standard WSL - Default FS - copy file concurrent-24",
  "Standard WSL - Default FS - copy file to output stream concurrent-1",
  "Standard WSL - Default FS - copy file to output stream concurrent-24",
  "Standard WSL - Default FS - copy input stream to file concurrent-1",
  "Standard WSL - Default FS - copy input stream to file concurrent-24",
  "Standard WSL - Default FS - create directory",
  "Standard WSL - Default FS - delete path",
  "Standard WSL - Default FS - Files Exists concurrent-1",
  "Standard WSL - Default FS - Files Exists concurrent-24",
  "Standard WSL - Default FS - move directories with files concurrent-1",
  "Standard WSL - Default FS - move directories with files concurrent-24",
  "Standard WSL - Default FS - move files concurrent-1",
  "Standard WSL - Default FS - move files concurrent-24",
  "Standard WSL - Default FS - newByteChannel concurrent-1",
  "Standard WSL - Default FS - newByteChannel concurrent-24",
  "Standard WSL - Default FS - newFileChannel concurrent-1",
  "Standard WSL - Default FS - newFileChannel concurrent-24",
  "Standard WSL - Default FS - readAttributes concurrent-1",
  "Standard WSL - Default FS - readAttributes concurrent-24",
  "Standard WSL - Default FS - Read file concurrent-1",
  "Standard WSL - Default FS - Read file concurrent-24",
  "Standard WSL - Default FS - upload file tree",

  "Local Linux X86_64 - copy file concurrent-1",
  "Local Linux X86_64 - copy file concurrent-12",
  "Local Linux X86_64 - copy file to output stream concurrent-1",
  "Local Linux X86_64 - copy file to output stream concurrent-12",
  "Local Linux X86_64 - copy input stream to file concurrent-1",
  "Local Linux X86_64 - copy input stream to file concurrent-12",
  "Local Linux X86_64 - create directory",
  "Local Linux X86_64 - delete path",
  "Local Linux X86_64 - Files Exists concurrent-1",
  "Local Linux X86_64 - Files Exists concurrent-12",
  "Local Linux X86_64 - move directories with files concurrent-1",
  "Local Linux X86_64 - move directories with files concurrent-12",
  "Local Linux X86_64 - move files concurrent-1",
  "Local Linux X86_64 - move files concurrent-12",
  "Local Linux X86_64 - newByteChannel concurrent-1",
  "Local Linux X86_64 - newByteChannel concurrent-12",
  "Local Linux X86_64 - newFileChannel concurrent-1",
  "Local Linux X86_64 - newFileChannel concurrent-12",
  "Local Linux X86_64 - readAttributes concurrent-1",
  "Local Linux X86_64 - readAttributes concurrent-12",
  "Local Linux X86_64 - Read file concurrent-1",
  "Local Linux X86_64 - Read file concurrent-12",
  "Local Linux X86_64 - Read file single request concurrent-1",
  "Local Linux X86_64 - Read file single request concurrent-12",
  "Local Linux X86_64 - upload file tree",

  "Local Eel (Posix) - Default FS - copy file concurrent-1",
  "Local Eel (Posix) - Default FS - copy file concurrent-12",
  "Local Eel (Posix) - Default FS - copy file to output stream concurrent-1",
  "Local Eel (Posix) - Default FS - copy file to output stream concurrent-12",
  "Local Eel (Posix) - Default FS - copy input stream to file concurrent-1",
  "Local Eel (Posix) - Default FS - copy input stream to file concurrent-12",
  "Local Eel (Posix) - Default FS - create directory",
  "Local Eel (Posix) - Default FS - delete path",
  "Local Eel (Posix) - Default FS - Files Exists concurrent-1",
  "Local Eel (Posix) - Default FS - Files Exists concurrent-12",
  "Local Eel (Posix) - Default FS - move directories with files concurrent-1",
  "Local Eel (Posix) - Default FS - move directories with files concurrent-12",
  "Local Eel (Posix) - Default FS - move files concurrent-1",
  "Local Eel (Posix) - Default FS - move files concurrent-12",
  "Local Eel (Posix) - Default FS - newByteChannel concurrent-1",
  "Local Eel (Posix) - Default FS - newByteChannel concurrent-12",
  "Local Eel (Posix) - Default FS - newFileChannel concurrent-1",
  "Local Eel (Posix) - Default FS - newFileChannel concurrent-12",
  "Local Eel (Posix) - Default FS - readAttributes concurrent-1",
  "Local Eel (Posix) - Default FS - readAttributes concurrent-12",
  "Local Eel (Posix) - Default FS - Read file concurrent-1",
  "Local Eel (Posix) - Default FS - Read file concurrent-12",
  "Local Eel (Posix) - Default FS - upload file tree",
  "Local Eel (Posix) - MultiRoutingFileSystem - copy file concurrent-1",
  "Local Eel (Posix) - MultiRoutingFileSystem - copy file concurrent-12",
  "Local Eel (Posix) - MultiRoutingFileSystem - copy file to output stream concurrent-1",
  "Local Eel (Posix) - MultiRoutingFileSystem - copy file to output stream concurrent-12",
  "Local Eel (Posix) - MultiRoutingFileSystem - copy input stream to file concurrent-1",
  "Local Eel (Posix) - MultiRoutingFileSystem - copy input stream to file concurrent-12",
  "Local Eel (Posix) - MultiRoutingFileSystem - create directory",
  "Local Eel (Posix) - MultiRoutingFileSystem - delete path",
  "Local Eel (Posix) - MultiRoutingFileSystem - Files Exists concurrent-1",
  "Local Eel (Posix) - MultiRoutingFileSystem - Files Exists concurrent-12",
  "Local Eel (Posix) - MultiRoutingFileSystem - move directories with files concurrent-1",
  "Local Eel (Posix) - MultiRoutingFileSystem - move directories with files concurrent-12",
  "Local Eel (Posix) - MultiRoutingFileSystem - move files concurrent-1",
  "Local Eel (Posix) - MultiRoutingFileSystem - move files concurrent-12",
  "Local Eel (Posix) - MultiRoutingFileSystem - newByteChannel concurrent-1",
  "Local Eel (Posix) - MultiRoutingFileSystem - newByteChannel concurrent-12",
  "Local Eel (Posix) - MultiRoutingFileSystem - newFileChannel concurrent-1",
  "Local Eel (Posix) - MultiRoutingFileSystem - newFileChannel concurrent-12",
  "Local Eel (Posix) - MultiRoutingFileSystem - readAttributes concurrent-1",
  "Local Eel (Posix) - MultiRoutingFileSystem - readAttributes concurrent-12",
  "Local Eel (Posix) - MultiRoutingFileSystem - Read file concurrent-1",
  "Local Eel (Posix) - MultiRoutingFileSystem - Read file concurrent-12",
  "Local Eel (Posix) - MultiRoutingFileSystem - upload file tree",

  "Local Eel (Windows) - Default FS - copy file concurrent-1",
  "Local Eel (Windows) - Default FS - copy file concurrent-24",
  "Local Eel (Windows) - Default FS - copy file to output stream concurrent-1",
  "Local Eel (Windows) - Default FS - copy file to output stream concurrent-24",
  "Local Eel (Windows) - Default FS - copy input stream to file concurrent-1",
  "Local Eel (Windows) - Default FS - copy input stream to file concurrent-24",
  "Local Eel (Windows) - Default FS - create directory",
  "Local Eel (Windows) - Default FS - delete path",
  "Local Eel (Windows) - Default FS - Files Exists concurrent-1",
  "Local Eel (Windows) - Default FS - Files Exists concurrent-24",
  "Local Eel (Windows) - Default FS - move directories with files concurrent-1",
  "Local Eel (Windows) - Default FS - move directories with files concurrent-24",
  "Local Eel (Windows) - Default FS - move files concurrent-1",
  "Local Eel (Windows) - Default FS - move files concurrent-24",
  "Local Eel (Windows) - Default FS - newByteChannel concurrent-1",
  "Local Eel (Windows) - Default FS - newByteChannel concurrent-24",
  "Local Eel (Windows) - Default FS - newFileChannel concurrent-1",
  "Local Eel (Windows) - Default FS - newFileChannel concurrent-24",
  "Local Eel (Windows) - Default FS - readAttributes concurrent-1",
  "Local Eel (Windows) - Default FS - readAttributes concurrent-24",
  "Local Eel (Windows) - Default FS - Read file concurrent-1",
  "Local Eel (Windows) - Default FS - Read file concurrent-24",
  "Local Eel (Windows) - Default FS - upload file tree",

  "Local Eel (Windows) - MultiRoutingFileSystem - copy file concurrent-1",
  "Local Eel (Windows) - MultiRoutingFileSystem - copy file concurrent-24",
  "Local Eel (Windows) - MultiRoutingFileSystem - copy file to output stream concurrent-1",
  "Local Eel (Windows) - MultiRoutingFileSystem - copy file to output stream concurrent-24",
  "Local Eel (Windows) - MultiRoutingFileSystem - copy input stream to file concurrent-1",
  "Local Eel (Windows) - MultiRoutingFileSystem - copy input stream to file concurrent-24",
  "Local Eel (Windows) - MultiRoutingFileSystem - create directory",
  "Local Eel (Windows) - MultiRoutingFileSystem - delete path",
  "Local Eel (Windows) - MultiRoutingFileSystem - Files Exists concurrent-1",
  "Local Eel (Windows) - MultiRoutingFileSystem - Files Exists concurrent-24",
  "Local Eel (Windows) - MultiRoutingFileSystem - move directories with files concurrent-1",
  "Local Eel (Windows) - MultiRoutingFileSystem - move directories with files concurrent-24",
  "Local Eel (Windows) - MultiRoutingFileSystem - move files concurrent-1",
  "Local Eel (Windows) - MultiRoutingFileSystem - move files concurrent-24",
  "Local Eel (Windows) - MultiRoutingFileSystem - newByteChannel concurrent-1",
  "Local Eel (Windows) - MultiRoutingFileSystem - newByteChannel concurrent-24",
  "Local Eel (Windows) - MultiRoutingFileSystem - newFileChannel concurrent-1",
  "Local Eel (Windows) - MultiRoutingFileSystem - newFileChannel concurrent-24",
  "Local Eel (Windows) - MultiRoutingFileSystem - readAttributes concurrent-1",
  "Local Eel (Windows) - MultiRoutingFileSystem - readAttributes concurrent-24",
  "Local Eel (Windows) - MultiRoutingFileSystem - Read file concurrent-1",
  "Local Eel (Windows) - MultiRoutingFileSystem - Read file concurrent-24",
  "Local Eel (Windows) - MultiRoutingFileSystem - upload file tree",
]

const testConfigurator = new SimpleMeasureConfigurator("project", null)
testConfigurator.initData(projects)

const charts = computed(() => {
  const chartsDeclaration: ChartDefinition[] = metricsDeclaration.map((metric) => {
    return {
      labels: [metric],
      measures: [metric],
      projects: testConfigurator.selected.value ?? [],
    }
  })
  return combineCharts(chartsDeclaration)
})
</script>
