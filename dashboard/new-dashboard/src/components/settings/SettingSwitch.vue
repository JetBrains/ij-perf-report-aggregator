<template>
  <div class="flex items-center justify-between w-full">
    <span v-tooltip.left="tooltip">{{ label }}:</span>
    <ToggleSwitch
      v-model="settingsStore[setting]"
      class="ml-4"
    />
  </div>
</template>

<script setup lang="ts">
import { BooleanSetting } from "./SettingConfigurator"
import { useSettingsStore } from "./settingsStore"

defineProps<{
  label: string
  tooltip: string
  setting: BooleanSetting
}>()

const settingsStore = useSettingsStore()
</script>

<style scoped></style>
