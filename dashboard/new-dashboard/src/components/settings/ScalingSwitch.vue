<template>
  <div class="flex items-center justify-between w-full">
    <span v-tooltip.left="'Scales each value in the dataset on its median. Each element is divided by the median and multiplied by 50.'">Scaling:</span>
    <ToggleSwitch
      v-model="settingsStore.scaling"
      class="ml-4"
    />
  </div>
</template>

<script setup lang="ts">
import { useSettingsStore } from "./settingsStore"

const settingsStore = useSettingsStore()
</script>

<style scoped></style>
