<template>
  <DashboardPage
    db-name="perfint"
    table="goland"
    persistent-id="goland_completion_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section>
      <GroupProjectsChart
        label="Typing Code Analysis"
        measure="typingCodeAnalyzing"
        :projects="['act/typing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Basic Completion"
        measure="completion"
        :projects="['caddy/completion/variable', 'caddy/completion/type', 'caddy/completion/return', 'caddy/completion/interface', 'caddy/completion/import']"
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
