<template>
  <DashboardPage
    db-name="perfint"
    table="goland"
    persistent-id="goland_dfa_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section>
      <GroupProjectsChart
        label="Inspection Processing Times (95 & 99 Percentile)"
        :measure="['go.dfa.inspection.file.processing.time.distribution.ms.95.percentile', 'go.dfa.inspection.file.processing.time.distribution.ms.99.percentile']"
        :projects="['kubernetes/kubernetes/dfa', 'cockroach/cockroach/dfa', 'mattermost-server/mattermost-server/dfa']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Interpretation Stage Time"
        measure="go.dfa.inspection.interpretation.stage.time.ms"
        :projects="['kubernetes/kubernetes/dfa', 'cockroach/cockroach/dfa', 'mattermost-server/mattermost-server/dfa']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Total Time"
        measure="go.dfa.inspection.total.time.ms"
        :projects="['kubernetes/kubernetes/dfa', 'cockroach/cockroach/dfa', 'mattermost-server/mattermost-server/dfa']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Resolve Function / Method Count"
        :measure="['go.dfa.resolve.function.count', 'go.dfa.resolve.method.count']"
        :projects="['kubernetes/kubernetes/dfa', 'cockroach/cockroach/dfa', 'mattermost-server/mattermost-server/dfa']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
