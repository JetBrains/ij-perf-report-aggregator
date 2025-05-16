<template>
  <DashboardPage
    db-name="perfintDev"
    table="ijent"
    persistent-id="ijent_benchmark_dashboard"
    initial-machine="windows-azure"
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
  "ijent.fileSystemClose.sum.ms",
  "ijent.providerCheckAccess.sum.ms",
  "ijent.providerCopy.sum.ms",
  "ijent.providerCreateDirectory.sum.ms",
  "ijent.providerDelete.sum.ms",
  "ijent.providerMove.sum.ms",
  "ijent.providerNewByteChannel.sum.ms",
  "ijent.providerReadAttributes.sum.ms",
  "ijent.seekableByteChannelClose.sum.ms",
  "ijent.seekableByteChannelNewPosition.sum.ms",
  "ijent.seekableByteChannelRead.sum.ms",
  "ijent.seekableByteChannelSize.sum.ms",
  "ijent.seekableByteChannelWrite.sum.ms",
]

const projects = [
  "IJent - Files Exists",
  "IJent - move files",
  "IJent - move directories with files",
  "IJent - delete path",
  "IJent - create directory",
  "IJent - upload file tree",
  "IJent - readAttributes",
  "IJent - readAttributes parallel",
  "IJent - newFileChannel",
  "IJent - newFileChannel parallel",
  "IJent - newByteChannel",
  "IJent - newByteChannel parallel",
  "IJent - copy input stream to file",
  "IJent - copy input stream to file parallel",
  "IJent - copy file to output stream",
  "IJent - copy file to output stream parallel",
  "IJent - copy file",
  "IJent - copy file parallel",
  "IJent - Read file quick",
  "IJent - Read file quick parallel",
  "IJent - Read file",
  "IJent - Read file parallel",
  "IJent - Read Java ZipFile",
  "IJent - Read Java ZipFile parallel",
  "IJent - Read JBZipFile",
  "IJent - Read JBZipFile file parallel",
  "WSL - Files Exists",
  "WSL - move files",
  "WSL - move directories with files",
  "WSL - delete path",
  "WSL - create directory",
  "WSL - upload file tree",
  "WSL - readAttributes",
  "WSL - readAttributes parallel",
  "WSL - newFileChannel",
  "WSL - newFileChannel parallel",
  "WSL - newByteChannel",
  "WSL - newByteChannel parallel",
  "WSL - copy input stream to file",
  "WSL - copy input stream to file parallel",
  "WSL - copy file to output stream",
  "WSL - copy file to output stream parallel",
  "WSL - copy file",
  "WSL - copy file parallel",
  "WSL - Read file quick",
  "WSL - Read file quick parallel",
  "WSL - Read file",
  "WSL - Read file parallel",
  "WSL - Read Java ZipFile",
  "WSL - Read Java ZipFile parallel",
  "WSL - Read JBZipFile",
  "WSL - Read JBZipFile file parallel",
  "Windows - Read file",
  "Windows - Read file parallel",
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
