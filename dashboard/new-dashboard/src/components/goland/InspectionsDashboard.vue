<template>
  <DashboardPage
    db-name="perfintDev"
    table="goland"
    :with-installer="false"
    persistent-id="goland_code_analyzes_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section>
      <GroupProjectsChart
        label="Local inspection execution time"
        measure="localInspections"
        :projects="['istio/localInspection/adsc.go', 'minotaur/localInspection/server.go']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Global inspection execution time"
        measure="globalInspections"
        :projects="['delve/inspection', 'kubernetes/inspection']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Single inspections execution time (cockroach)"
        measure="globalInspections"
        :projects="['singleInspection/cockroach/VgoDependencyUpdateAvailable']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
