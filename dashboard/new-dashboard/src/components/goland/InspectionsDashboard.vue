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
        :projects="['goland/delve/inspection', 'goland/kubernetes/inspection']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Single inspections execution time (cockroach)"
        measure="globalInspections"
        :projects="[
          'goland/singleInspection/cockroach/VgoDependencyUpdateAvailable',
          'goland/singleInspection/cockroach/GoUnnecessarilyExportedIdentifiers',
          'goland/singleInspection/cockroach/GoNilness',
          'goland/singleInspection/cockroach/GoBoolExpressions',
        ]"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
