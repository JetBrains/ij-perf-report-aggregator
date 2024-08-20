<template>
  <DashboardPage
    db-name="perfintDev"
    table="ml"
    :initial-machine="null"
    persistent-id="api-tests"
    :with-installer="false"
  >
    <section>
      <GroupProjectsChart
        label="Current token quota usage %"
        :measure="['currentPercents']"
        :projects="['stage: tokenQuota', 'prod: tokenQuota']"
      />
    </section>
    <ApiTestsDashboardGroupCharts :definitions="nameSuggestProjects" />
    <ApiTestsDashboardGroupCharts :definitions="testProjects" />
    <ApiTestsDashboardGroupCharts :definitions="timingProjects" />
  </DashboardPage>
</template>
<script setup lang="ts">
import DashboardPage from "../../common/DashboardPage.vue"
import ApiTestsDashboardGroupCharts from "./components/ApiTestsDashboardGroupCharts.vue"
import { nameSuggestProjects, testProjects, timingProjects } from "./api-tests-projects"
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
</script>
