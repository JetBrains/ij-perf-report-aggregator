<template>
  <ElFormItem>
    <ElButton
      title="Updated automatically, but you can force data reloading"
      type="primary"
      icon="el-icon-refresh"
      @click="doLoad"
    />
  </ElFormItem>
</template>
<script lang="ts">
import { defineComponent, inject, PropType, computed } from "vue"
import { dataQueryExecutorKey } from "../injectionKeys"

export default defineComponent({
  name: "ReloadButton",
  props: {
    load: {
      type: Function as PropType<() => void>,
      default: null,
    },
  },
  setup(props) {
    return {
      doLoad: computed(() => {
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
    }
  }
})
</script>