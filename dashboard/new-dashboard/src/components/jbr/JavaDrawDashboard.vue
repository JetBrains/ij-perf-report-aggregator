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
        <div class="flex-grow border-t border-gray-400" />
        <span class="flex-shrink mx-4 text-lg">{{ metric }}</span>
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
  </DashboardPage>
  >
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import { jbrLinuxConfigurations, jbrMacConfigurations, jbrWindowsConfigurations } from "./configurations"

const metricsNames = ["Plus_200_Random_Small_Circles", "Plus_2_SweepGradient_Circles", "Plus_320_Long_Lines", "Plus_4000_Random_Small_Circles"]
const ubuntuConfigurations = jbrLinuxConfigurations.map((config) => "JavaDraw_" + config)
const macOSConfigurations = jbrMacConfigurations.map((config) => "JavaDraw_" + config)
const windowsConfigurations = jbrWindowsConfigurations.map((config) => "JavaDraw_" + config)
</script>
