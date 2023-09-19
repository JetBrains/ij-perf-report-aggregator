<template>
  <VTooltip theme="info">
    <span>Smoothing:</span>
    <template #popper><span class="text-sm">Applies exponential smoothing to the dataset.</span></template>
  </VTooltip>
  <InputSwitch v-model="settingsStore.smoothing" />
</template>

<script setup lang="ts">
import { SmoothingConfigurator } from "../../configurators/SmoothingConfigurator"
import { useSettingsStore } from "./settingsStore"

const settingsStore = useSettingsStore()

const emit = defineEmits(["update:configurators"])
emit("update:configurators", new SmoothingConfigurator())
</script>

<style scoped></style>
