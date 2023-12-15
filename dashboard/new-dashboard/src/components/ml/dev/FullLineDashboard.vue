<template>
  <DashboardPage
    db-name="perfintDev"
    table="ml"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    persistent-id="fullLineDashboardDev"
    :with-installer="false"
  >
    <section>
      <GroupProjectsChart
        label="Full line completion "
        :measure="['inlineCompletionShow#mean_value']"
        :projects="['kotlin_language_server/full-line-completion/DefinitionTest/empty-fun-full-line']"
      />
    </section>
  </DashboardPage>
</template>
<script setup lang="ts">
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"
</script>
