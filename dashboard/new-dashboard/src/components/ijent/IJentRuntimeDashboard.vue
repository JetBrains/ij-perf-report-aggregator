<template>
  <DashboardPage
    db-name="perfintDev"
    table="ijent"
    persistent-id="ijent_runtime_dashboard"
    initial-machine="Linux Munich i7-13700, 64 Gb"
    :with-installer="false"
    :with-mode="false"
    :charts="charts"
  >
    <template #configurator>
      <PreRenameDataSwitch v-model="preRenameData" />
    </template>
    <section>
      <GroupProjectsChart
        v-for="chart in charts"
        :key="chart.definition.label"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
        :aliases="chart.aliases"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import { computed, ref } from "vue"
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import PreRenameDataSwitch from "../settings/PreRenameDataSwitch.vue"

type Env = "Local" | "Docker" | "WSL"

function envSeries(envMap: Record<Env, string[]>): { projects: string[]; aliases: string[] } {
  const projects: string[] = []
  const aliases: string[] = []
  for (const env of ["Local", "Docker", "WSL"] as const) {
    for (const p of envMap[env]) {
      projects.push(p)
      aliases.push(env)
    }
  }
  return { projects, aliases }
}

const PRE_RENAME_ALIASES: Record<string, string[]> = {
  "intellij-community/indexing/Local": ["ijent-import-intellij-Local"],
  "intellij-community/indexing/Docker": ["ijent-import-intellij-Docker"],
  "intellij-community/indexing/WSL": [],
}

function withPreRenameAliases(canonical: string, include: boolean): string[] {
  return include ? [canonical, ...(PRE_RENAME_ALIASES[canonical] ?? [])] : [canonical]
}

const preRenameData = ref(false)

const charts = computed(() => {
  const include = preRenameData.value

  const communityIndexing = envSeries({
    Local: withPreRenameAliases("intellij-community/indexing/Local", include),
    Docker: withPreRenameAliases("intellij-community/indexing/Docker", include),
    WSL: withPreRenameAliases("intellij-community/indexing/WSL", include),
  })

  const chartsDeclaration: ChartDefinition[] = [
    {
      labels: [
        "IDE Resident Memory (95p, MB)",
        "IDE Virtual Memory (95p, MB)",
        "Avg RAM (MB)",
        "Avg File-Mappings RAM (MB)",
        "Max Heap (MB)",
        "Max Threads",
        "Total CPU Time (ms)",
        "GC Time (ms)",
        "GC Collections",
        "GC Pause Total (ms)",
        "GC Pause Count",
        "Full GC Pause (ms)",
        "Freed by GC",
        "AWT Dispatch Total (ms)",
        "IJent Events Count",
      ],
      measures: [
        "Memory | IDE | RESIDENT SIZE (MB) 95th pctl",
        "Memory | IDE | VIRTUAL SIZE (MB) 95th pctl",
        "MEM.avgRamMegabytes",
        "MEM.avgFileMappingsRamMegabytes",
        "JVM.maxHeapMegabytes",
        "JVM.maxThreadCount",
        "JVM.totalCpuTimeMs",
        "JVM.GC.collectionTimesMs",
        "JVM.GC.collections",
        "gcPause",
        "gcPauseCount",
        "fullGCPause",
        "freedMemoryByGC",
        "AWTEventQueue.dispatchTimeTotal",
        "ijent.events.count",
      ],
      ...communityIndexing,
    },
  ]

  return combineCharts(chartsDeclaration)
})
</script>
