<template>
  <DashboardPage
    db-name="perfintDev"
    table="ml"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    persistent-id="llmPhpStormDashboardDev"
    :with-installer="false"
  >
    <section>
      <GroupProjectsChart
        label="LLM completion "
        :measure="['inlineCompletionShow#mean_value']"
        :projects="['Symfony/completion/echoStatement', 'Symfony/completion/simpleResolve']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
