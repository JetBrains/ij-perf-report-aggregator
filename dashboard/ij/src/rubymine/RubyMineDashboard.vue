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
      label="Indexing"
      measure="indexing"
      :projects="['diaspora-project-test/indexing', 'gem-rbs-collection-indexing-test/indexing', 'gitlab-project-test/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Number Of Indexed Files"
      measure="numberOfIndexedFiles"
      :projects="['diaspora-project-test/indexing', 'gem-rbs-collection-indexing-test/indexing', 'gitlab-project-test/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Scanning"
      measure="scanning"
      :projects="['diaspora-project-test/indexing', 'gem-rbs-collection-indexing-test/indexing', 'gitlab-project-test/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Inspection"
      measure="globalInspections"
      :projects="['diaspora-project-inspections-test/inspection-RubyResolve-app', 'diaspora-project-inspections-test/inspection-app', 'gitlab-project-inspections-test/inspection-RubyResolve-app', 'gitlab-project-inspections-test/inspection-app']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Find Usages: execution time"
      measure="findUsages"
      :projects="['RUBY-23764-Case1/ruby-23764-findusages-case1', 'RUBY-23764-Case2/ruby-23764-findusages-case2']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Find Usages: number of found usages"
      measure="findUsages#number"
      :projects="['RUBY-23764-Case1/ruby-23764-findusages-case1', 'RUBY-23764-Case2/ruby-23764-findusages-case2']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Completion Diaspora"
      measure="completion"
      :projects="['diaspora-project-test/completion/basic_completion', 'diaspora-project-test/completion/exceptions', 'diaspora-project-test/completion/localization']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Completion Gitlab"
      measure="completion"
      :projects="['gitlab-project-test/completion/basic_completion', 'gitlab-project-test/completion/exceptions', 'gitlab-project-test/completion/localization']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Typing: average delay"
      measure="test#average_awt_delay"
      :projects="['RUBY-26170/typing', 'RUBY-29334/typing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Typing: total time"
      measure="typing"
      :projects="['RUBY-26170/typing', 'RUBY-29334/typing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Get Symbol Members: execution time"
      measure="getSymbolMembers"
      :projects="['diaspora-project-test/getSymbolMembers-ApplicationController-false', 'diaspora-project-test/getSymbolMembers-ApplicationController-true', 'gitlab-project-test/getSymbolMembers-ApplicationController-false', 'gitlab-project-test/getSymbolMembers-ApplicationController-true']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Get Symbol Members: number"
      measure="getSymbolMembers#number"
      :projects="['diaspora-project-test/getSymbolMembers-ApplicationController-false', 'diaspora-project-test/getSymbolMembers-ApplicationController-true', 'gitlab-project-test/getSymbolMembers-ApplicationController-false', 'gitlab-project-test/getSymbolMembers-ApplicationController-true']"
      :server-configurator="serverConfigurator"
    />
  </Dashboard>
</template>

<script lang="ts" setup>
import { initDataComponent } from "shared/src/DataQueryExecutor"
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import { chartDefaultStyle } from "shared/src/chart"
import Dashboard from "shared/src/components/Dashboard.vue"
import DimensionHierarchicalSelect from "shared/src/components/DimensionHierarchicalSelect.vue"
import DimensionSelect from "shared/src/components/DimensionSelect.vue"
import GroupLineChart from "shared/src/components/GroupLineChart.vue"
import TimeRangeSelect from "shared/src/components/TimeRangeSelect.vue"
import { dimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { privateBuildConfigurator } from "shared/src/configurators/PrivateBuildConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { chartStyleKey } from "shared/src/injectionKeys"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
import { provide } from "vue"
import { useRouter } from "vue-router"

provide(chartStyleKey, {
  ...chartDefaultStyle,
})

provideReportUrlProvider()

const persistentStateManager = new PersistentStateManager("rubymine_dashboard", {
  machine: "Linux Munich i7-3770, 32 Gb",
  project: [],
  branch: "master",
}, useRouter())

const serverConfigurator = new ServerConfigurator("perfint", "ruby")
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