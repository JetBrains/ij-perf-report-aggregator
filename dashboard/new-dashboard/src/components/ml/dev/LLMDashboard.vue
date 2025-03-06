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
        :projects="['SimpleInlineCompletionTest_gradle-calculator/simple inline completion']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Code generation"
        measure="ai-generate-code#mean_value"
        :projects="['CodeGenerationPerformanceTest_gradle-calculator/generate code']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Do generate"
        measure="doGenerate#mean_value"
        :projects="['CodeGenerationPerformanceTest_gradle-calculator/generate code']"
      ></GroupProjectsChart>
    </section>
    <section>
      <GroupProjectsChart
        label="AI Ignore Sum"
        :measure="['ai-ignore.sum.ms', 'ai-ignore-accumulative-x100.sum.ms']"
        :projects="[
          'kotlinx_coroutines_k2_dev_AiIgnoreCoroutinesTest/test aiignore performance',
          'kotlinx_coroutines_k2_dev_AiIgnoreCoroutinesAccumulativeTest/test aiignore performance with repetitions on each file',
        ]"
      ></GroupProjectsChart>
    </section>
    <section>
      <GroupProjectsChart
        label="AI Ignore Single"
        :measure="['ai-ignore.mean.ns', 'ai-ignore.median.ns', 'ai-ignore.95.percentile.ns']"
        :projects="['kotlinx_coroutines_k2_dev_AiIgnoreCoroutinesTest/test aiignore performance']"
      ></GroupProjectsChart>
    </section>
    <section>
      <GroupProjectsChart
        label="AI Ignore Multiple X100"
        :measure="['ai-ignore-accumulative-x100.mean.ns', 'ai-ignore-accumulative-x100.median.ns', 'ai-ignore-accumulative-x100.95.percentile.ns']"
        :projects="['kotlinx_coroutines_k2_dev_AiIgnoreCoroutinesAccumulativeTest/test aiignore performance with repetitions on each file']"
      ></GroupProjectsChart>
    </section>
  </DashboardPage>
</template>
<script setup lang="ts">
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"
</script>
