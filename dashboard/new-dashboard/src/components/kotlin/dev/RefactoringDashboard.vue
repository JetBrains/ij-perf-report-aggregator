<template>
  <DashboardPage
    v-slot="{ averagesConfigurators }"
    db-name="perfintDev"
    table="kotlin"
    persistent-id="kotlinDev_dashboard"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  >
    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'performInlineRename\_%'"
          :aggregated-project="'%\_k1'"
          :is-like="true"
          :title="'mean all rename K1'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'performInlineRename\_%'"
          :aggregated-project="'%\_k2'"
          :is-like="true"
          :title="'mean all rename K2'"
        />
      </div>

      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'prepareForRename\_%'"
          :aggregated-project="'%\_k1'"
          :is-like="true"
          :title="'mean all prepare-rename K1'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'performInlineRename\_%'"
          :aggregated-project="'%\_k2'"
          :is-like="true"
          :title="'mean all prepare-rename K2'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'startInlineRename\_%'"
          :aggregated-project="'%\_k1'"
          :is-like="true"
          :title="'mean all start-rename K1'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'startInlineRename\_%'"
          :aggregated-project="'%\_k2'"
          :is-like="true"
          :title="'mean all start-rename K2'"
        />
      </div>
    </section>

    <DashboardGroupCharts :definitions="refactoringCharts" />
  </DashboardPage>
</template>

<script setup lang="ts">
import AggregationChart from "../../charts/AggregationChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"
import DashboardGroupCharts from "./DashboardGroupCharts.vue"
import { refactoringCharts } from "./projects"
</script>
