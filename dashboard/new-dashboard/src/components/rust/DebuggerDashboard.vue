<template>
  <DashboardPage
    db-name="perfintDev"
    table="rust"
    persistent-id="rust_debugger_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :with-installer="false">
    <section>
      <GroupProjectsChart
        label="Launch Debug"
        measure="debugRunConfiguration"
        :projects=rustDebuggerTestCases
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Step Into"
        measure="debugStep_into"
        :projects=rustDebuggerTestCases
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Step Out"
        measure="debugStep_out"
        :projects=rustDebuggerTestCases
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Step Over"
        measure="debugStep_over"
        :projects=rustDebuggerTestCases
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Evaluate Expression"
        measure="evaluateExpression"
        :projects=rustDebuggerTestCases
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import { rustDebuggerTestCases } from "./RustTestCases"
</script>
