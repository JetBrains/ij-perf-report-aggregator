<template>
  <div class="flex items-center justify-between w-full">
    <span v-tooltip.left="'Make zero on Y axis flexible and not always equals to 0.'">Flexible zero on Y axis:</span>
    <ToggleSwitch
      v-model="settingsStore.flexibleYZero"
      class="ml-4"
    />
  </div>
</template>

<script setup lang="ts">
import { useSettingsStore } from "./settingsStore"

const settingsStore = useSettingsStore()
</script>

<style scoped></style>
