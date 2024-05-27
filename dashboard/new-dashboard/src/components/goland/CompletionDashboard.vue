<template>
  <DashboardPage
    db-name="perfint"
    table="goland"
    persistent-id="goland_completion_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <section>
      <GroupProjectsChart
        label="Auto-import completion"
        measure="completion"
        :projects="['autoImportCompletion']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Smart completion"
        measure="completion"
        :projects="['smartCompletionSpeed']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
