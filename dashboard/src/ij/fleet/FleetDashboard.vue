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
import { initDataComponent } from "../../shared/DataQueryExecutor"
import { PersistentStateManager } from "../../shared/PersistentStateManager"
import Dashboard from "../../shared/components/Dashboard.vue"
import DimensionHierarchicalSelect from "../../shared/components/DimensionHierarchicalSelect.vue"
import LineChartCard from "../../shared/components/LineChartCard.vue"
import TimeRangeSelect from "../../shared/components/TimeRangeSelect.vue"
import { MachineConfigurator } from "../../shared/configurators/MachineConfigurator"
import { ServerConfigurator } from "../../shared/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "../../shared/configurators/TimeRangeConfigurator"
import { provideReportUrlProvider } from "../../shared/lineChartTooltipLinkProvider"
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
