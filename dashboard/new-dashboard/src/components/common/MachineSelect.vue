<template>
  <DimensionHierarchicalSelect
    label="Machine"
    :dimension="machineConfigurator"
  >
    <template #icon>
      <ComputerDesktopIcon class="w-4 h-4 text-gray-500" />
    </template>
  </DimensionHierarchicalSelect>
</template>

<script setup lang="ts">
import { DimensionConfigurator } from "../../configurators/DimensionConfigurator"
import DimensionHierarchicalSelect from "../charts/DimensionHierarchicalSelect.vue"

interface Props {
  machineConfigurator: DimensionConfigurator
}

defineProps<Props>()
</script>
