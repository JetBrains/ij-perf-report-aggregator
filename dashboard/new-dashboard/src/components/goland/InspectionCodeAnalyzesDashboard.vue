<template>
  <DashboardPage
    db-name="perfint"
    table="goland"
    persistent-id="goland_code_analyzes_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <section>
      <GroupProjectsChart
        label="Local inspection execution time"
        measure="localInspections"
        :projects="['goland/istio/localInspection', 'goland/minotaur/localInspection']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Global inspection execution time: light projects"
        measure="globalInspections"
        :projects="['delve/inspection', 'kubernetes/inspection']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
