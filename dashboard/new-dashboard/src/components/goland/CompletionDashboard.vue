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
        label="Auto-import Completion (kubernetes)"
        measure="completion"
        :projects="['kubernetesImport/completion']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Basic Completion"
        measure="completion"
        :projects="['caddyCompletion/completion/variable', 'caddyCompletion/completion/type', 'caddyCompletion/completion/method', 'caddyCompletion/completion/interface']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Smart Completion"
        measure="completion"
        :projects="['smartCompletion/completion']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
