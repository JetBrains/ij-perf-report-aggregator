<template>
  <StartupPage>
    <Divider label="Bootstrap" />
    <section class="grid grid-cols-2 gap-x-6">
      <LineChartCard
        :measures='["appInitPreparation_d", "appInit_d", "pluginDescriptorLoading_d", "app initialization.end"]'
      />
      <LineChartCard
        :measures='["bootstrap_d", "appStarter_d", "pluginDescriptorInitV18_d", "RunManager initialization", "euaShowing_d"]'
      />
    </section>

    <Divider label="Class and Resource Loading" />
    <LineChartCard
      :measures='["classLoadingTime", "classLoadingSearchTime", "classLoadingDefineTime"]'
    />
    <section class="grid grid-cols-2 gap-x-6">
      <LineChartCard
        :measures='["classLoadingCount", "resourceLoadingCount", "classLoadingPreparedCount", "classLoadingLoadedCount"]'
      />
      <LineChartCard
        :measures='["metrics.classLoadingMetrics/inlineCount", "metrics.classLoadingMetrics/companionCount",
                    "metrics.classLoadingMetrics/lambdaCount", "metrics.classLoadingMetrics/methodHandleCount"]'
      />
    </section>

    <Divider label="Services" />
    <section class="grid grid-cols-2 gap-x-6">
      <LineChartCard
        :skip-zero-values="false"
        :measures='["serviceSyncPreloading_d", "serviceAsyncPreloading_d", "projectServiceSyncPreloading_d", "projectServiceAsyncPreloading_d"]'
      />
      <LineChartCard
        :measures='["projectDumbAware", "appComponentCreation_d", "projectComponentCreation_d"]'
      />
    </section>

    <Divider label="Post-opening" />
    <section class="grid grid-cols-2 gap-x-6">
      <LineChartCard
        :measures='["editorRestoring", "editorRestoringTillPaint", "file opening in EDT"]'
      />
      <LineChartCard
        :measures='["splash_i", "startUpCompleted", "metrics.totalOpeningTime/timeFromAppStartTillAnalysisFinished"]'
      />
    </section>

    <Divider label="Highlighting Passes" />
    <LineChartCard
      :measures="highlightingPasses"
    />

    <Divider label="Exit" />
    <LineChartCard
      :measures='["metrics.exitMetrics/application.exit", "metrics.exitMetrics/saveSettingsOnExit", "metrics.exitMetrics/disposeProjects"]'
    />

    <slot />
  </StartupPage>
</template>
<script setup lang="ts">
import { ref } from "vue"
import LineChartCard from "../charts/LineChartCard.vue"
import Divider from "../common/Divider.vue"
import StartupPage from "./StartupPage.vue"
import { fetchHighlightingPasses } from "./utils"

const highlightingPasses = ref<string[]>()
fetchHighlightingPasses(highlightingPasses)
</script>