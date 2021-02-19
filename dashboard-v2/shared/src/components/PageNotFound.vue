<template>
  <div>Page not found.</div>
</template>
<script lang="ts">
import "element-plus/packages/theme-chalk/src/message.scss"

import ElMessage from "element-plus/es/el-message"
import { IMessageHandle } from "element-plus/lib/el-message/src/types"
import { defineComponent, onMounted, onUnmounted } from "vue"

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