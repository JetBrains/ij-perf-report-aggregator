<template>
  <DashboardPage
    v-slot="{ dashboardConfigurators, warnings}"
    db-name="perfint"
    table="goland"
    persistent-id="goland_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <section>
      <GroupProjectsChart
        label="Indexing: Lightweight projects"
        measure="indexing"
        :projects="['flux/indexing', 'delve/indexing', 'istio/indexing']"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number Of Indexed Files: Lightweight projects"
        measure="numberOfIndexedFiles"
        :projects="['flux/indexing', 'delve/indexing', 'istio/indexing']"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Indexing: Heavyweight projects"
        measure="indexing"
        :projects="['moby/indexing', 'mattermost-server/indexing', 'cockroach/indexing', 'kubernetes/indexing']"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Number Of Indexed Files: Heavyweight projects"
        measure="numberOfIndexedFiles"
        :projects="['moby/indexing', 'mattermost-server/indexing', 'cockroach/indexing', 'kubernetes/indexing']"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Inspection execution time: Lightweight projects"
        measure="globalInspections"
        :projects="['istio/inspection', 'moby/inspection', 'flux/inspection', 'delve/inspection']"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Inspection execution time: Heavyweight projects"
        measure="globalInspections"
        :projects="['cockroach/inspection', 'kubernetes/inspection']"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Local inspection execution time"
        measure="localInspections"
        :projects="['kubernetes/localInspection', 'mattermost-server/localInspection', 'GO-5422/localInspection']"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing: average responsiveness time"
        measure="test#average_awt_delay"
        :projects="['mattermost-server/typing']"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Typing: total time"
        measure="typing"
        :projects="['mattermost-server/typing']"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages execution time"
        measure="findUsages"
        :projects="['vault/findUsages/Backend', 'vault/findUsages/List', 'vault/findUsages/Path', 'vault/findUsages/String']"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Find Usages number of found usages"
        measure="findUsages#number"
        :projects="['vault/findUsages/Backend', 'vault/findUsages/List', 'vault/findUsages/Path', 'vault/findUsages/String']"
        :configurators="dashboardConfigurators"
        :accidents="warnings"
      />
    </section>
  </DashboardPage>>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

</script>