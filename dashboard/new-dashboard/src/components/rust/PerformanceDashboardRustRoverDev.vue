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
        label="UI Lags in Local Inspection: Test Itself (metric 'test#max_awt_delay')"
        measure="test#max_awt_delay"
        :projects="rustLocalInspectionCases"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="UI Lags in Local Inspection: Whole IDE Run (metric 'ide#max_awt_delay')"
        measure="ide#max_awt_delay"
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
        label="UI Lags in Global Inspection: Test Itself (metric 'test#max_awt_delay')"
        measure="test#max_awt_delay"
        :projects="rustGlobalInspectionProjects.map((project) => `global-inspection/${project}-inspection`)"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="UI Lags in Global Inspection: Whole IDE Run (metric 'ide#max_awt_delay')"
        measure="ide#max_awt_delay"
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
        label="UI Lags in Completion: Test Itself (metric 'test#max_awt_delay')"
        measure="test#max_awt_delay"
        :projects="rustCompletionCases"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="UI Lags in Completion: Whole IDE Run (metric 'ide#max_awt_delay')"
        measure="ide#max_awt_delay"
        :projects="rustCompletionCases"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages"
        measure="findUsages"
        :projects="rustFindUsagesCases"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="UI Lags in Find Usages: Test Itself (metric 'test#max_awt_delay')"
        measure="test#max_awt_delay"
        :projects="rustFindUsagesCases"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="UI Lags in Find Usages: Whole IDE Run (metric 'ide#max_awt_delay')"
        measure="ide#max_awt_delay"
        :projects="rustFindUsagesCases"
      />
    </section>

    <section>
      <GroupProjectsChart
        label="Typing Latency (mean value)"
        measure="typing#latency#mean_value"
        :projects="rustTypingCases"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing Latency (max value)"
        measure="typing#latency#max"
        :projects="rustTypingCases"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="UI Lags in Typing: Test Itself (metric 'test#max_awt_delay')"
        measure="test#max_awt_delay"
        :projects="rustTypingCases"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="UI Lags in Typing: Whole IDE Run (metric 'ide#max_awt_delay')"
        measure="ide#max_awt_delay"
        :projects="rustTypingCases"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import { rustLocalInspectionCases, rustGlobalInspectionProjects, rustCompletionCases, rustFindUsagesCases, rustTypingCases } from "./RustTestCases"
</script>
