<template>
  <StartupPage>
    <template #default="{ configurators }">
      <Divider label="Bootstrap" />
      <section class="grid grid-cols-2 gap-x-6">
        <LineChart
          :measures="['appInit_d', 'pluginDescriptorLoading_d', 'app initialization.end', 'connect FSRecords']"
          :configurators="configurators"
          title=""
          tooltip-trigger="axis"
        />
        <LineChart
          :measures="['bootstrap_d', 'appStarter_d', 'pluginDescriptorInitV18_d', 'euaShowing_d']"
          :configurators="configurators"
          title=""
          tooltip-trigger="axis"
        />
      </section>

      <section class="grid grid-cols-2 gap-x-6">
        <LineChart
          :measures="['PHM classes preloading', 'SvgCache creation', 'RunManager initialization']"
          :configurators="configurators"
          title=""
          tooltip-trigger="axis"
        />
        <LineChart
          :measures="['classLoadingTime', 'classLoadingSearchTime', 'classLoadingDefineTime']"
          :configurators="configurators"
          title=""
          tooltip-trigger="axis"
        />
      </section>

      <Divider label="Class and Resource Loading" />
      <section class="grid grid-cols-2 gap-x-6">
        <LineChart
          :measures="['classLoadingCount', 'resourceLoadingCount', 'classLoadingPreparedCount', 'classLoadingLoadedCount']"
          :configurators="configurators"
          title=""
          tooltip-trigger="axis"
        />
        <LineChart
          :measures="[
            'metrics.classLoadingMetrics/inlineCount',
            'metrics.classLoadingMetrics/companionCount',
            'metrics.classLoadingMetrics/lambdaCount',
            'metrics.classLoadingMetrics/methodHandleCount',
          ]"
          :configurators="configurators"
          title=""
          tooltip-trigger="axis"
        />
      </section>

      <Divider label="Services" />
      <section class="grid grid-cols-2 gap-x-6">
        <LineChart
          :skip-zero-values="false"
          :measures="['serviceSyncPreloading_d', 'serviceAsyncPreloading_d', 'projectServiceSyncPreloading_d', 'projectServiceAsyncPreloading_d']"
          :configurators="configurators"
          title=""
          tooltip-trigger="axis"
        />
        <LineChart
          :measures="['projectDumbAware', 'appComponentCreation_d', 'projectComponentCreation_d']"
          :configurators="configurators"
          title=""
          tooltip-trigger="axis"
        />
      </section>

      <Divider label="Post-opening" />
      <section class="grid grid-cols-2 gap-x-6">
        <LineChart
          :measures="['editorRestoring', 'editorRestoringTillPaint', 'file opening in EDT']"
          :configurators="configurators"
          title=""
          tooltip-trigger="axis"
        />
        <LineChart
          :measures="['splash_i', 'startUpCompleted', 'metrics.totalOpeningTime/timeFromAppStartTillAnalysisFinished']"
          :configurators="configurators"
          title=""
          tooltip-trigger="axis"
        />
      </section>

      <span v-if="highlightingPasses">
        <Divider label="Highlighting Passes" />
        <LineChart
          :measures="highlightingPasses"
          :configurators="configurators"
          title=""
        />
        <LineChart
          :measures="['metrics.codeAnalysisDaemon/fusExecutionTime', 'metrics.runDaemon/executionTime']"
          :configurators="configurators"
          title=""
          tooltip-trigger="axis"
        />
      </span>

      <Divider label="Notifications" />
      <LineChart
        :measures="['metrics.notifications/number']"
        :skip-zero-values="false"
        :configurators="configurators"
        title=""
        :with-measure-name="true"
        tooltip-trigger="axis"
      />

      <Divider label="Exit" />
      <LineChart
        :measures="['metrics.exitMetrics/application.exit', 'metrics.exitMetrics/saveSettingsOnExit', 'metrics.exitMetrics/disposeProjects']"
        :configurators="configurators"
        title=""
        tooltip-trigger="axis"
      />
    </template>
  </StartupPage>
</template>
<script setup lang="ts">
import LineChart from "../charts/LineChart.vue"
import Divider from "../common/Divider.vue"
import StartupPage from "./StartupPage.vue"
import { fetchHighlightingPasses } from "./utils"

const highlightingPasses = fetchHighlightingPasses()
</script>
