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
        :projects="['kubernetes/localInspection', 'GO-5422/localInspection']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="GloBbal inspection execution time: light projects"
        measure="globalInspections"
        :projects="['delve/inspection', 'flux/inspection', 'istio/inspection']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
