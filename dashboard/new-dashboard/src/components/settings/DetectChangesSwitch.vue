<template>
  <div class="flex items-center justify-between w-full">
    <VTooltip theme="info">
      <span>Detect Changes (Beta):</span>
      <template #popper><span class="text-sm">Apply change detector algorithm. </span></template>
    </VTooltip>
    <InputSwitch
      v-model="settingsStore.detectChanges"
      class="ml-4"
    />
  </div>
</template>

<script setup lang="ts">
import { useSettingsStore } from "./settingsStore"

const settingsStore = useSettingsStore()
</script>

<style scoped></style>
