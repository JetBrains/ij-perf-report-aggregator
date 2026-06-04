<template>
  <DashboardPage
    db-name="perfintDev"
    table="ijent"
    persistent-id="ijent_project_loading_dashboard"
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
        "Project Opening",
        "JPS — Aggregate Sync Duration",
        "JPS — Aggregate Counters",
        "JPS — Project Serializers Load",
        "JPS — Config Reader Load",
        "JPS — Content Reader Load",
        "Workspace Model — Total Loading",
        "Workspace Model — Module Bridge Loader",
        "Workspace Model — Module Manager Bridge",
      ],
      measures: [
        "project.opening",
        "jps.aggregate.sync.duration",
        "jps.aggregate.counters",
        "jps.project.serializers.load.ms",
        "jps.storage.jps.conf.reader.load.component.ms",
        "jps.app.storage.content.reader.load.component.ms",
        "workspaceModel.loading.total.ms",
        "workspaceModel.moduleBridgeLoader.loading.modules.ms",
        "workspaceModel.moduleManagerBridge.load.module.ms",
      ],
      ...communityIndexing,
    },
  ]

  return combineCharts(chartsDeclaration)
})
</script>
