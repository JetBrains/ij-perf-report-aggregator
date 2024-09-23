<template>
  <DashboardPage
    db-name="perfint"
    table="goland"
    persistent-id="goland_findusages_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <section>
      <GroupProjectsChart
        label="Total Execution Time"
        measure="findUsages"
        :projects="[
          'vault/backend(interface)',
          'vault/list(method)',
          'vault/path(struct)',
          'vault/string(method)',
          'vault/unmarshalJSON(method)',
          'pocketbase/write(method)',
          'pocketbase/start(method)',
          'pocketbase/open(method)',
          'pocketbase/file(struct)',
          'pocketbase/close(method)',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="First Usage Time"
        measure="findUsages_firstUsage"
        :projects="[
          'vault/backend(interface)',
          'vault/list(method)',
          'vault/path(struct)',
          'vault/string(method)',
          'vault/unmarshalJSON(method)',
          'pocketbase/write(method)',
          'pocketbase/start(method)',
          'pocketbase/open(method)',
          'pocketbase/file(struct)',
          'pocketbase/close(method)',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number of Usages"
        measure="findUsages#number"
        :projects="[
          'vault/backend(interface)',
          'vault/list(method)',
          'vault/path(struct)',
          'vault/string(method)',
          'vault/unmarshalJSON(method)',
          'pocketbase/write(method)',
          'pocketbase/start(method)',
          'pocketbase/open(method)',
          'pocketbase/file(struct)',
          'pocketbase/close(method)',
        ]"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
