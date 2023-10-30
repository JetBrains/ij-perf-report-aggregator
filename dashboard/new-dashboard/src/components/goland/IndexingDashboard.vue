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
        measure="numberOfIndexedFiles"
        :projects="['cockroach/indexing', 'delve/indexing', 'mattermost-server/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number Of Indexed Files: heavy projects"
        measure="numberOfIndexedFiles"
        :projects="['kubernetes/indexing', 'flux/indexing', 'istio/indexing']"
      />
    </section>

    <section>
      <GroupProjectsChart
        label="Number Of Runs Of Indexing: light projects"
        measure="numberOfRunsOfIndexing"
        :projects="['cockroach/indexing', 'delve/indexing', 'mattermost-server/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number Of Runs Of Indexing: heavy projects"
        measure="numberOfRunsOfIndexing"
        :projects="['kubernetes/indexing', 'flux/indexing', 'istio/indexing']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
