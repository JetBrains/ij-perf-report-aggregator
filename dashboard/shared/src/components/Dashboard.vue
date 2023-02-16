<template>
  <div class="flex items-center justify-between">
    <div class="flex items-baseline gap-x-2">
      <slot name="toolbar" />
    </div>
    <ReloadButton />
  </div>
  <slot />
  <ChartTooltip
    ref="tooltip"
  />
</template>

<script setup lang="ts">
import { provide, ref } from "vue"
import { chartToolTipKey } from "../injectionKeys"
import ChartTooltip from "./ChartTooltip.vue"
import ReloadButton from "./ReloadButton.vue"

const tooltip = ref<typeof ChartTooltip>()
provide(chartToolTipKey, tooltip)
</script>
