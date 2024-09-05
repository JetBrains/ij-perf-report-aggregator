<template>
  <StickyToolbar>
    <template #start>
      <TimeRangeSelect :timerange-configurator="timeRangeConfigurator" />
      <BranchSelect
        v-if="branchConfigurator != null"
        :branch-configurator="branchConfigurator"
        :release-configurator="releaseConfigurator"
        :triggered-by-configurator="triggeredByConfigurator"
      />
      <MachineSelect
        v-if="machineConfigurator != null"
        :machine-configurator="machineConfigurator"
      />
      <slot name="configurator" />
      <CopyLink :timerange-configurator="timeRangeConfigurator" />
    </template>
    <template #end>
      <slot name="toolbar" />
    </template>
  </StickyToolbar>
</template>

<script setup lang="ts">
import { BranchConfigurator } from "../../configurators/BranchConfigurator"
import { BuildConfigurator } from "../../configurators/BuildConfigurator"
import { MachineConfigurator } from "../../configurators/MachineConfigurator"
import { ReleaseNightlyConfigurator } from "../../configurators/ReleaseNightlyConfigurator"
import { TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"
import CopyLink from "../settings/CopyLink.vue"
import BranchSelect from "./BranchSelect.vue"
import MachineSelect from "./MachineSelect.vue"
import StickyToolbar from "./StickyToolbar.vue"
import TimeRangeSelect from "./TimeRangeSelect.vue"

const { timeRangeConfigurator, branchConfigurator, releaseConfigurator, triggeredByConfigurator, machineConfigurator } = defineProps<{
  timeRangeConfigurator: TimeRangeConfigurator
  branchConfigurator: BranchConfigurator | null
  releaseConfigurator?: ReleaseNightlyConfigurator
  triggeredByConfigurator: BuildConfigurator
  machineConfigurator?: MachineConfigurator
}>()
</script>
