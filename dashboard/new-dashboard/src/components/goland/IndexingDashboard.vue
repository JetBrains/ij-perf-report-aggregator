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
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Indexed Files"
          measure="numberOfIndexedFilesWritingIndexValue"
          :projects="['cockroach/indexing', 'delve/indexing', 'mattermost/indexing', 'kubernetes/indexing', 'flux/indexing', 'istio/indexing']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Index Size"
          measure="indexSize"
          :projects="['cockroach/indexing', 'delve/indexing', 'mattermost/indexing', 'kubernetes/indexing', 'flux/indexing', 'istio/indexing']"
        />
      </div>
    </section>
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Processing Time"
          measure="processingTime#Go"
          :projects="['kubernetes/indexing', 'flux/indexing', 'istio/indexing', 'cockroach/indexing', 'delve/indexing', 'mattermost/indexing']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Processing Speed"
          measure="processingSpeedAvg#Go"
          :projects="['kubernetes/indexing', 'flux/indexing', 'istio/indexing', 'cockroach/indexing', 'delve/indexing', 'mattermost/indexing']"
        />
      </div>
    </section>
    <section>
      <GroupProjectsChart
        label="GC Pause, ms"
        measure="gcPause"
        :projects="['cockroach/indexing', 'delve/indexing', 'mattermost/indexing', 'kubernetes/indexing', 'flux/indexing', 'istio/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="GC Memory Collected, Mb"
        measure="freedMemoryByGC"
        :projects="['cockroach/indexing', 'delve/indexing', 'mattermost/indexing', 'kubernetes/indexing', 'flux/indexing', 'istio/indexing']"
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
