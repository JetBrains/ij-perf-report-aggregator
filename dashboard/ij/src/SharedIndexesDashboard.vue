<template>
  <div>
    <el-row>
      <el-col :span="18">
        <el-form
          :inline="true"
          size="small"
        >
          <DimensionSelect
            label="Scenarios"
            :dimension="scenarioConfigurator"
          />
          <MeasureSelect :configurator="measureConfigurator" />
          <DimensionHierarchicalSelect
            label="Machine"
            :dimension="machineConfigurator"
          />
          <TimeRangeSelect :configurator="timeRangeConfigurator" />

          <ReloadButton />
        </el-form>
      </el-col>
    </el-row>

    <template
      v-for="metric in measureConfigurator.value.value"
      :key="metric"
    >
      <el-divider>
        {{ metric }}
      </el-divider>
      <el-row
        :gutter="5"
      >
        <el-col :span="14">
          <LineChartCard
            :measures="[metric]"
          />
        </el-col>
        <el-col :span="10">
          <BarChartCard
            :height="chartHeight"
            :measures="[metric]"
          />
        </el-col>
      </el-row>
    </template>
  </div>
</template>

<script lang="ts">
import { DataQueryExecutor, initDataComponent } from "shared/src/DataQueryExecutor"
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import { DEFAULT_LINE_CHART_HEIGHT } from "shared/src/chart"
import BarChartCard from "shared/src/components/BarChartCard.vue"
import DimensionHierarchicalSelect from "shared/src/components/DimensionHierarchicalSelect.vue"
import DimensionSelect from "shared/src/components/DimensionSelect.vue"
import LineChartCard from "shared/src/components/LineChartCard.vue"
import MeasureSelect from "shared/src/components/MeasureSelect.vue"
import ReloadButton from "shared/src/components/ReloadButton.vue"
import TimeRangeSelect from "shared/src/components/TimeRangeSelect.vue"
import { AggregationOperatorConfigurator } from "shared/src/configurators/AggregationOperatorConfigurator"
import { DimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { MeasureConfigurator } from "shared/src/configurators/MeasureConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { aggregationOperatorConfiguratorKey } from "shared/src/injectionKeys"
import { defineComponent, provide } from "vue"
import { useRouter } from "vue-router"

export default defineComponent({
  components: {ReloadButton, DimensionSelect, DimensionHierarchicalSelect, TimeRangeSelect, MeasureSelect, LineChartCard, BarChartCard},
  setup() {
    const persistentStateManager = new PersistentStateManager("sharedIndexes-dashboard", {
      machine: "macMini 2018",
      project: [
        "ijx-intellij-speed/shared-indexes",
        "ijx-intellij-speed/usual-indexes",
        "ijx-intellij-speed/shared-indexes-with-archive-and-git-hashes",
      ],
      measure: ["indexing", "scanning"],
    }, useRouter())

    const serverConfigurator = new ServerConfigurator("sharedIndexes")
    const scenarioConfigurator = new DimensionConfigurator("project", serverConfigurator, persistentStateManager, true)

    const machineConfigurator = new MachineConfigurator(new DimensionConfigurator("machine", serverConfigurator, persistentStateManager),
      persistentStateManager)

    const measureConfigurator = new MeasureConfigurator(serverConfigurator, persistentStateManager, scenarioConfigurator)
    const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)

    // median by default, no UI control to change is added (insert <AggregationOperatorSelect /> if needed)
    provide(aggregationOperatorConfiguratorKey, new AggregationOperatorConfigurator(persistentStateManager))

    const dataQueryExecutor = new DataQueryExecutor([
      serverConfigurator,
      scenarioConfigurator,
      machineConfigurator,
      timeRangeConfigurator,
    ], true)

    initDataComponent(persistentStateManager, dataQueryExecutor)
    return {
      chartHeight: DEFAULT_LINE_CHART_HEIGHT,
      scenarioConfigurator,
      machineConfigurator,
      measureConfigurator,
      timeRangeConfigurator,
      loadData: dataQueryExecutor.scheduleLoadIncludingConfiguratorsFunctionReference,
      valueToGroup(value: string): string {
        const separatorIndex = value.indexOf("/")
        return separatorIndex > 0 ? value.substring(0, separatorIndex) : value
      }
    }
  },
})
</script>