<template>
  <Dashboard>
    <template #toolbar>
      <TimeRangeSelect :configurator="timeRangeConfigurator" />
      <DimensionHierarchicalSelect
        label="Machine"
        :dimension="machineConfigurator"
      />
      <DimensionSelect
        label="Triggered by"
        :dimension="triggeredByConfigurator"
      />
    </template>
    <div class="relative flex py-5 items-center">
      <div class="flex-grow border-t border-gray-400" />
      <span class="flex-shrink mx-4 text-gray-400 text-lg">Core</span>
      <div class="flex-grow border-t border-gray-400" />
    </div>
    <GroupLineChart
      label="Typing (time)"
      measure="fleet.test"
      :projects="['multiCaretTyping', 'stressEnter', 'stressTyping']"
      value-unit="ns"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Typing (average delay)"
      measure="awt.delay"
      value-unit="ns"
      :projects="['multiCaretTyping', 'stressEnter', 'stressTyping']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Typing (max delay)"
      measure="max.awt.delay"
      value-unit="ns"
      :projects="['multiCaretTyping', 'stressEnter', 'stressTyping']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Highlighting"
      value-unit="ns"
      measure="fleet.test"
      :projects="['stressHighlighting']"
      :server-configurator="serverConfigurator"
    />
    <div class="relative flex py-5 items-center">
      <div class="flex-grow border-t border-gray-400" />
      <span class="flex-shrink mx-4 text-gray-400 text-lg">PHP</span>
      <div class="flex-grow border-t border-gray-400" />
    </div>
    <GroupLineChart
      label="Typing (time)"
      value-unit="ns"
      measure="fleet.test"
      :projects="['Typing in mPDF', 'Typing in mPDF With Backend', 'Pressing Enter in mPDF']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Typing (average delay)"
      value-unit="ns"
      measure="awt.delay"
      :projects="['Typing in mPDF', 'Typing in mPDF With Backend', 'Pressing Enter in mPDF']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Typing (max delay)"
      value-unit="ns"
      measure="max.awt.delay"
      :projects="['Typing in mPDF', 'Typing in mPDF With Backend', 'Pressing Enter in mPDF']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Other"
      value-unit="ns"
      measure="fleet.test"
      :projects="['Open mPDF', 'Frontend Completion in mPDF']"
      :server-configurator="serverConfigurator"
    />
  </Dashboard>
</template>

<script lang="ts" setup>
import { initDataComponent } from "shared/src/DataQueryExecutor"
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import { chartDefaultStyle } from "shared/src/chart"
import Dashboard from "shared/src/components/Dashboard.vue"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { privateBuildConfigurator } from "shared/src/configurators/PrivateBuildConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { chartStyleKey } from "shared/src/injectionKeys"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
import { provide } from "vue"
import { useRouter } from "vue-router"
import DimensionHierarchicalSelect from "../../../shared/src/components/DimensionHierarchicalSelect.vue"
import DimensionSelect from "../../../shared/src/components/DimensionSelect.vue"
import GroupLineChart from "../../../shared/src/components/GroupLineChart.vue"
import TimeRangeSelect from "../../../shared/src/components/TimeRangeSelect.vue"

provide(chartStyleKey, {
  ...chartDefaultStyle,
})

provideReportUrlProvider(false)

const persistentStateManager = new PersistentStateManager("fleet_perf_dashboard", {
  machine: "linux-blade-hetzner",
  project: [],
  branch: "master",
}, useRouter())

const serverConfigurator = new ServerConfigurator("fleet", "measure")
const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager, [])
const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const triggeredByConfigurator = privateBuildConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator])
const configurators = [
  serverConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  triggeredByConfigurator
]
initDataComponent(configurators)
</script>