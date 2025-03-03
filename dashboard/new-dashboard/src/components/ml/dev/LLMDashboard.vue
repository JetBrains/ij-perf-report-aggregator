<template>
  <DashboardPage
    db-name="perfintDev"
    table="ml"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    persistent-id="llmDashboardDev"
    :with-installer="false"
  >
    <section>
      <GroupProjectsChart
        label="Inline completion"
        measure="callInlineCompletionOnCompletion#mean_value"
        :projects="['gradle-calculator_SimpleInlineCompletionTest/simple inline completion']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Code generation"
        measure="ai-generate-code#mean_value"
        :projects="['gradle-calculator_CodeGenerationPerformanceTest/generate code']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Do generate"
        measure="doGenerate#mean_value"
        :projects="['gradle-calculator_CodeGenerationPerformanceTest/generate code']"
      ></GroupProjectsChart>
    </section>
  </DashboardPage>
</template>
<script setup lang="ts">
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"
</script>
