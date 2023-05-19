<template>
  <DashboardPage
    v-slot="{serverConfigurator, dashboardConfigurators}"
    db-name="jbr"
    table="report"
    persistent-id="jbr_dacapo_dashboard"
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
          :server-configurator="serverConfigurator"
          :configurators="dashboardConfigurators"
        />
      </section>
      <section class="flex gap-x-6">
        <div class="flex-1 min-w-0">
          <GroupProjectsChart
            label="Ubuntu"
            :measure="metric"
            :projects="ubuntuConfigurations"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </div>
        <div class="flex-1 min-w-0">
          <GroupProjectsChart
            label="Windows"
            :measure="metric"
            :projects="windowsConfigurations"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </div>
      </section>
    </div>
  </DashboardPage>>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const metricsNames = ["avrora", "fop", "h2", "jython", "luindex", "lusearch", "pmd", "sunflow", "xalan"]
const ubuntuConfigurations = ["Ubuntu2004x64", "Ubuntu2004x64OGL", "Ubuntu2204x64", "Ubuntu2204x64OGL"].map(config => "DaCapo_" + config)
const macOSConfigurations = ["macOS13x64OGL", "macOS13x64Metal", "macOS13aarch64OGL", "macOS13aarch64Metal", "macOS12x64OGL", "macOS12x64Metal", "macOS12aarch64OGL",
  "macOS12aarch64Metal"].map(config => "DaCapo_" + config)
const windowsConfigurations = ["Windows10x64"].map(config => "DaCapo_" + config)


</script>