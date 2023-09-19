<template>
  <VTooltip theme="info">
    <span>Scaling:</span>
    <template #popper> <span class="text-sm">Scales each value in the dataset on its median. Each element is divided by the median and multiplied by 50.</span> </template>
  </VTooltip>

  <InputSwitch v-model="settingsStore.scaling" />
</template>

<script setup lang="ts">
import { ScalingConfigurator } from "../../configurators/ScalingConfigurator"
import { useSettingsStore } from "./settingsStore"

const settingsStore = useSettingsStore()

const emit = defineEmits(["update:configurators"])
emit("update:configurators", new ScalingConfigurator())
</script>

<style scoped></style>
