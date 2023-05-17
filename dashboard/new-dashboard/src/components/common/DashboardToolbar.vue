<template>
  <Toolbar class="customToolbar">
    <template #start>
      <TimeRangeSelect
        :ranges="TimeRangeConfigurator.timeRanges"
        :value="timeRangeConfigurator.value.value"
        :on-change="onChangeRange"
      >
        <template #icon>
          <CalendarIcon class="w-4 h-4 text-gray-500" />
        </template>
      </TimeRangeSelect>
      <BranchSelect
        :branch-configurator="branchConfigurator"
        :release-configurator="releaseConfigurator"
        :triggered-by-configurator="triggeredByConfigurator"
      />
      <DimensionHierarchicalSelect
        v-if="machineConfigurator != null"
        label="Machine"
        :dimension="machineConfigurator"
      >
        <template #icon>
          <ComputerDesktopIcon class="w-4 h-4 text-gray-500" />
        </template>
      </DimensionHierarchicalSelect>
    </template>
  </Toolbar>
</template>

<script setup lang="ts">
import DimensionHierarchicalSelect from "shared/src/components/DimensionHierarchicalSelect.vue"
import { BranchConfigurator } from "shared/src/configurators/BranchConfigurator"
import { BuildConfigurator } from "shared/src/configurators/BuildConfigurator"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { ReleaseNightlyConfigurator } from "shared/src/configurators/ReleaseNightlyConfigurator"
import { TimeRange, TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import BranchSelect from "./BranchSelect.vue"
import TimeRangeSelect from "./TimeRangeSelect.vue"

const props = defineProps<{
  timeRangeConfigurator: TimeRangeConfigurator
  branchConfigurator: BranchConfigurator
  releaseConfigurator?: ReleaseNightlyConfigurator
  triggeredByConfigurator: BuildConfigurator
  machineConfigurator?: MachineConfigurator
  onChangeRange: (value: TimeRange) => void
}>()

</script>

<style>
.customToolbar {
  background-color: transparent;
  border: none;
  padding: 0;
}
</style>