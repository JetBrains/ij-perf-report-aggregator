<template>
  <DashboardPage
    db-name="jbr"
    table="report"
    persistent-id="jbr_team_dashboard"
    :with-installer="false"
    :is-build-number-exists="true"
  >
    <div
      v-for="metric in metricsNames"
      :key="metric"
    >
      <div class="relative flex py-5 items-center">
        <div class="grow border-t" />
        <span class="shrink mx-4 text-lg">{{ metric }}</span>
        <div class="grow border-t" />
      </div>
      <section>
        <GroupProjectsChart
          label="macOS"
          :measure="metric"
          :projects="macOSConfigurations"
        />
      </section>
      <section class="flex gap-x-6">
        <div class="flex-1 min-w-0">
          <GroupProjectsChart
            label="Ubuntu"
            :measure="metric"
            :projects="ubuntuConfigurations"
          />
        </div>
        <div class="flex-1 min-w-0">
          <GroupProjectsChart
            label="Windows"
            :measure="metric"
            :projects="windowsConfigurations"
          />
        </div>
      </section>
    </div>
  </DashboardPage>
  >
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import { jbrLinuxConfigurations, jbrMacConfigurations, jbrWindowsConfigurations } from "./configurations"

const metricsNames = [
  "CircleTests",
  "EllipseTests-fill-false",
  "EllipseTests-fill-true",
  "dc_boulder_2013-13-30-06-13-17",
  "dc_boulder_2013-13-30-06-13-20",
  "dc_shp_alllayers_2013-00-30-07-00-43",
  "dc_shp_alllayers_2013-00-30-07-00-47",
  "dc_spearfish_2013-11-30-06-11-15",
  "dc_spearfish_2013-11-30-06-11-19",
  "dc_topp:states_2013-11-30-06-11-06",
  "dc_topp:states_2013-11-30-06-11-07",
  "spiralTest-dash-false",
  "spiralTest-fill",
  "test_z_625k",
].map((metric) => metric + ".ser")
const ubuntuConfigurations = jbrLinuxConfigurations.map((config) => "Mapbench_" + config)
const macOSConfigurations = jbrMacConfigurations.map((config) => "Mapbench_" + config)
const windowsConfigurations = jbrWindowsConfigurations.map((config) => "Mapbench_" + config)
</script>
