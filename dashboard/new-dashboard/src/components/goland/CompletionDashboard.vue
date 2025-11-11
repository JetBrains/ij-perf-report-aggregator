<template>
  <DashboardPage
    db-name="perfintDev"
    table="goland"
    persistent-id="goland_completion_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :with-installer="false"
  >
    <section>
      <GroupProjectsChart
        label="Typing: Total Time"
        measure="typing"
        :projects="['act/typing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Completion"
        measure="completion#mean_value"
        :projects="['caddy/completion/variable', 'caddy/completion/type', 'caddy/completion/return', 'caddy/completion/interface', 'caddy/completion/import']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Completion First Element"
        measure="completion#firstElementShown#mean_value"
        :projects="['caddy/completion/variable', 'caddy/completion/type', 'caddy/completion/return', 'caddy/completion/interface', 'caddy/completion/import']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Type Completion"
        measure="completion#mean_value"
        :projects="['permify/completion/method', 'permify/completion/any']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Type Completion First Element"
        measure="completion#firstElementShown#mean_value"
        :projects="['permify/completion/method', 'permify/completion/any']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
