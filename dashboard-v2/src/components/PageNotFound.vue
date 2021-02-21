<template>
  <div>Page not found.</div>
</template>
<script lang="ts">
import { defineComponent, onMounted, onUnmounted } from "vue"
import { ElMessage } from "element-plus"
import { IMessageHandle } from "element-plus/lib/el-message/src/types"

export default defineComponent({
  name: "PageNotFound",
  setup() {
    let messageHandle: IMessageHandle | null = null
    onMounted(() => {
      messageHandle = ElMessage({
        message: "Page Not Found",
        type: "error",
        duration: 0,
        showClose: false,
      })
    })

    onUnmounted(() => {
      messageHandle?.close()
    })
    return {}
  }
})
</script>