<template>
  <div class="flex items-center justify-between w-full">
    <span v-tooltip.left="'Remove values that are larger than 95 and lower than 5 percentile'">Remove outliers (Alpha):</span>
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
