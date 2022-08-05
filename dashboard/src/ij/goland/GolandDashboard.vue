<template>
  <Dashboard>
    <template #toolbar>
      <TimeRangeSelect :configurator="timeRangeConfigurator" />
      <DimensionSelect
        label="Branch"
        :dimension="branchConfigurator"
      />
      <DimensionHierarchicalSelect
        label="Machine"
        :dimension="machineConfigurator"
      />
      <DimensionSelect
        label="Triggered by"
        :dimension="triggeredByConfigurator"
      />
    </template>
    <GroupLineChart
      label="Indexing: Lightweight projects"
      measure="indexing"
      :projects="['flux/indexing', 'delve/indexing', 'istio/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Number Of Indexed Files: Lightweight projects"
      measure="numberOfIndexedFiles"
      :projects="['flux/indexing', 'delve/indexing', 'istio/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Indexing: Heavyweight projects"
      measure="indexing"
      :projects="['moby/indexing', 'mattermost-server/indexing', 'cockroach/indexing', 'kubernetes/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Number Of Indexed Files: Heavyweight projects"
      measure="numberOfIndexedFiles"
      :projects="['moby/indexing', 'mattermost-server/indexing', 'cockroach/indexing', 'kubernetes/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Inspection execution time: Lightweight projects"
      measure="globalInspections"
      :projects="['istio/inspection', 'moby/inspection', 'flux/inspection', 'delve/inspection']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Inspection execution time: Heavyweight projects"
      measure="globalInspections"
      :projects="['cockroach/inspection', 'kubernetes/inspection']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Local inspection execution time"
      measure="localInspections"
      :projects="['kubernetes/localInspection', 'mattermost-server/localInspection', 'GO-5422/localInspection']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Typing: average responsiveness time"
      measure="test#average_awt_delay"
      :projects="['mattermost-server/typing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Typing: total time"
      measure="typing"
      :projects="['mattermost-server/typing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Find Usages execution time"
      measure="findUsages"
      :projects="['vault/findUsage']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Find Usages number of found usages"
      measure="findUsages#number"
      :projects="['vault/findUsage']"
      :server-configurator="serverConfigurator"
    />
  </Dashboard>
</template>

<script lang="ts" setup>
import { initDataComponent } from "../../shared/DataQueryExecutor"
import { PersistentStateManager } from "../../shared/PersistentStateManager"
import { chartDefaultStyle } from "../../shared/chart"
import Dashboard from "../../shared/components/Dashboard.vue"
import DimensionHierarchicalSelect from "../../shared/components/DimensionHierarchicalSelect.vue"
import DimensionSelect from "../../shared/components/DimensionSelect.vue"
import GroupLineChart from "../../shared/components/GroupLineChart.vue"
import TimeRangeSelect from "../../shared/components/TimeRangeSelect.vue"
import { dimensionConfigurator } from "../../shared/configurators/DimensionConfigurator"
import { MachineConfigurator } from "../../shared/configurators/MachineConfigurator"
import { privateBuildConfigurator } from "../../shared/configurators/PrivateBuildConfigurator"
import { ServerConfigurator } from "../../shared/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "../../shared/configurators/TimeRangeConfigurator"
import { chartStyleKey } from "../../shared/injectionKeys"
import { provideReportUrlProvider } from "../../shared/lineChartTooltipLinkProvider"
import { provide } from "vue"
import { useRouter } from "vue-router"

provide(chartStyleKey, {
  ...chartDefaultStyle,
})

provideReportUrlProvider()

const persistentStateManager = new PersistentStateManager("goland_dashboard", {
  machine: "linux-blade-hetzner",
  project: [],
  branch: "master",
}, useRouter())

const serverConfigurator = new ServerConfigurator("perfint", "goland")
const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const branchConfigurator = dimensionConfigurator("branch", serverConfigurator, persistentStateManager, true, [timeRangeConfigurator], (a, _) => {
  return a.includes("/") ? 1 : -1
})
const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator, branchConfigurator])
const triggeredByConfigurator = privateBuildConfigurator(serverConfigurator, persistentStateManager, [branchConfigurator, timeRangeConfigurator])
const configurators = [
  serverConfigurator,
  branchConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  triggeredByConfigurator
]
initDataComponent(configurators)
</script>