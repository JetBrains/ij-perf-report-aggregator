<template>
  Scaling:
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
