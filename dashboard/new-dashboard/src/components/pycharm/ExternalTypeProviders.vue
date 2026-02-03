<template>
  <DashboardPage
    db-name="perfintDev"
    table="pycharm"
    persistent-id="pycharm_dashboard"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  >
    <section>
      <GroupProjectsChart
        label="Completion"
        :measure="['fus_time_to_show_90p']"
        :projects="[
          'edx-platform (Django)/completion/testCompletionInModel/pyrefly',
          'edx-platform (Django)/completion/testCompletionInModel/pycharm',

          'edx-platform (Django)/completion/testCompletionInSettings/pyrefly',
          'edx-platform (Django)/completion/testCompletionInSettings/pycharm',

          'pandas-examples/completion/testCompletionInGroupBy/pyrefly',
          'pandas-examples/completion/testCompletionInGroupBy/pycharm',

          'pandas-examples/completion/testCompletionInMerge/pyrefly',
          'pandas-examples/completion/testCompletionInMerge/pycharm',
        ]"
      />
    </section>

    <section>
      <GroupProjectsChart
        label="Global Inspections"
        measure="globalInspections"
        :projects="[
          'edx-platform (Django)/inspection/testInspection/pyrefly',
          'edx-platform (Django)/inspection/testInspection/pycharm',

          'pandas-examples/inspection/testInspection/pyrefly',
          'pandas-examples/inspection/testInspection/pycharm',
        ]"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="Typing Code Analysis"
        measure="typingCodeAnalyzing"
        :projects="[
          'edx-platform (Django)/typing/typingCodeAnalysisCourse/pyrefly',
          'edx-platform (Django)/typing/typingCodeAnalysisCourse/pycharm',

          'edx-platform (Django)/typing/typingCodeAnalysisAdmin/pyrefly',
          'edx-platform (Django)/typing/typingCodeAnalysisAdmin/pycharm',

          'pandas-examples/typing/typingCodeAnalysisGroupBy/pyrefly',
          'pandas-examples/typing/typingCodeAnalysisGroupBy/pycharm',

          'pandas-examples/typing/typingCodeAnalysisMerge/pyrefly',
          'pandas-examples/typing/typingCodeAnalysisMerge/pycharm',
        ]"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import GroupProjectsWithClientChart from "../charts/GroupProjectsWithClientChart.vue"
</script>
