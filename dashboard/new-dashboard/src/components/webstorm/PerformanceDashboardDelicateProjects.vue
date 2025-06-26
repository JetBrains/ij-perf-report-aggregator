<template>
  <DashboardPage
    db-name="perfintDev"
    :with-installer="false"
    table="webstorm"
    persistent-id="webstorm_dashboard_delicate_projects"
    initial-machine="linux-blade-hetzner"
  >
    <template
      v-for="group in groups"
      :key="group.measure"
    >
      <Divider :title="group.label" />
      <section
        v-for="(groupOf3, groupOf3index) in groupBy3(group.projects)"
        :key="groupOf3index"
        class="flex gap-x-6 flex-col md:flex-row"
      >
        <div
          v-for="project in groupOf3"
          :key="project"
          class="flex-1 min-w-0"
        >
          <GroupProjectsChart
            :label="project"
            :measure="group.measure"
            :projects="[project, project + 'NEXT']"
          />
        </div>
      </section>
    </template>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
import { groupBy3 } from "./utils"

const groups = [
  {
    label: "FirstCodeAnalysis",
    measure: "firstCodeAnalysis",
    projects: ["clickUp_frontend/localInspection/app.component.html", "clickUp_frontend/localInspection/app.component.ts"],
  },
]
</script>
