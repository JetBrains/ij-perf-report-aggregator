<template>
  <VTooltip theme="info">
    <span>Scaling:</span>
    <template #popper> <span class="text-sm">Scales each value in the dataset on its median. Each element is divided by the median and multiplied by 50.</span> </template>
  </VTooltip>

  <InputSwitch v-model="scalingEnabled" />
</template>

<script setup lang="ts">
import { useStorage } from "@vueuse/core/index"
import { ScalingConfigurator } from "../../configurators/ScalingConfigurator"

const scalingEnabled = useStorage("scalingEnabled", true)

const emit = defineEmits(["update:configurators"])
emit("update:configurators", new ScalingConfigurator())
</script>

<style scoped></style>
