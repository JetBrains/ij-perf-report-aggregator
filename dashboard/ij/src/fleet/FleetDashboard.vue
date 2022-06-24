<template>
  <Dashboard>
    <template #toolbar>
      <DimensionHierarchicalSelect
        label="Machine"
        :dimension="machineConfigurator"
      />
      <TimeRangeSelect :configurator="timeRangeConfigurator" />
    </template>

    <!-- :skip-zero-values="false" because computed measures cannot be filtered -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 pt-4">
      <LineChartCard
        :measures='["editor appeared.end"]'
        :skip-zero-values="false"
      />
      <LineChartCard
        :measures='["window appeared.end"]'
        :skip-zero-values="false"
      />
    </div>
  </Dashboard>
</template>

<script setup lang="ts">
import { initDataComponent } from "shared/src/DataQueryExecutor"
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import Dashboard from "shared/src/components/Dashboard.vue"
import LineChartCard from "shared/src/components/LineChartCard.vue"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
import { useRouter } from "vue-router"
import DimensionHierarchicalSelect from "../../../shared/src/components/DimensionHierarchicalSelect.vue"
import TimeRangeSelect from "../../../shared/src/components/TimeRangeSelect.vue"

const persistentStateManager = new PersistentStateManager("fleet-dashboard", {
  machine: "macMini 2018",
}, useRouter())
const serverConfigurator = new ServerConfigurator("fleet")
provideReportUrlProvider()

const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager)
const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)

initDataComponent([
  serverConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
])
</script>
