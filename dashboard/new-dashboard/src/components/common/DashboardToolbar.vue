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
    </template>
  </Toolbar>
</template>

<script setup lang="ts">
import { useStorage } from "@vueuse/core"
import { watch } from "vue"
import { BranchConfigurator } from "../../configurators/BranchConfigurator"
import { BuildConfigurator } from "../../configurators/BuildConfigurator"
import { MachineConfigurator } from "../../configurators/MachineConfigurator"
import { ReleaseNightlyConfigurator } from "../../configurators/ReleaseNightlyConfigurator"
import { TimeRange, TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"
import { useSmoothingStore } from "../../shared/storage"
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
watch(smoothingEnabled, (value) => {
  window.location.reload()
  useSmoothingStore().isSmoothingEnabled = value
})
useSmoothingStore().isSmoothingEnabled = smoothingEnabled.value
</script>

<style>
.customToolbar {
  background-color: transparent;
  border: none;
  padding: 0;
}
</style>
