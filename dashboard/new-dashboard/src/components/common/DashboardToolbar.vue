<template>
  <Toolbar class="customToolbar">
    <template #start>
      <TimeRangeSelect
        :ranges="TimeRangeConfigurator.timeRanges"
        :value="props.timeRangeConfigurator.value.value"
        :on-change="onChangeRange"
      />
      <BranchSelect
        :branch-configurator="props.branchConfigurator"
        :release-configurator="props.releaseConfigurator"
        :triggered-by-configurator="props.triggeredByConfigurator"
      />
      <MachineSelect
        v-if="machineConfigurator != null"
        :machine-configurator="machineConfigurator"
      />
    </template>
    <template #end>
      Smoothing:
      <InputSwitch v-model="smoothingEnabled" />
      Scaling:
      <InputSwitch v-model="scalingEnabled" />
    </template>
  </Toolbar>
</template>

<script setup lang="ts">
import { useStorage } from "@vueuse/core"
import { BranchConfigurator } from "../../configurators/BranchConfigurator"
import { BuildConfigurator } from "../../configurators/BuildConfigurator"
import { MachineConfigurator } from "../../configurators/MachineConfigurator"
import { ReleaseNightlyConfigurator } from "../../configurators/ReleaseNightlyConfigurator"
import { TimeRange, TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"
import BranchSelect from "./BranchSelect.vue"
import MachineSelect from "./MachineSelect.vue"
import TimeRangeSelect from "./TimeRangeSelect.vue"

const props = defineProps<{
  timeRangeConfigurator: TimeRangeConfigurator
  branchConfigurator: BranchConfigurator
  releaseConfigurator?: ReleaseNightlyConfigurator
  triggeredByConfigurator: BuildConfigurator
  machineConfigurator?: MachineConfigurator
  onChangeRange: (value: TimeRange) => void
}>()

const smoothingEnabled = useStorage("smoothingEnabled", true)
const scalingEnabled = useStorage("scalingEnabled", true)
</script>

<style>
.customToolbar {
  background-color: transparent;
  border: none;
  padding: 0;
}
</style>
