<template>
  <DashboardPage
    db-name="jbr"
    table="report"
    persistent-id="jbr_team_dashboard"
    :with-installer="false"
    :is-build-number-exists="true"
  >
    <div
      v-for="mode in renderingModes"
      :key="mode"
    >
      <div
        v-for="metric in metricsNames"
        :key="metric"
      >
        <div class="relative flex py-5 items-center">
          <div class="grow border-t" />
          <span class="shrink mx-4 text-lg">{{ metric }}_{{ mode }}</span>
          <div class="grow border-t" />
        </div>
        <section>
          <GroupProjectsChart
            label="macOS"
            :measure="metric + '_' + mode"
            :projects="macOSConfigurations"
          />
        </section>
        <section class="flex gap-x-6">
          <div class="flex-1 min-w-0">
            <GroupProjectsChart
              label="Ubuntu"
              :measure="metric + '_' + mode"
              :projects="ubuntuConfigurations"
            />
          </div>
          <div class="flex-1 min-w-0">
            <GroupProjectsChart
              label="Windows"
              :measure="metric + '_' + mode"
              :projects="windowsConfigurations"
            />
          </div>
        </section>
      </div>
    </div>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import { jbrLinuxConfigurations, jbrMacConfigurations, jbrWindowsConfigurations } from "./configurations"

const renderingModes = ["buffer", "onscreen", "volatile"]

const metricsNames = [
  "ClipFlatOval",
  "FlatBox",
  "Image",
  "VolImage",
  "LinGrad3RotatedOval",
  "RadGrad3RotatedOvalAA",
  "WiredBubbles",
  "FlatQuadAA",
  "WiredQuadAA",
  "TextNoAA",
  "TextLCD",
  "TextGray",
  "LargeTextNoAA",
  "FlatOval_XOR",
  "TextWiredQuadMix",
  "VolImageFlatBoxMix",
  "VolImageTextNoAABat",
]

const ubuntuConfigurations = jbrLinuxConfigurations.map((config) => "Render_" + config)
const macOSConfigurations = jbrMacConfigurations.map((config) => "Render_" + config)
const windowsConfigurations = jbrWindowsConfigurations.map((config) => "Render_" + config)
</script>
