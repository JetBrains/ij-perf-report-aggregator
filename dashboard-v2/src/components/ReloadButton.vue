<template>
  <el-form-item>
    <el-button
      title="Updated automatically, but you can force data reloading"
      type="primary"
      icon="el-icon-refresh"
      @click="load"
    />
  </el-form-item>
</template>
<script lang="ts">
import { defineComponent } from "vue"
export default defineComponent({
  name: "ReloadButton",
  props: {
    load: {
      type: Function,
      required: true,
    },
  },
})
</script>