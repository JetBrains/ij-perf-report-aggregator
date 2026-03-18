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
      />
    </section>

    <section>
      <GroupProjectsWithClientChart
        label="Pandas DataFrame.mean() method"
        measure="attempt.median.ms"
        :projects="['com.intellij.python.junit5Tests.performance.typeEngine.PyTypeEnginePerformanceTest.complex type inference performance - Pyrefly', 'com.intellij.python.junit5Tests.performance.typeEngine.PyTypeEnginePerformanceTest.complex type inference performance - PyCharm']"
      />
    </section>

    <section>
      <GroupProjectsWithClientChart
        label="Pandas DataFrame class"
        measure="attempt.median.ms"
        :projects="['com.intellij.python.junit5Tests.performance.typeEngine.PyTypeEnginePerformanceTest.data frame type inference performance - Pyrefly', 'com.intellij.python.junit5Tests.performance.typeEngine.PyTypeEnginePerformanceTest.data frame type inference performance - PyCharm']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import GroupProjectsWithClientChart from "../charts/GroupProjectsWithClientChart.vue"
</script>
