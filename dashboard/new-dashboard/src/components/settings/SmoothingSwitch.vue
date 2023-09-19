<template>
  <div class="flex items-center justify-between w-full">
    <VTooltip theme="info">
      <span>Smoothing:</span>
      <template #popper><span class="text-sm">Applies exponential smoothing to the dataset.</span></template>
    </VTooltip>
    <InputSwitch
      v-model="settingsStore.smoothing"
      class="ml-4"
    />
  </div>
</template>

<script setup lang="ts">
import { useSettingsStore } from "./settingsStore"

const settingsStore = useSettingsStore()
</script>

<style scoped></style>
