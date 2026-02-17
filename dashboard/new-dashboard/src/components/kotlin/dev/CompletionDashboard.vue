<template>
  <DashboardPage
    v-slot="{ averagesConfigurators }"
    db-name="perfintDev"
    table="kotlin"
    persistent-id="kotlinDev_dashboard"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  >
    <section class="flex-1 min-w-0">
      <AggregationChart
        :configurators="averagesConfigurators"
        :aggregated-measure="'completion\_%'"
        :aggregated-project="'%\_k2'"
        :is-like="true"
        :title="'mean all project completion'"
      />
    </section>

    <KotlinDashboardGroupCharts :definitions="completionCharts" />
  </DashboardPage>
</template>

<script setup lang="ts">
import AggregationChart from "../../charts/AggregationChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"
import KotlinDashboardGroupCharts from "../KotlinDashboardGroupCharts.vue"
import { completionCharts } from "../projects"
</script>
