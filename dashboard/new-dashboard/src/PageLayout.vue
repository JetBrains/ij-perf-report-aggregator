<template>
  <PageHeader />
  <slot />

  <ChartTooltip
    ref="tooltip"
  />
</template>
<script setup lang="ts">
import ChartTooltip from "shared/src/components/ChartTooltip.vue"
import { chartToolTipKey } from "shared/src/injectionKeys"
import { provide, ref } from "vue"
import PageHeader from "./PageHeader.vue"

const tooltip = ref<typeof ChartTooltip>()

provide(chartToolTipKey, tooltip)
</script>