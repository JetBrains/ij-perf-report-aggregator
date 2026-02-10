<template>
  <DashboardPage
    db-name="perfUnitTests"
    table="report"
    persistent-id="lspDashboard"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  >
    <section>
      <GroupProjectsChart
        label="Diagnostic Performance"
        measure="attempt.mean.ms"
        :projects="lspProjects"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../common/DashboardPage.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"

const lspProjects = ["com.jetbrains.ls.lsp.performanceTest.DiagnosticPerformanceTest$huge$1.huge", "com.jetbrains.ls.lsp.performanceTest.DiagnosticPerformanceTest$tiny$1.tiny"]
</script>
