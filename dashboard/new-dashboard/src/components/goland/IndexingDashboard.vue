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
        :projects="['goland/cockroach/indexing', 'goland/delve/indexing', 'goland/mattermost/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Indexing Time Without Pauses: heavy projects"
        measure="indexingTimeWithoutPauses"
        :projects="['goland/kubernetes/indexing', 'goland/flux/indexing', 'goland/istio/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number Of Indexed Files: light projects"
        measure="numberOfIndexedFilesWritingIndexValue"
        :projects="['goland/cockroach/indexing', 'goland/delve/indexing', 'goland/mattermost/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number Of Indexed Files: heavy projects"
        measure="numberOfIndexedFilesWritingIndexValue"
        :projects="['goland/kubernetes/indexing', 'goland/flux/indexing', 'goland/istio/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Processing time for GO files: light projects"
        measure="processingTime#Go"
        :projects="['goland/cockroach/indexing', 'goland/delve/indexing', 'goland/mattermost/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Processing time for GO files: heavy projects"
        measure="processingTime#Go"
        :projects="['goland/kubernetes/indexing', 'goland/flux/indexing', 'goland/istio/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number Of Indexed Files - second IDE run"
        measure="numberOfIndexedFilesWritingIndexValue"
        :projects="[
          'goland/kubernetes/secondScanning',
          'goland/flux/secondScanning',
          'goland/istio/secondScanning',
          'goland/cockroach/secondScanning',
          'goland/delve/secondScanning',
          'goland/mattermost/secondScanning',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number Of Indexed Files - third IDE run"
        measure="numberOfIndexedFilesWritingIndexValue"
        :projects="[
          'goland/kubernetes/thirdScanning',
          'goland/flux/thirdScanning',
          'goland/istio/thirdScanning',
          'goland/cockroach/thirdScanning',
          'goland/delve/thirdScanning',
          'goland/mattermost/thirdScanning',
        ]"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
