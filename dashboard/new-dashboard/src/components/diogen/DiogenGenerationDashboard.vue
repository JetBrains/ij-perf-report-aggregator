<template>
  <DashboardPage
    db-name="diogen"
    table="report"
    persistent-id="diogen_generation_dashboard"
    :with-installer="false"
    :initial-machine="null"
  >
    <section>
      <GroupProjectsChart
        v-for="project in PROJECTS"
        :key="project"
        :label="project"
        :measure="MEASURES"
        :projects="[project]"
        :value-unit="'counter'"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../common/DashboardPage.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"

const PROJECTS = ["ijAllCur", "freezeFlat"]
const MEASURES = ["hosts-last-created-issue", "max-hosts-unsorted"]
</script>
