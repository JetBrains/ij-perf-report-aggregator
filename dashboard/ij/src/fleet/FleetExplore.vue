<template>
  <Dashboard>
    <template #toolbar>
      <MeasureSelect :configurator="measureConfigurator" />
      <DimensionHierarchicalSelect
        label="Machine"
        :dimension="machineConfigurator"
      />
      <TimeRangeSelect :configurator="timeRangeConfigurator" />
    </template>

    <div class="grid grid-cols-2 gap-4">
      <LineChartCard />
    </div>
  </Dashboard>
</template>

<script setup lang="ts">
import { initDataComponent } from "shared/src/DataQueryExecutor"
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import Dashboard from "shared/src/components/Dashboard.vue"
import DimensionHierarchicalSelect from "shared/src/components/DimensionHierarchicalSelect.vue"
import LineChartCard from "shared/src/components/LineChartCard.vue"
import MeasureSelect from "shared/src/components/MeasureSelect.vue"
import TimeRangeSelect from "shared/src/components/TimeRangeSelect.vue"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { MeasureConfigurator } from "shared/src/configurators/MeasureConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"

provideReportUrlProvider()

const persistentStateManager = new PersistentStateManager("fleet-explore")
const serverConfigurator = new ServerConfigurator("fleet")
const measureConfigurator = new MeasureConfigurator(serverConfigurator, persistentStateManager)
const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager)
const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)

initDataComponent([
  serverConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  measureConfigurator,
])
</script>
