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
        :aggregated-measure="'performInlineRename\_%'"
        :aggregated-project="'%\_k2'"
        :is-like="true"
        :title="'mean all rename'"
      />
      <AggregationChart
        :configurators="averagesConfigurators"
        :aggregated-measure="'performInlineRename\_%'"
        :aggregated-project="'%\_k2'"
        :is-like="true"
        :title="'mean all prepare-rename'"
      />
      <AggregationChart
        :configurators="averagesConfigurators"
        :aggregated-measure="'startInlineRename\_%'"
        :aggregated-project="'%\_k2'"
        :is-like="true"
        :title="'mean all start-rename'"
      />
    </section>

    <KotlinDashboardGroupCharts :definitions="refactoringCharts" />
  </DashboardPage>
</template>

<script setup lang="ts">
import AggregationChart from "../../charts/AggregationChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"
import KotlinDashboardGroupCharts from "../KotlinDashboardGroupCharts.vue"
import { refactoringCharts } from "../projects"
</script>
