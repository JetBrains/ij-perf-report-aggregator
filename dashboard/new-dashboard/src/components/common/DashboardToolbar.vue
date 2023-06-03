<template>
  <Toolbar class="customToolbar">
    <template #start>
      <TimeRangeSelect
        :ranges="TimeRangeConfigurator.timeRanges"
        :value="props.timeRangeConfigurator.value.value"
        :on-change="onChangeRange"
      >
        <template #icon>
          <CalendarIcon class="w-4 h-4 text-gray-500" />
        </template>
      </TimeRangeSelect>
      <BranchSelect
        :branch-configurator="props.branchConfigurator"
        :release-configurator="props.releaseConfigurator"
        :triggered-by-configurator="props.triggeredByConfigurator"
      />
      <DimensionHierarchicalSelect
        v-if="props.machineConfigurator != null"
        label="Machine"
        :dimension="props.machineConfigurator"
      >
        <template #icon>
          <ComputerDesktopIcon class="w-4 h-4 text-gray-500" />
        </template>
      </DimensionHierarchicalSelect>
    </template>
  </Toolbar>
</template>

<script setup lang="ts">
import { BranchConfigurator } from "../../configurators/BranchConfigurator"
import { BuildConfigurator } from "../../configurators/BuildConfigurator"
import { MachineConfigurator } from "../../configurators/MachineConfigurator"
import { ReleaseNightlyConfigurator } from "../../configurators/ReleaseNightlyConfigurator"
import { TimeRange, TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"
import DimensionHierarchicalSelect from "../charts/DimensionHierarchicalSelect.vue"
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