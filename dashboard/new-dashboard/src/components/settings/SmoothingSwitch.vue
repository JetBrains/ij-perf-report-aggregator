<template>
  <div class="flex items-center justify-between w-full">
    <span v-tooltip.left="'Applies exponential smoothing to the dataset.'">Smoothing:</span>
    <ToggleSwitch
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
