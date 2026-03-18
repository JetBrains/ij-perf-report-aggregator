<template>
  <DashboardPage
    db-name="perfUnitTests"
    table="report"
    persistent-id="pycharm_dashboard"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  >
    <section>
      <GroupProjectsWithClientChart
        label="Big pandas file"
        measure="attempt.median.ms"
        :projects="['com.intellij.python.junit5Tests.performance.typeEngine.PyTypeEnginePerformanceTest.testBigPandasFile - Pyrefly', 'com.intellij.python.junit5Tests.performance.typeEngine.PyTypeEnginePerformanceTest.testBigPandasFile - PyCharm']"
        :legend-formatter="legendFormatter"
      />
    </section>

    <section>
      <GroupProjectsWithClientChart
        label="Pandas DataFrame.mean() method"
        measure="attempt.median.ms"
        :projects="['com.intellij.python.junit5Tests.performance.typeEngine.PyTypeEnginePerformanceTest.complex type inference performance - Pyrefly', 'com.intellij.python.junit5Tests.performance.typeEngine.PyTypeEnginePerformanceTest.complex type inference performance - PyCharm']"
        :legend-formatter="legendFormatter"
      />
    </section>

    <section>
      <GroupProjectsWithClientChart
        label="Pandas DataFrame class"
        measure="attempt.median.ms"
        :projects="['com.intellij.python.junit5Tests.performance.typeEngine.PyTypeEnginePerformanceTest.data frame type inference performance - Pyrefly', 'com.intellij.python.junit5Tests.performance.typeEngine.PyTypeEnginePerformanceTest.data frame type inference performance - PyCharm']"
        :legend-formatter="legendFormatter"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../common/DashboardPage.vue"
import GroupProjectsWithClientChart from "../charts/GroupProjectsWithClientChart.vue"

const legendFormatter = (name: string) => {
  if (name.includes("- Pyrefly")) {
    return "Pyrefly"
  }
  if (name.includes("- PyCharm")) {
    return "PyCharm"
  }
  if (name.includes("- ty")) {
    return "ty"
  }
  return name
}
</script>
