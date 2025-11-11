<template>
  <DashboardPage
    db-name="perfintDev"
    table="goland"
    :with-installer="false"
    persistent-id="goland_code_analyzes_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section>
      <GroupProjectsChart
        label="File Analysis on Open"
        measure="firstCodeAnalysis"
        :projects="[
          'istio/localInspection/adsc.go',
          'minotaur/localInspection/server.go',
          'caddy/localInspection/replacer.go',
          'cockroach/localInspection/tochar.go',
          'flux/localInspection/sourcesecret.go',
          'kubernetes/localInspection/client_test.go',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Global Inspections"
        measure="globalInspections"
        :projects="[
          'delve/inspection',
          'kubernetes/inspection',
          'caddy/inspection',
          'cockroach/inspection',
          'k8sDevice/inspection',
          'mattermost-server/inspection',
          'milvus/inspection',
          'rclone/inspection',
          'tempo/inspection',
          'volcano/inspection',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Global Inspections without DFA"
        measure="globalInspections"
        :projects="[
          'cockroach/noDfa/inspection',
          'k8sDevice/noDfa/inspection',
          'kubernetes/noDfa/inspection',
          'mattermost-server/noDfa/inspection',
          'milvus/noDfa/inspection',
          'rclone/noDfa/inspection',
          'tempo/noDfa/inspection',
          'volcano/noDfa/inspection',
        ]"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
