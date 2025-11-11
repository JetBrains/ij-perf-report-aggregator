<template>
  <DashboardPage
    db-name="perfintDev"
    :with-installer="false"
    table="goland"
    persistent-id="goland_indexing_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
  >
    <section>
      <GroupProjectsChart
        label="Indexing Time"
        measure="indexingTimeWithoutPauses"
        :projects="['cockroach/indexing', 'delve/indexing', 'mattermost/indexing', 'kubernetes/indexing', 'flux/indexing', 'istio/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Index Size"
        measure="indexingSize"
        :projects="['cockroach/indexing', 'delve/indexing', 'mattermost/indexing', 'kubernetes/indexing', 'flux/indexing', 'istio/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Processing Time"
        measure="processingTime#Go"
        :projects="['kubernetes/indexing', 'flux/indexing', 'istio/indexing', 'cockroach/indexing', 'delve/indexing', 'mattermost/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Processing Speed"
        measure="processingSpeedAvg#Go"
        :projects="['kubernetes/indexing', 'flux/indexing', 'istio/indexing', 'cockroach/indexing', 'delve/indexing', 'mattermost/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number Of Indexed Files"
        measure="numberOfIndexedFilesWritingIndexValue"
        :projects="['cockroach/indexing', 'delve/indexing', 'mattermost/indexing', 'kubernetes/indexing', 'flux/indexing', 'istio/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number Of Indexed Files - second IDE run"
        measure="numberOfIndexedFilesWritingIndexValue"
        :projects="['kubernetes/secondScanning', 'flux/secondScanning', 'istio/secondScanning', 'cockroach/secondScanning', 'delve/secondScanning', 'mattermost/secondScanning']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number Of Indexed Files - third IDE run"
        measure="numberOfIndexedFilesWritingIndexValue"
        :projects="['kubernetes/thirdScanning', 'flux/thirdScanning', 'istio/thirdScanning', 'cockroach/thirdScanning', 'delve/thirdScanning', 'mattermost/thirdScanning']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Scanning Time"
        measure="scanningTimeWithoutPauses"
        :projects="['cockroach/indexing', 'delve/indexing', 'mattermost/indexing', 'kubernetes/indexing', 'flux/indexing', 'istio/indexing']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
