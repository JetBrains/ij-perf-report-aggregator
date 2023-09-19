<template>
  <div
    v-if="sidebarVm != null"
    class="flex items-center justify-between w-full"
  >
    <VTooltip theme="info">
      <span>Sidebar:</span>
      <template #popper> <span class="text-sm">Enable sidebar instead of popup</span> </template>
    </VTooltip>

    <InputSwitch
      v-model="sidebarEnabled"
      class="ml-4"
    />
  </div>
</template>

<script setup lang="ts">
import { useStorage } from "@vueuse/core/index"
import { watch } from "vue"
import { injectOrNull } from "../../shared/injectionKeys"
import { sidebarStartupKey } from "../../shared/keys"

const sidebarVm = injectOrNull(sidebarStartupKey)
const sidebarEnabled = useStorage("sidebarEnabled", true)
if (sidebarVm != null) {
  watch(sidebarEnabled, (value) => {
    if (!value) {
      sidebarVm.close()
    }
  })
}
</script>

<style scoped></style>
