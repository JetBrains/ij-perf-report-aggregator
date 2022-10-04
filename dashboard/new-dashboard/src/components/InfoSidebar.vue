<template>
  <div
    v-if="visible.value"
    class="min-w-80 p-5 ml-5 border border-solid rounded-md border-zinc-200"
  >
    <span class="uppercase text-xs mb-4">
      Named entities adding (2000000)
    </span>
  </div>
</template>
<script setup lang="ts">
import { inject, onMounted, ref } from "vue"
import { sidebarVmKey } from "../shared/keys"
import { InfoSidebarVmImpl } from "./InfoSidebarVm"

const vm = inject(sidebarVmKey) || new InfoSidebarVmImpl()
const visible = ref(false)

onMounted(() => {
  setTimeout(() => {
    console.log("kek")
    // vm.toggle()
    // visible.value = true
  }, 5000)
})
</script>