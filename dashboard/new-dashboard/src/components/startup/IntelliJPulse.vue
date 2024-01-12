<template>
  <StartupPage>
    <template #default="{ configurators }">
      <Divider label="Bootstrap" />
      <section class="grid grid-cols-2 gap-x-6">
        <PerformanceLineChart
          :measures="['appInit_d', 'pluginDescriptorLoading_d', 'app initialization.end', 'connect FSRecords']"
          :configurators="configurators"
          title=""
        />
        <PerformanceLineChart
          :measures="['bootstrap_d', 'appStarter_d', 'pluginDescriptorInitV18_d', 'euaShowing_d']"
          :configurators="configurators"
          title=""
        />
      </section>

      <section class="grid grid-cols-2 gap-x-6">
        <PerformanceLineChart
          :measures="['PHM classes preloading', 'SvgCache creation', 'RunManager initialization']"
          :configurators="configurators"
          title=""
        />
        <PerformanceLineChart
          :measures="['classLoadingTime', 'classLoadingSearchTime', 'classLoadingDefineTime']"
          :configurators="configurators"
          title=""
        />
      </section>

      <Divider label="Class and Resource Loading" />
      <section class="grid grid-cols-2 gap-x-6">
        <PerformanceLineChart
          :measures="['classLoadingCount', 'resourceLoadingCount', 'classLoadingPreparedCount', 'classLoadingLoadedCount']"
          :configurators="configurators"
          title=""
        />
        <PerformanceLineChart
          :measures="[
            'metrics.classLoadingMetrics/inlineCount',
            'metrics.classLoadingMetrics/companionCount',
            'metrics.classLoadingMetrics/lambdaCount',
            'metrics.classLoadingMetrics/methodHandleCount',
          ]"
          :configurators="configurators"
          title=""
        />
      </section>

      <Divider label="Services" />
      <section class="grid grid-cols-2 gap-x-6">
        <PerformanceLineChart
          :skip-zero-values="false"
          :measures="['serviceSyncPreloading_d', 'serviceAsyncPreloading_d', 'projectServiceSyncPreloading_d', 'projectServiceAsyncPreloading_d']"
          :configurators="configurators"
          title=""
        />
        <PerformanceLineChart
          :measures="['projectDumbAware', 'appComponentCreation_d', 'projectComponentCreation_d']"
          :configurators="configurators"
          title=""
        />
      </section>

      <Divider label="Post-opening" />
      <section class="grid grid-cols-2 gap-x-6">
        <PerformanceLineChart
          :measures="['editorRestoring', 'editorRestoringTillPaint', 'file opening in EDT']"
          :configurators="configurators"
          title=""
        />
        <PerformanceLineChart
          :measures="['splash_i', 'startUpCompleted', 'metrics.totalOpeningTime/timeFromAppStartTillAnalysisFinished']"
          :configurators="configurators"
          title=""
        />
      </section>

      <span v-if="highlightingPasses">
        <Divider label="Highlighting Passes" />
        <PerformanceLineChart
          :measures="highlightingPasses"
          :configurators="configurators"
          title=""
        />
        <PerformanceLineChart
          :measures="['metrics.codeAnalysisDaemon/fusExecutionTime', 'metrics.runDaemon/executionTime']"
          :configurators="configurators"
          title=""
        />
      </span>

      <Divider label="Notifications" />
      <PerformanceLineChart
        :measures="['metrics.notifications/number']"
        :skip-zero-values="false"
        :configurators="configurators"
        title=""
        :with-measure-name="true"
      />

      <Divider label="Exit" />
      <PerformanceLineChart
        :measures="['metrics.exitMetrics/application.exit', 'metrics.exitMetrics/saveSettingsOnExit', 'metrics.exitMetrics/disposeProjects']"
        :configurators="configurators"
        title=""
      />
    </template>
  </StartupPage>
</template>
<script setup lang="ts">
import PerformanceLineChart from "../charts/PerformanceLineChart.vue"
import Divider from "../common/Divider.vue"
import StartupPage from "./StartupPage.vue"
import { fetchHighlightingPasses } from "./utils"

const highlightingPasses = fetchHighlightingPasses()
</script>
