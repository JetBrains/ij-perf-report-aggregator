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
        <div class="flex-grow border-t" />
        <span class="flex-shrink mx-4 text-lg">{{ metric }}</span>
        <div class="flex-grow border-t" />
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
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import { jbrLinuxConfigurations, jbrMacConfigurations, jbrWindowsConfigurations } from "./configurations"

const metricsNames = [
  "ArgbSurfaceBlitImage",
  "LinGrad3RotatedOvalAA",
  "LinGradRotatedOval",
  "LinGradRotatedOvalAA",
  "ArgbSwBlitImage",
  "BgrSurfaceBlitImage",
  "LinGrad3RotatedOval",
  "RadGrad3RotatedOval",
  "RadGrad3RotatedOvalAA",
  "FlatBox",
  "RotatedOval",
  "WiredBubbles",
  "ClipFlatOvalAA",
  "ClipFlatBoxAA",
  "FlatBoxAA",
  "FlatOvalAA",
  "ClipFlatOval",
  "RotatedOvalAA",
  "VolImageAA",
  "ImageAA",
  "RotatedBox",
  "RotatedBoxAA",
  "WiredBox",
  "FlatOval",
  "WiredBoxAA",
  "WiredBubblesAA",
  "Lines",
  "Image",
  "ClipFlatBox",
  "VolImage",
  "LargeTextGray",
  "LargeTextNoAA",
  "Image_XOR",
  "WhiteTextGray",
  "LargeTextLCD",
  "BgrSwBlitImage",
  "FlatQuad",
  "TextLCD",
  "WhiteTextLCD",
  "TextNoAA",
  "FlatOval_XOR",
  "TextGray",
  "WhiteTextNoAA",
  "LinesAA",
  "WiredQuadAA",
  "Lines_XOR",
  "RotatedBox_XOR",
  "WiredQuad",
  "TextNoAA_XOR",
  "FlatQuadAA",
  "TextLCD_XOR",
]
const ubuntuConfigurations = jbrLinuxConfigurations.map((config) => "Render_" + config)
const macOSConfigurations = jbrMacConfigurations.map((config) => "Render_" + config)
const windowsConfigurations = jbrWindowsConfigurations.map((config) => "Render_" + config)
</script>
