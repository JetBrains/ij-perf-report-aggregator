<template>
  <DashboardPage
    db-name="perfintDev"
    table="rust"
    persistent-id="rust_plugin_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :with-installer="false"
  >
    <section>
      <GroupProjectsChart
        label="Local Inspections (on file open, metric 'firstCodeAnalysis')"
        measure="firstCodeAnalysis"
        :projects="rustLocalInspectionCases"
      />
    </section>

    <section>
      <GroupProjectsChart
        label="Local Inspections (on typing top-level, metric 'typingCodeAnalyzing#mean_value')"
        measure="typingCodeAnalyzing#mean_value"
        :projects="rustLocalInspectionCases.map((testCase) => `${testCase}-top-level-typing`)"
      />
    </section>

    <section>
      <GroupProjectsChart
        label="Local Inspections (on typing stmt in function, metric 'typingCodeAnalyzing#mean_value')"
        measure="typingCodeAnalyzing#mean_value"
        :projects="rustLocalInspectionCases"
      />
    </section>

    <section>
      <GroupProjectsChart
        label="Global Inspection execution time (metric 'globalInspections')"
        measure="globalInspections"
        :projects="rustGlobalInspectionProjects.map((project) => `global-inspection/${project}-inspection`)"
      />
    </section>

    <section>
      <GroupProjectsChart
        label="Completion"
        measure="completion#mean_value"
        :projects="rustCompletionCases"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages"
        measure="findUsages"
        :projects="['find-usages/yew', 'find-usages/wasm']"
      />
    </section>

    <section>
      <GroupProjectsChart
        label="Typing Latency (mean value)"
        measure="typing#latency#mean_value"
        :projects="['typing/nalgebra/typing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing Latency (max value)"
        measure="typing#latency#max"
        :projects="['typing/nalgebra/typing']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import { rustLocalInspectionCases, rustGlobalInspectionProjects, rustCompletionCases } from "./RustTestCases"
</script>
