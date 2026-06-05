<template>
  <DashboardPage
    db-name="perfUnitTests"
    table="report"
    persistent-id="py_perf_unit_tests"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  >
    <Divider label="Pandas" />
    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsWithClientChart
          label="Completion: pandas/groupBy"
          measure="attempt.median.ms"
          :projects="[
            'com.intellij.python.junit5Tests.performance.PyPandasExamplesPerformanceTest.testCompletionInGroupBy - PyCharm',
            'com.intellij.python.junit5Tests.performance.PyPandasExamplesPerformanceTest.testCompletionInGroupBy - Pyrefly',
          ]"
          :legend-formatter="legendFormatter"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsWithClientChart
          label="Completion: pandas/merge"
          measure="attempt.median.ms"
          :projects="[
            'com.intellij.python.junit5Tests.performance.PyPandasExamplesPerformanceTest.testCompletionInMerge - PyCharm',
            'com.intellij.python.junit5Tests.performance.PyPandasExamplesPerformanceTest.testCompletionInMerge - Pyrefly',
          ]"
          :legend-formatter="legendFormatter"
        />
      </div>
    </section>
    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsWithClientChart
          label="Typing Code Analysis: pandas/groupBy"
          measure="attempt.median.ms"
          :projects="[
            'com.intellij.python.junit5Tests.performance.PyPandasExamplesPerformanceTest.testTypingCodeAnalysisGroupBy - PyCharm',
            'com.intellij.python.junit5Tests.performance.PyPandasExamplesPerformanceTest.testTypingCodeAnalysisGroupBy - Pyrefly',
          ]"
          :legend-formatter="legendFormatter"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsWithClientChart
          label="Typing Code Analysis: pandas/merge"
          measure="attempt.median.ms"
          :projects="[
            'com.intellij.python.junit5Tests.performance.PyPandasExamplesPerformanceTest.testTypingCodeAnalysisMerge - PyCharm',
            'com.intellij.python.junit5Tests.performance.PyPandasExamplesPerformanceTest.testTypingCodeAnalysisMerge - Pyrefly',
          ]"
          :legend-formatter="legendFormatter"
        />
      </div>
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="Typing-with-completion"
        measure="attempt.median.ms"
        :projects="[
          'com.intellij.python.junit5Tests.performance.PyPandasTypingPerformanceTest.testPandas233 - PyCharm',
          'com.intellij.python.junit5Tests.performance.PyPandasTypingPerformanceTest.testPandas300rc0 - PyCharm',
        ]"
      />
    </section>

    <Divider label="Boto3" />
    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsWithClientChart
          label="Completion: boto3/resource.py"
          measure="attempt.median.ms"
          :projects="['com.intellij.python.junit5Tests.performance.PyBoto3PerformanceTest.testResourceCompletion - PyCharm']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsWithClientChart
          label="Completion: boto3/client.py"
          measure="attempt.median.ms"
          :projects="['com.intellij.python.junit5Tests.performance.PyBoto3PerformanceTest.testClientCompletion - PyCharm']"
        />
      </div>
    </section>

    <Divider label="Generative" />
    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsWithClientChart
          label="Highlighting (after each type)"
          measure="attempt.median.ms"
          :projects="['com.intellij.python.junit5Tests.performance.PyHighlightingPerformanceTest.testGenerativeHighlighting - PyCharm']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsWithClientChart
          label="Typing"
          measure="attempt.median.ms"
          :projects="['com.intellij.python.junit5Tests.performance.PyTypingPerformanceTest.testGenerativeTest - PyCharm']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsWithClientChart
          label="Typing Code Analysis"
          measure="attempt.median.ms"
          :projects="['com.intellij.python.junit5Tests.performance.PyTypingCodeAnalysisPerformanceTest.testGenerativeTest - PyCharm']"
        />
      </div>
    </section>

    <Divider label="Other" />
    <section>
      <GroupProjectsWithClientChart
        label="Completion: basic (venv)"
        measure="attempt.median.ms"
        :projects="['com.intellij.python.junit5Tests.performance.PyCompletionBasicPerformanceTest.testEmptyProjectWithVenv - PyCharm']"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="Completion: overloads"
        measure="attempt.median.ms"
        :projects="['com.intellij.python.junit5Tests.performance.PyOverloadsPerformanceTest.testCompletion - PyCharm']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../common/DashboardPage.vue"
import Divider from "../common/Divider.vue"
import GroupProjectsWithClientChart from "../charts/GroupProjectsWithClientChart.vue"

const legendFormatter = (name: string) => {
  if (name.includes("- Pyrefly")) return "Pyrefly"
  if (name.includes("- PyCharm")) return "PyCharm"
  return name
}
</script>
