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
        :projects="['goland/findUsages/unmarshalJSON', 'goland/findUsages/path', 'goland/findUsages/backend', 'goland/findUsages/string', 'goland/findUsages/list']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages (vault): First Usage Time"
        measure="findUsages_firstUsage"
        :projects="['goland/findUsages/unmarshalJSON', 'goland/findUsages/path', 'goland/findUsages/backend', 'goland/findUsages/string', 'goland/findUsages/list']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages (vault): Number of Usages"
        measure="findUsages#number"
        :projects="['goland/findUsages/unmarshalJSON', 'goland/findUsages/path', 'goland/findUsages/backend', 'goland/findUsages/string', 'goland/findUsages/list']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages (pocketbase): Execution Time"
        measure="findUsages"
        :projects="[
          'goland/findUsages/write(method)',
          'goland/findUsages/start(method)',
          'goland/findUsages/open(method)',
          'goland/findUsages/close(method)',
          'goland/findUsages/file(struct)',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages (pocketbase): First Usage Time"
        measure="findUsages_firstUsage"
        :projects="[
          'goland/findUsages/write(method)',
          'goland/findUsages/start(method)',
          'goland/findUsages/open(method)',
          'goland/findUsages/close(method)',
          'goland/findUsages/file(struct)',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages (pocketbase): Number of Usages"
        measure="findUsages#number"
        :projects="[
          'goland/findUsages/write(method)',
          'goland/findUsages/start(method)',
          'goland/findUsages/open(method)',
          'goland/findUsages/close(method)',
          'goland/findUsages/file(struct)',
        ]"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
