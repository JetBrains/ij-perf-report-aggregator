<template>
  <DashboardPage
    db-name="perfint"
    table="goland"
    persistent-id="goland_completion_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <section class="flex gap-x-6 flex-col md:flex-row">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Auto-import Completion"
          measure="completion"
          :projects="['goland/autoImportCompletion/autoImportCompletion']"
        />
      </div>

      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Basic Completion"
          measure="completion"
          :projects="['goland/basicCompletionSpeed/basicCompletionSpeed']"
        />
      </div>

      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Smart Completion"
          measure="completion"
          :projects="['goland/smartCompletionSpeed/smartCompletionSpeed']"
        />
      </div>
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
