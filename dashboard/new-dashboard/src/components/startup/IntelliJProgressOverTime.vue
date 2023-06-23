<template>
  <StartupPage>
    <template #toolbar>
      <AggregationOperatorSelect />
    </template>
    <div class="grid grid-cols-2 gap-6">
      <BarChartCard :measures="['bootstrap_d', 'appInitPreparation_d', 'appInit_d', 'pluginDescriptorLoading_d']" />
      <BarChartCard :measures="['splash_i', 'startUpCompleted']" />

      <!-- todo "moduleLoading_d" when it will be fixed -->
      <BarChartCard :measures="['appStarter_d', 'serviceSyncPreloading_d', 'serviceAsyncPreloading_d', 'projectServiceSyncPreloading_d', 'projectServiceAsyncPreloading_d']" />
      <BarChartCard :measures="['projectDumbAware', 'editorRestoring', 'editorRestoringTillPaint']" />
    </div>
  </StartupPage>
</template>
<script setup lang="ts">
import { provide } from "vue"
import { useRouter } from "vue-router"
import { AggregationOperatorConfigurator } from "../../configurators/AggregationOperatorConfigurator"
import { aggregationOperatorConfiguratorKey } from "../../shared/injectionKeys"
import AggregationOperatorSelect from "../charts/AggregationOperatorSelect.vue"
import BarChartCard from "../charts/BarChartCard.vue"
import { PersistentStateManager } from "../common/PersistentStateManager"
import StartupPage from "./StartupPage.vue"

const persistentStateManager = new PersistentStateManager(
  "ij-progress-over-time",
  {
    product: "IU",
    project: "simple for IJ",
    machine: "macMini M1, 16GB",
    branch: "master",
  },
  useRouter()
)
provide(aggregationOperatorConfiguratorKey, new AggregationOperatorConfigurator(persistentStateManager))
</script>
