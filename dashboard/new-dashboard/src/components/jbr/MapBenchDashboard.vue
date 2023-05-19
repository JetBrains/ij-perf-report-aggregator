<template>
  <DashboardPage
    db-name="jbr"
    table="report"
    persistent-id="jbr_mapbench_dashboard"
    :with-installer="false"
    :is-build-number-exists="true"
  >
    <div
      v-for="metric in metricsNames"
      :key="metric"
    >
      <div class="relative flex py-5 items-center">
        <div class="flex-grow border-t border-gray-400" />
        <span class="flex-shrink mx-4 text-gray-400 text-lg">{{ metric }}</span>
        <div class="flex-grow border-t border-gray-400" />
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
  </DashboardPage>>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const metricsNames = ["CircleTests", "EllipseTests-fill-false", "EllipseTests-fill-true", "spiralTest-dash-false", "spiralTest-fill", "dc_boulder_2013-13-30-06-13-17",
  "dc_boulder_2013-13-30-06-13-20", "dc_shp_alllayers_2013-00-30-07-00-43", "dc_shp_alllayers_2013-00-30-07-00-47", "dc_spearfish_2013-11-30-06-11-15",
  "dc_spearfish_2013-11-30-06-11-19", "test_z_625k"].flatMap(test => {
  return ["ser.Pct95"].map(stat => test + "." + stat)
})
const ubuntuConfigurations = ["Ubuntu2004x64", "Ubuntu2004x64OGL", "Ubuntu2204x64", "Ubuntu2204x64OGL"].map(config => "Mapbench_" + config)
const macOSConfigurations = ["macOS13x64OGL", "macOS13x64Metal", "macOS13aarch64OGL", "macOS13aarch64Metal", "macOS12x64OGL", "macOS12x64Metal", "macOS12aarch64OGL",
  "macOS12aarch64Metal"].map(config => "Mapbench_" + config)
const windowsConfigurations = ["Windows10x64"].map(config => "Mapbench_" + config)

</script>