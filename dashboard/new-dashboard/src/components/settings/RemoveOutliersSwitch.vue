<template>
  <div class="flex items-center justify-between w-full">
    <span v-tooltip.left="'Remove outliers based on rolling MAD (Median Value Deviation) score'">Remove outliers (Alpha):</span>
    <InputSwitch
      v-model="settingsStore.removeOutliers"
      class="ml-4"
    />
  </div>
</template>

<script setup lang="ts">
import { useSettingsStore } from "./settingsStore"

const settingsStore = useSettingsStore()
</script>

<style scoped></style>
