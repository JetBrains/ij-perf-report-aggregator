<template>
  <DashboardPage
    db-name="perfintDev"
    table="clion"
    persistent-id="clion_goto_declaration_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :with-installer="false"
  >
    <section
      v-for="chart in charts"
      :key="chart.label"
      class="flex gap-x-6 flex-col md:flex-row"
    >
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          :label="`Go to: ${chart.label}`"
          :measure="measure"
          :projects="chart.projects"
        />
      </div>
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

interface ChartDef {
  label: string
  projects: string[]
}

const measure = ["clionGotoDeclaration"]

const charts: ChartDef[] = [
  { label: "Constructor declaration", projects: ["radler/luau/gotoDeclaration/AstStatDeclareFunction.ctor"] },
  { label: "Method declaration", projects: ["radler/luau/gotoDeclaration/TypeChecker.getScopes"] },
  { label: "std::string declaration", projects: ["radler/luau/gotoDeclaration/std.string"] },
  { label: "Macro declaration", projects: ["radler/luau/gotoDeclaration/LUAU_ASSERT"] },
  { label: "Include header", projects: ["radler/luau/gotoDeclaration/time.h"] },
]
</script>
