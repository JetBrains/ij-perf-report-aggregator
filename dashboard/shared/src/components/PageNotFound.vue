<template>
  <div>Page not found.</div>
</template>
<script lang="ts">
// import "element-plus/theme-chalk/el-message.css"

import { defineComponent, onMounted, onUnmounted } from "vue"

export default defineComponent({
  name: "PageNotFound",
  setup() {
    let messageHandle: MessageHandle | null = null
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