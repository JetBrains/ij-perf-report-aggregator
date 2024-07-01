<template>
  <DashboardPage
    db-name="perfint"
    table="goland"
    persistent-id="goland_completion_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <section>
      <GroupProjectsChart
        label="Typing Latency (mattermost)"
        measure="typing#latency#mean_value"
        :projects="['goland/typing/mattermost']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Auto-import Completion"
        measure="completion"
        :projects="['goland/autoImportCompletion/autoImportCompletion']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Basic Completion"
        measure="completion"
        :projects="['goland/basicCompletionSpeed/basicCompletionSpeed']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Smart Completion"
        measure="completion"
        :projects="['goland/smartCompletionSpeed/smartCompletionSpeed']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
