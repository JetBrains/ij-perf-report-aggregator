<template>
  <Toolbar>
    <template #start>
      <slot name="toolbar" />
    </template>
    <template #end>
      <ReloadButton />
    </template>
  </Toolbar>
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
