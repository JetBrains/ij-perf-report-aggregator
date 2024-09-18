<template>
  <div class="flex items-center justify-between w-full">
    <span v-tooltip.left="'Apply change detector algorithm. '">Detect Changes:</span>
    <ToggleSwitch
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
