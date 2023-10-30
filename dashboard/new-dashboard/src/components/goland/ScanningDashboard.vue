<template>
  <DashboardPage
    db-name="perfint"
    table="goland"
    persistent-id="goland_scanning_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <section>
      <GroupProjectsChart
        label="Scanning time for the first, second and third IDE runs (kubernetes project)"
        measure="scanningTimeWithoutPauses"
        :projects="['kubernetes/indexing', 'kubernetes/secondScanning', 'kubernetes/thirdScanning']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Scanning time without pauses"
        measure="scanningTimeWithoutPauses"
        :projects="['cockroach/indexing', 'delve/indexing', 'flux/indexing', 'istio/indexing', 'mattermost-server/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Numbers of runs of scanning"
        measure="numberOfRunsOfScannning"
        :projects="['cockroach/indexing', 'delve/indexing', 'flux/indexing', 'istio/indexing', 'mattermost-server/indexing']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
