<template>
  <StartupPage>
    <Divider label="Bootstrap" />
    <section class="grid grid-cols-2 gap-x-6">
      <LineChartCard
        :measures='["bootstrap_d", "appInitPreparation_d", "appInit_d", "pluginDescriptorLoading_d", "euaShowing_d", "appStarter_d"]'
      />
      <LineChartCard
        :measures='["pluginDescriptorInitV18_d", "RunManager initialization"]'
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

    <Divider label="Exit" />
    <LineChartCard
      :measures='["metrics.exitMetrics/application.exit", "metrics.exitMetrics/saveSettingsOnExit", "metrics.exitMetrics/disposeProjects"]'
    />
    <Divider label="Memory" />
    <LineChartCard
      :skip-zero-values="false"
      :measures='["metrics.memory/usedMb", "metrics.memory/metaspaceMb","metrics.memory/maxMb"]'
    />
    <Divider label="GC" />
    <div class="grid grid-cols-2 gap-6">
      <LineChartCard
        class="col-span-2"
        :skip-zero-values="false"
        :measures='["metrics.gc/totalHeapUsedMax", "metrics.gc/freedMemoryByGC"]'
      />
      <LineChartCard
        :skip-zero-values="false"
        label="Number of pauses"
        :measures='["metrics.gc/gcPauseCount"]'
      />
      <LineChartCard
        :skip-zero-values="false"
        :measures='["metrics.gc/fullGCPause","metrics.gc/gcPause"]'
      />
    </div>
  </StartupPage>
</template>
<script setup lang="ts">
import LineChartCard from "shared/src/components/LineChartCard.vue"
import Divider from "../common/Divider.vue"
import StartupPage from "./StartupPage.vue"
</script>