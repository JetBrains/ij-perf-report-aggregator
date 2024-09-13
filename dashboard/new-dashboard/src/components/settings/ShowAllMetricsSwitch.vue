<template>
  <div class="flex items-center justify-between w-full">
    <span v-tooltip.left="'Disable filtering of detailed metrics like completion_32.'">Show all metrics:</span>
    <InputSwitch
      v-model="settingsStore.showAllMetrics"
      class="ml-4"
    />
  </div>
</template>

<script setup lang="ts">
import { useSettingsStore } from "./settingsStore"

const settingsStore = useSettingsStore()
</script>

<style scoped></style>
