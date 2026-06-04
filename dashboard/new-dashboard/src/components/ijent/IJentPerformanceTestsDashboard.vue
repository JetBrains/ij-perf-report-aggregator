<template>
  <DashboardPage
    db-name="perfintDev"
    table="ijent"
    persistent-id="ijent_performance_dashboard"
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

// Pre-rename raw project names from before the 2026-05-25 unification (IJPL-245044). Each
// canonical name maps to zero or more aliases. Each alias is rendered as its own line series in
// echarts, sharing the legend label with its canonical sibling; sparse pre-rename series (1-2
// points) will draw long diagonal strokes. Opt in via the "Pre-rename data" toolbar toggle.
const PRE_RENAME_ALIASES: Record<string, string[]> = {
  "intellij-community/indexing/Local": ["ijent-import-intellij-Local"],
  "intellij-community/indexing/Docker": ["ijent-import-intellij-Docker"],
  "intellij-community/indexing/WSL": [],
  "intellij-community/build/Local": [],
  "intellij-community/build/Docker": ["ijent-build-intellij-Docker"],
  "intellij-community/build/WSL": ["ijent-build-intellij-WSL"],
  "php-project/indexing/Local": ["php-project/indexing/Local/indexingLocal", "indexing-php-project/indexingLocal"],
  "php-project/indexing/Docker": ["php-project/indexing/Docker/indexingDocker", "indexing-php-project/indexingDocker"],
  "php-project/indexing/WSL": ["php-project/indexing/WSL/indexingWSL", "indexing-php-project/indexingWSL"],
  "spring-pet-clinic-maven/indexing/Local": ["spring-pet-clinic-maven/indexing/Local/indexingLocal", "spring-pet-clinic-maven/indexingLocal"],
  "spring-pet-clinic-maven/indexing/Docker": ["spring-pet-clinic-maven/indexing/Docker/indexingDocker", "spring-pet-clinic-maven/indexingDocker"],
  "spring-pet-clinic-maven/indexing/WSL": ["spring-pet-clinic-maven/indexing/WSL/indexingWSL", "spring-pet-clinic-maven/indexingWSL"],
  "spring-pet-clinic-gradle/indexing/Local": ["spring-pet-clinic-gradle/indexing/Local/indexingLocal", "spring-pet-clinic-gradle/indexingLocal"],
  "spring-pet-clinic-gradle/indexing/Docker": ["spring-pet-clinic-gradle/indexing/Docker/indexingDocker", "spring-pet-clinic-gradle/indexingDocker"],
  "spring-pet-clinic-gradle/indexing/WSL": ["spring-pet-clinic-gradle/indexing/WSL/indexingWSL", "spring-pet-clinic-gradle/indexingWSL"],
  "jps-1000-modules/import/Local": ["ijent-import-jps-1000-modules-Local", "nio_default-import-jps-1000-modules-Local"],
  "jps-1000-modules/import/Docker": ["ijent-import-jps-1000-modules-Docker"],
  "jps-1000-modules/import/WSL": ["wsl-import-jps-1000-modules-WSL"],
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

  const communityBuild = envSeries({
    Local: withPreRenameAliases("intellij-community/build/Local", include),
    Docker: withPreRenameAliases("intellij-community/build/Docker", include),
    WSL: withPreRenameAliases("intellij-community/build/WSL", include),
  })

  const phpIndexing = envSeries({
    Local: withPreRenameAliases("php-project/indexing/Local", include),
    Docker: withPreRenameAliases("php-project/indexing/Docker", include),
    WSL: withPreRenameAliases("php-project/indexing/WSL", include),
  })

  const springMavenIndexing = envSeries({
    Local: withPreRenameAliases("spring-pet-clinic-maven/indexing/Local", include),
    Docker: withPreRenameAliases("spring-pet-clinic-maven/indexing/Docker", include),
    WSL: withPreRenameAliases("spring-pet-clinic-maven/indexing/WSL", include),
  })

  const springGradleIndexing = envSeries({
    Local: withPreRenameAliases("spring-pet-clinic-gradle/indexing/Local", include),
    Docker: withPreRenameAliases("spring-pet-clinic-gradle/indexing/Docker", include),
    WSL: withPreRenameAliases("spring-pet-clinic-gradle/indexing/WSL", include),
  })

  const jpsImport = envSeries({
    Local: withPreRenameAliases("jps-1000-modules/import/Local", include),
    Docker: withPreRenameAliases("jps-1000-modules/import/Docker", include),
    WSL: withPreRenameAliases("jps-1000-modules/import/WSL", include),
  })

  const chartsDeclaration: ChartDefinition[] = [
    {
      labels: ["First Scanning — Community", "Indexing — Community", "Dumb Mode — Community", "Indexed Files — Community"],
      measures: ["scanningTimeWithoutPauses", "indexingTimeWithoutPauses", "dumbModeTimeWithPauses", "numberOfIndexedFiles"],
      ...communityIndexing,
    },
    {
      labels: ["Indexing — PHP"],
      measures: ["indexingTimeWithoutPauses"],
      ...phpIndexing,
    },
    {
      labels: ["Indexing — Spring Pet Clinic (Maven)"],
      measures: ["indexingTimeWithoutPauses"],
      ...springMavenIndexing,
    },
    {
      labels: ["Indexing — Spring Pet Clinic (Gradle)"],
      measures: ["indexingTimeWithoutPauses"],
      ...springGradleIndexing,
    },
    {
      labels: ["Build — Community"],
      measures: ["build_compilation_duration"],
      ...communityBuild,
    },
    {
      labels: ["Project Opening — JPS 1000 Modules", "Indexing — JPS 1000 Modules"],
      measures: ["project.opening", "indexingTimeWithoutPauses"],
      ...jpsImport,
    },
  ]

  return combineCharts(chartsDeclaration)
})
</script>
