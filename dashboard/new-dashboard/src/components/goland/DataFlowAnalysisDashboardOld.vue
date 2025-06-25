<template>
  <DashboardPage
    db-name="perfint"
    table="goland"
    persistent-id="goland_dfa_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section>
      <GroupProjectsChart
        label="General Total Time"
        measure="go.dfa.general.total.time.ms"
        :projects="[
          'kubernetes/kubernetes/dfa',
          'cockroach/cockroach/dfa',
          'mattermost-server/mattermost-server/dfa',
          'volcano/volcano/dfa',
          'tempo/tempo/dfa',
          'rclone/rclone/dfa',
          'milvus/milvus/dfa',
          'k8sDevice/k8sDevice/dfa',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="General Average Time"
        measure="go.dfa.general.avg.time.ms"
        :projects="[
          'kubernetes/kubernetes/dfa',
          'cockroach/cockroach/dfa',
          'mattermost-server/mattermost-server/dfa',
          'volcano/volcano/dfa',
          'tempo/tempo/dfa',
          'rclone/rclone/dfa',
          'milvus/milvus/dfa',
          'k8sDevice/k8sDevice/dfa',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="General Average Time Without Summary Load Time"
        measure="go.dfa.general.avg.without.summary.load.time.ms"
        :projects="[
          'kubernetes/kubernetes/dfa',
          'cockroach/cockroach/dfa',
          'mattermost-server/mattermost-server/dfa',
          'volcano/volcano/dfa',
          'tempo/tempo/dfa',
          'rclone/rclone/dfa',
          'milvus/milvus/dfa',
          'k8sDevice/k8sDevice/dfa',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="General Computed File Gists Count"
        measure="go.dfa.general.computed.file.gists.count"
        :projects="[
          'kubernetes/kubernetes/dfa',
          'cockroach/cockroach/dfa',
          'mattermost-server/mattermost-server/dfa',
          'volcano/volcano/dfa',
          'tempo/tempo/dfa',
          'rclone/rclone/dfa',
          'milvus/milvus/dfa',
          'k8sDevice/k8sDevice/dfa',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="General Files Count"
        measure="go.dfa.general.files.count"
        :projects="[
          'kubernetes/kubernetes/dfa',
          'cockroach/cockroach/dfa',
          'mattermost-server/mattermost-server/dfa',
          'volcano/volcano/dfa',
          'tempo/tempo/dfa',
          'rclone/rclone/dfa',
          'milvus/milvus/dfa',
          'k8sDevice/k8sDevice/dfa',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="General Functions Count"
        measure="go.dfa.general.functions.count"
        :projects="[
          'kubernetes/kubernetes/dfa',
          'cockroach/cockroach/dfa',
          'mattermost-server/mattermost-server/dfa',
          'volcano/volcano/dfa',
          'tempo/tempo/dfa',
          'rclone/rclone/dfa',
          'milvus/milvus/dfa',
          'k8sDevice/k8sDevice/dfa',
        ]"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
