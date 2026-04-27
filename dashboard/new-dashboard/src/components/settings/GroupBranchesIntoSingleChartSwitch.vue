<template>
  <div class="flex items-center justify-between w-full">
    <span v-tooltip.left="'Merge selected branches into one series for each metric'">Merge branches:</span>
    <ToggleSwitch
      v-model="settingsStore.groupBranchesIntoSingleChart"
      class="ml-4"
    />
  </div>
</template>

<script setup lang="ts">
import { useSettingsStore } from "./settingsStore"

const settingsStore = useSettingsStore()
</script>

<style scoped></style>
