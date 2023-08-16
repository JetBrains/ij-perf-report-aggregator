<template>
  <DashboardPage
    v-slot="{ averagesConfigurators }"
    db-name="perfint"
    table="kotlin"
    persistent-id="kotlin_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'completion\_%'"
          :aggregated-project="'%\_k1'"
          :is-like="true"
          :title="'completion K1'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'completion\_%'"
          :aggregated-project="'%\_k2'"
          :is-like="true"
          :title="'completion K2'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'highlighting\_%'"
          :aggregated-project="'%\_k1'"
          :is-like="true"
          :title="'highlighting K1'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'highlighting\_%'"
          :aggregated-project="'%\_k2'"
          :is-like="true"
          :title="'highlighting K2'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'findUsage\_%'"
          :aggregated-project="'%\_k1'"
          :is-like="true"
          :title="'findUsage K1'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'findUsage\_%'"
          :aggregated-project="'%\_k2'"
          :is-like="true"
          :title="'findUsage K2'"
        />
      </div>
    </section>
    <DashboardGroupCharts :definitions="completionCharts" />
    <DashboardGroupCharts :definitions="localInspectionsCharts" />
    <DashboardGroupCharts :definitions="highlightingCharts" />
    <DashboardGroupCharts :definitions="findUsagesCharts" />
    <DashboardGroupCharts :definitions="refactoringCharts" />
  </DashboardPage>
</template>

<script setup lang="ts">
import AggregationChart from "../charts/AggregationChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import DashboardGroupCharts from "./DashboardGroupCharts.vue"
import { completionCharts, findUsagesCharts, highlightingCharts, localInspectionsCharts, refactoringCharts } from "./projects"
</script>
