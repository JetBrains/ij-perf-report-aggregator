<template>
  <DashboardPage
    db-name="perfintDev"
    table="goland"
    :with-installer="false"
    persistent-id="goland_scanning_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section>
      <GroupProjectsChart
        label="Scanning time: first, second and third IDE runs (kubernetes)"
        measure="scanningTimeWithoutPauses"
        :projects="['kubernetes/indexing', 'kubernetes/secondScanning', 'kubernetes/thirdScanning']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Scanning time without pauses"
        measure="scanningTimeWithoutPauses"
        :projects="['cockroach/indexing', 'delve/indexing', 'flux/indexing', 'istio/indexing', 'mattermost/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Numbers of runs of scanning"
        measure="numberOfRunsOfScannning"
        :projects="['cockroach/indexing', 'delve/indexing', 'flux/indexing', 'istio/indexing', 'mattermost/indexing']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
