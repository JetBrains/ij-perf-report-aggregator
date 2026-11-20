<template>
  <DashboardPage
    db-name="jbr"
    table="report"
    persistent-id="jbr_team_dashboard"
    :with-installer="false"
    :is-build-number-exists="true"
  >
    <JbrPlatformCharts
      v-for="metric in metricsNames"
      :key="metric"
      :measure="metric"
      benchmark="Mapbench"
    />
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../common/DashboardPage.vue"
import JbrPlatformCharts from "./JbrPlatformCharts.vue"

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
</script>
