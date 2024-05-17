<template>
  <DimensionHierarchicalSelect
    label="Machine"
    :dimension="machineConfigurator"
  />
</template>

<script setup lang="ts">
import { MachineConfigurator } from "../../configurators/MachineConfigurator"
import DimensionHierarchicalSelect from "../charts/DimensionHierarchicalSelect.vue"

interface Props {
  machineConfigurator: MachineConfigurator
}

defineProps<Props>()
</script>
