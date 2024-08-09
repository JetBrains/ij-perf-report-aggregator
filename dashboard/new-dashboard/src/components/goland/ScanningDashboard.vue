<template>
  <DashboardPage
    db-name="perfint"
    table="goland"
    persistent-id="goland_scanning_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <section>
      <GroupProjectsChart
        label="Scanning time: first, second and third IDE runs"
        measure="scanningTimeWithoutPauses"
        :projects="['goland/kubernetes/indexing', 'goland/kubernetes/secondScanning', 'goland/kubernetes/thirdScanning']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Scanning time without pauses"
        measure="scanningTimeWithoutPauses"
        :projects="['goland/cockroach/indexing', 'goland/delve/indexing', 'goland/flux/indexing', 'goland/istio/indexing', 'goland/mattermost/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Numbers of runs of scanning"
        measure="numberOfRunsOfScannning"
        :projects="['goland/cockroach/indexing', 'goland/delve/indexing', 'goland/flux/indexing', 'goland/istio/indexing', 'goland/mattermost/indexing']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
