<template>
  <div class="flex items-center justify-between w-full">
    <span v-tooltip.left="'Wrap metrics reported as a mean in a ±1 standard deviation band, so the spread between the runs behind the mean is visible'">Std dev interval:</span>
    <ToggleSwitch
      v-model="settingsStore.stdDevInterval"
      class="ml-4"
    />
  </div>
</template>

<script setup lang="ts">
import { useSettingsStore } from "./settingsStore"

const settingsStore = useSettingsStore()
</script>
