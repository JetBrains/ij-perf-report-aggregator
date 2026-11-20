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
      <JbrPlatformCharts
        v-for="metric in metricsNames"
        :key="metric"
        :measure="metric + '_' + mode"
        benchmark="Render"
      />
    </div>
  </DashboardPage>
</template>

<script setup lang="ts">
import DashboardPage from "../common/DashboardPage.vue"
import JbrPlatformCharts from "./JbrPlatformCharts.vue"

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
</script>
