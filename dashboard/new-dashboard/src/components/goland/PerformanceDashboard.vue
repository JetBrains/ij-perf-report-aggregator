<template>
  <DashboardPage
    db-name="perfint"
    table="goland"
    persistent-id="goland_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <section>
      <GroupProjectsChart
        label="Indexing: Lightweight projects"
        measure="indexingWithoutPauses"
        :projects="['flux/indexing', 'delve/indexing', 'istio/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number Of Indexed Files: Lightweight projects"
        measure="numberOfIndexedFiles"
        :projects="['flux/indexing', 'delve/indexing', 'istio/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Indexing: Heavyweight projects"
        :measure="['indexingWithoutPauses', 'indexing']"
        :projects="['moby/indexing', 'mattermost-server/indexing', 'cockroach/indexing', 'kubernetes/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number Of Indexed Files: Heavyweight projects"
        measure="numberOfIndexedFiles"
        :projects="['moby/indexing', 'mattermost-server/indexing', 'cockroach/indexing', 'kubernetes/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Inspection execution time: Lightweight projects"
        measure="globalInspections"
        :projects="['istio/inspection', 'moby/inspection', 'flux/inspection', 'delve/inspection']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Inspection execution time: Heavyweight projects"
        measure="globalInspections"
        :projects="['cockroach/inspection', 'kubernetes/inspection']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Local inspection execution time"
        measure="localInspections"
        :projects="['kubernetes/localInspection', 'mattermost-server/localInspection', 'GO-5422/localInspection']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing: average responsiveness time"
        measure="test#average_awt_delay"
        :projects="['mattermost-server/typing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing: total time"
        measure="typing"
        :projects="['mattermost-server/typing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages execution time"
        measure="findUsages"
        :projects="['vault/findUsages/Backend', 'vault/findUsages/List', 'vault/findUsages/Path', 'vault/findUsages/String']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages number of found usages"
        measure="findUsages#number"
        :projects="['vault/findUsages/Backend', 'vault/findUsages/List', 'vault/findUsages/Path', 'vault/findUsages/String']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
