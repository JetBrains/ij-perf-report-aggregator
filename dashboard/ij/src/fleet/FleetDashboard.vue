<template>
  <Dashboard>
    <template #toolbar>
      <DimensionHierarchicalSelect
        label="Machine"
        :dimension="machineConfigurator"
      />
      <TimeRangeSelect :configurator="timeRangeConfigurator" />
    </template>

    <div class="relative flex py-5 items-center">
      <div class="flex-grow border-t border-gray-400" />
      <span class="flex-shrink mx-4 text-gray-400 text-lg">Remote Mode</span>
      <div class="flex-grow border-t border-gray-400" />
    </div>
    <!-- :skip-zero-values="false" because computed measures cannot be filtered -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 pt-4">
      <LineChartCard
        label="editor appeared"
        :measures='["editor appeared.end"]'
        :skip-zero-values="false"
      />
      <LineChartCard
        label="time to edit"
        :measures='["time to edit.end"]'
        :skip-zero-values="false"
      />
      <LineChartCard
        label="terminal ready"
        :measures='["terminal ready.end"]'
        :skip-zero-values="false"
      />
      <LineChartCard
        label="file tree rendered"
        :measures='["file tree rendered.end"]'
        :skip-zero-values="false"
      />
      <LineChartCard
        label="highlighting done"
        :measures='["highlighting done.end"]'
        :skip-zero-values="false"
      />
      <LineChartCard
        label="window appeared"
        :measures='["window appeared.end"]'
        :skip-zero-values="false"
      />
    </div>

    <div class="relative flex py-5 items-center">
      <div class="flex-grow border-t border-gray-400" />
      <span class="flex-shrink mx-4 text-gray-400 text-lg">ShortCircuit</span>
      <div class="flex-grow border-t border-gray-400" />
    </div>
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 pt-4">
      <LineChartCard
        label="editor appeared"
        :measures='["shortCircuit.editor appeared.end"]'
        :skip-zero-values="false"
      />
      <LineChartCard
        label="time to edit"
        :measures='["shortCircuit.time to edit.end"]'
        :skip-zero-values="false"
      />
      <LineChartCard
        label="terminal ready"
        :measures='["shortCircuit.terminal ready.end"]'
        :skip-zero-values="false"
      />
      <LineChartCard
        label="file tree rendered"
        :measures='["shortCircuit.file tree rendered.end"]'
        :skip-zero-values="false"
      />
      <LineChartCard
        label="highlighting done"
        :measures='["shortCircuit.highlighting done.end"]'
        :skip-zero-values="false"
      />
      <LineChartCard
        label="window appeared"
        :measures='["shortCircuit.window appeared.end"]'
        :skip-zero-values="false"
      />
    </div>
  </Dashboard>
</template>

<script setup lang="ts">
import { initDataComponent } from "shared/src/DataQueryExecutor"
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import Dashboard from "shared/src/components/Dashboard.vue"
import DimensionHierarchicalSelect from "shared/src/components/DimensionHierarchicalSelect.vue"
import LineChartCard from "shared/src/components/LineChartCard.vue"
import TimeRangeSelect from "shared/src/components/TimeRangeSelect.vue"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
import { useRouter } from "vue-router"

const persistentStateManager = new PersistentStateManager("fleet-dashboard", {
  machine: "macMini 2018",
}, useRouter())
const serverConfigurator = new ServerConfigurator("fleet")
provideReportUrlProvider()

const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator])

initDataComponent([
  serverConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
])
</script>
