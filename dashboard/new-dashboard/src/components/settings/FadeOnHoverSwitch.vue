<template>
  <div class="flex items-center justify-between w-full">
    <span v-tooltip.left="'Fade out other series when hovering over one'">Fade others on hover:</span>
    <ToggleSwitch
      v-model="settingsStore.fadeOnHover"
      class="ml-4"
    />
  </div>
</template>

<script setup lang="ts">
import { useSettingsStore } from "./settingsStore"

const settingsStore = useSettingsStore()
</script>

<style scoped></style>
