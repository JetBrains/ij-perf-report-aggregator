<template>
  <StartupPage>
    <Divider label="Bootstrap" />
    <section class="grid grid-cols-2 gap-x-6">
      <LineChartCard :measures="['appInit_d', 'pluginDescriptorLoading_d', 'app initialization.end', 'connect FSRecords']" />
      <LineChartCard :measures="['bootstrap_d', 'appStarter_d', 'pluginDescriptorInitV18_d', 'euaShowing_d']" />
    </section>

    <section class="grid grid-cols-2 gap-x-6">
      <LineChartCard :measures="['PHM classes preloading', 'SvgCache creation', 'RunManager initialization']" />
      <LineChartCard :measures="['classLoadingTime', 'classLoadingSearchTime', 'classLoadingDefineTime']" />
    </section>

    <Divider label="Class and Resource Loading" />
    <section class="grid grid-cols-2 gap-x-6">
      <LineChartCard :measures="['classLoadingCount', 'resourceLoadingCount', 'classLoadingPreparedCount', 'classLoadingLoadedCount']" />
      <LineChartCard
        :measures="[
          'metrics.classLoadingMetrics/inlineCount',
          'metrics.classLoadingMetrics/companionCount',
          'metrics.classLoadingMetrics/lambdaCount',
          'metrics.classLoadingMetrics/methodHandleCount',
        ]"
      />
    </section>

    <Divider label="Services" />
    <section class="grid grid-cols-2 gap-x-6">
      <LineChartCard
        :skip-zero-values="false"
        :measures="['serviceSyncPreloading_d', 'serviceAsyncPreloading_d', 'projectServiceSyncPreloading_d', 'projectServiceAsyncPreloading_d']"
      />
      <LineChartCard :measures="['projectDumbAware', 'appComponentCreation_d', 'projectComponentCreation_d']" />
    </section>

    <Divider label="Post-opening" />
    <section class="grid grid-cols-2 gap-x-6">
      <LineChartCard :measures="['editorRestoring', 'editorRestoringTillPaint', 'file opening in EDT']" />
      <LineChartCard :measures="['splash_i', 'startUpCompleted', 'metrics.totalOpeningTime/timeFromAppStartTillAnalysisFinished']" />
    </section>

    <span v-if="highlightingPasses">
      <Divider label="Highlighting Passes" />
      <LineChartCard :measures="highlightingPasses" />
      <LineChartCard :measures="['metrics.codeAnalysisDaemon/fusExecutionTime', 'metrics.runDaemon/executionTime']" />
    </span>

    <Divider label="Notifications" />
    <LineChartCard
      :measures="['metrics.notifications/number']"
      :skip-zero-values="false"
    />

    <Divider label="Exit" />
    <LineChartCard :measures="['metrics.exitMetrics/application.exit', 'metrics.exitMetrics/saveSettingsOnExit', 'metrics.exitMetrics/disposeProjects']" />

    <slot />
  </StartupPage>
</template>
<script setup lang="ts">
import LineChartCard from "../charts/LineChartCard.vue"
import Divider from "../common/Divider.vue"
import StartupPage from "./StartupPage.vue"
import { fetchHighlightingPasses } from "./utils"

const highlightingPasses = fetchHighlightingPasses()
</script>
