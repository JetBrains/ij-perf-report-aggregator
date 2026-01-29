<template>
  <DashboardPage
    v-slot="{ averagesConfigurators }"
    db-name="perfintDev"
    table="kotlin"
    persistent-id="kotlinDev_completion_dashboard"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  >
    <ConfiguratorRegistration
      :configurator="projectConfigurator"
      :data="Object.values(PROJECT_CATEGORIES).flatMap((c) => c.label)"
    />
    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'completion\_%'"
          :aggregated-project="'%\_k1'"
          :is-like="true"
          :title="'mean all project completion K1'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'completion\_%'"
          :aggregated-project="'%\_k2'"
          :is-like="true"
          :title="'mean all project completion K2'"
        />
      </div>
    </section>

    <K1K2DashboardGroupCharts :definitions="completionCharts" />
  </DashboardPage>
</template>

<script setup lang="ts">
import AggregationChart from "../../charts/AggregationChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"
import K1K2DashboardGroupCharts from "../K1K2DashboardGroupCharts.vue"
import { createKotlinCharts, PROJECT_CATEGORIES } from "../projects"
import { SimpleMeasureConfigurator } from "../../../configurators/SimpleMeasureConfigurator"
import ConfiguratorRegistration from "../ConfiguratorRegistration.vue"

const projectConfigurator = new SimpleMeasureConfigurator("project", null)
const { completionCharts } = createKotlinCharts(projectConfigurator)
</script>
