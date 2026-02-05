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
        label="Global Inspections (edx-platform)"
        measure="globalInspections"
        :projects="['edx-platform (Django)/inspection/testInspection/pyrefly', 'edx-platform (Django)/inspection/testInspection/pycharm']"
      />
      <GroupProjectsChart
        label="Global Inspections (pandas-examples)"
        measure="globalInspections"
        :projects="['pandas-examples/inspection/testInspection/pyrefly', 'pandas-examples/inspection/testInspection/pycharm']"
      />
    </section>
    <section>
      <GroupProjectsWithClientChart
        label="Typing Code Analysis (edx-platform/course)"
        measure="typingCodeAnalyzing"
        :projects="['edx-platform (Django)/typing/testTypingCodeAnalysisCourse/pyrefly', 'edx-platform (Django)/typing/testTypingCodeAnalysisCourse/pycharm']"
      />
      <GroupProjectsWithClientChart
        label="Typing Code Analysis (edx-platform/admin)"
        measure="typingCodeAnalyzing"
        :projects="['edx-platform (Django)/typing/testTypingCodeAnalysisAdmin/pyrefly', 'edx-platform (Django)/typing/testTypingCodeAnalysisAdmin/pycharm']"
      />
      <GroupProjectsWithClientChart
        label="Typing Code Analysis (pandas-examples/groupBy)"
        measure="typingCodeAnalyzing"
        :projects="['pandas-examples/typing/testTypingCodeAnalysisGroupBy/pyrefly', 'pandas-examples/typing/testTypingCodeAnalysisGroupBy/pycharm']"
      />
      <GroupProjectsWithClientChart
        label="Typing Code Analysis (pandas-examples/merge)"
        measure="typingCodeAnalyzing"
        :projects="['pandas-examples/typing/testTypingCodeAnalysisMerge/pyrefly', 'pandas-examples/typing/testTypingCodeAnalysisMerge/pycharm']"
      />
    </section>

    <section>
      <GroupProjectsChart
        label="Completion (edx-platform/model)"
        :measure="['fus_time_to_show_90p']"
        :projects="['edx-platform (Django)/completion/testCompletionInModel/pyrefly', 'edx-platform (Django)/completion/testCompletionInModel/pycharm']"
      />
      <GroupProjectsChart
        label="Completion (edx-platform/settings)"
        :measure="['fus_time_to_show_90p']"
        :projects="['edx-platform (Django)/completion/testCompletionInSettings/pyrefly', 'edx-platform (Django)/completion/testCompletionInSettings/pycharm']"
      />
      <GroupProjectsChart
        label="Completion (pandas-examples/groupBy)"
        :measure="['fus_time_to_show_90p']"
        :projects="['pandas-examples/completion/testCompletionInGroupBy/pyrefly', 'pandas-examples/completion/testCompletionInGroupBy/pycharm']"
      />
      <GroupProjectsChart
        label="Completion (pandas-examples/merge)"
        :measure="['fus_time_to_show_90p']"
        :projects="['pandas-examples/completion/testCompletionInMerge/pyrefly', 'pandas-examples/completion/testCompletionInMerge/pycharm']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import GroupProjectsWithClientChart from "../charts/GroupProjectsWithClientChart.vue"
</script>
