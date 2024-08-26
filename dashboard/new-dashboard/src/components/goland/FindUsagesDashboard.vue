<template>
  <DashboardPage
    db-name="perfint"
    table="goland"
    persistent-id="goland_findusages_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <section>
      <GroupProjectsChart
        label="Find Usages (vault): Execution Time"
        measure="findUsages"
        :projects="['findUsages/unmarshalJSON', 'findUsages/path', 'findUsages/backend', 'findUsages/string', 'findUsages/list']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages (vault): First Usage Time"
        measure="findUsages_firstUsage"
        :projects="['findUsages/unmarshalJSON', 'findUsages/path', 'findUsages/backend', 'findUsages/string', 'findUsages/list']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages (vault): Number of Usages"
        measure="findUsages#number"
        :projects="['findUsages/unmarshalJSON', 'findUsages/path', 'findUsages/backend', 'findUsages/string', 'findUsages/list']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages (pocketbase): Execution Time"
        measure="findUsages"
        :projects="['findUsages/write(method)', 'findUsages/start(method)', 'findUsages/open(method)', 'findUsages/close(method)', 'findUsages/file(struct)']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages (pocketbase): First Usage Time"
        measure="findUsages_firstUsage"
        :projects="['findUsages/write(method)', 'findUsages/start(method)', 'findUsages/open(method)', 'findUsages/close(method)', 'findUsages/file(struct)']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages (pocketbase): Number of Usages"
        measure="findUsages#number"
        :projects="['findUsages/write(method)', 'findUsages/start(method)', 'findUsages/open(method)', 'findUsages/close(method)', 'findUsages/file(struct)']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
