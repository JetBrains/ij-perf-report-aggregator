<template>
  <DashboardPage
    db-name="perfint"
    table="goland"
    persistent-id="goland_completion_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <section>
      <GroupProjectsChart
        label="Typing Latency"
        measure="typing#latency#mean_value"
        :projects="['typing/typing/mattermost']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Basic Completion"
        measure="completion"
        :projects="['caddy/completion/variable', 'caddy/completion/type', 'caddy/completion/method', 'caddy/completion/interface']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Smart Completion"
        measure="completion"
        :projects="['permify/completion/method', 'permify/completion/any']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
