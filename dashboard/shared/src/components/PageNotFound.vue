<template>
  <div>Page not found.</div>
</template>
<script lang="ts">
// import "element-plus/theme-chalk/el-message.css"

import { ElMessage } from "element-plus/es/components/message"
import { IMessageHandle } from "element-plus/es/components/message/src/types"
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