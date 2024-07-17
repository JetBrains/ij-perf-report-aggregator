<template>
  <DashboardPage
    db-name="perfint"
    table="datagrip"
    persistent-id="datagrip_indexing_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <section>
      <GroupProjectsChart
        label="Indexing Time Without Pauses"
        measure="indexingTimeWithoutPauses"
        :projects="['empty_project/indexing', 'dumps_project/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Scanning Time Without Pauses"
        measure="scanningTimeWithoutPauses"
        :projects="['empty_project/indexing', 'dumps_project/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number Of Indexed Files"
        measure="numberOfIndexedFilesWritingIndexValue"
        :projects="['empty_project/indexing', 'dumps_project/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Processing Time: SQL"
        measure="processingTime#SQL"
        :projects="['empty_project/indexing', 'dumps_project/indexing']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
