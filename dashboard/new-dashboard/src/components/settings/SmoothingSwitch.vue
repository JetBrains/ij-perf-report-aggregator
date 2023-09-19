<template>
  <VTooltip theme="info">
    <span>Smoothing:</span>
    <template #popper><span class="text-sm">Applies exponential smoothing to the dataset.</span></template>
  </VTooltip>
  <InputSwitch v-model="smoothingEnabled" />
</template>

<script setup lang="ts">
import { useStorage } from "@vueuse/core/index"
import { SmoothingConfigurator } from "../../configurators/SmoothingConfigurator"

const smoothingEnabled = useStorage("smoothingEnabled", true)

const emit = defineEmits(["update:configurators"])
emit("update:configurators", new SmoothingConfigurator())
</script>

<style scoped></style>
