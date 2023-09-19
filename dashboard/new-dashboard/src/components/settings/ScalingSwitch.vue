<template>
  <div class="flex items-center justify-between w-full">
    <VTooltip theme="info">
      <span>Scaling:</span>
      <template #popper> <span class="text-sm">Scales each value in the dataset on its median. Each element is divided by the median and multiplied by 50.</span> </template>
    </VTooltip>

    <InputSwitch
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
