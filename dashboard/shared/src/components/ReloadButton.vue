<template>
  <Button
    title="Updated automatically, but you can force data reloading"
    icon="pi pi-refresh"
    class="p-button-rounded p-button-text"
    @click="doLoad"
  />
</template>
<script setup lang="ts">
import { computed, inject, PropType } from "vue"
import { dataQueryExecutorKey } from "../injectionKeys"

const props = defineProps({
  load: {
    type: Function as PropType<() => void>,
    default: null,
  },
})

const doLoad = computed(() => {
  const load = props.load
  if (load != null) {
    return load
  }

  const executor = inject(dataQueryExecutorKey)
  if (executor == null) {
    throw new Error("Neither `load` function is set, nor `dataQueryExecutor` is provided")
  }
  return executor.scheduleLoadIncludingConfiguratorsFunctionReference
})
</script>