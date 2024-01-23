<template>
  <DashboardPage
    db-name="perfint"
    table="goland"
    persistent-id="goland_indexing_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <section>
      <GroupProjectsChart
        label="Indexing Time Without Pauses: light projects"
        measure="indexingTimeWithoutPauses"
        :projects="['cockroach/indexing', 'delve/indexing', 'mattermost-server/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Indexing Time Without Pauses: heavy projects"
        measure="indexingTimeWithoutPauses"
        :projects="['kubernetes/indexing', 'flux/indexing', 'istio/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number Of Indexed Files: light projects"
        measure="numberOfIndexedFilesWritingIndexValue"
        :projects="['cockroach/indexing', 'delve/indexing', 'mattermost-server/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number Of Indexed Files: heavy projects"
        measure="numberOfIndexedFilesWritingIndexValue"
        :projects="['kubernetes/indexing', 'flux/indexing', 'istio/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Processing time for GO files: light projects"
        measure="processingTime#Go"
        :projects="['cockroach/indexing', 'delve/indexing', 'mattermost-server/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Processing time for GO files: heavy projects"
        measure="processingTime#Go"
        :projects="['kubernetes/indexing', 'flux/indexing', 'istio/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number Of Indexed Files - second IDE run"
        measure="numberOfIndexedFilesWritingIndexValue"
        :projects="[
          'kubernetes/secondScanning',
          'flux/secondScanning',
          'istio/secondScanning',
          'cockroach/secondScanning',
          'delve/secondScanning',
          'mattermost-server/secondScanning',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number Of Indexed Files - third IDE run"
        measure="numberOfIndexedFilesWritingIndexValue"
        :projects="['kubernetes/thirdScanning', 'flux/thirdScanning', 'istio/thirdScanning', 'cockroach/thirdScanning', 'delve/thirdScanning', 'mattermost-server/thirdScanning']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
