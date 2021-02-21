<template>
  <el-form-item label="Server">
    <el-input
      v-model="value"
      data-lpignore="true"
      placeholder="Enter the aggregated stats server URL..."
    />
  </el-form-item>
</template>
<script lang="ts">
import { computed, defineComponent } from "vue"

export default defineComponent({
  name: "ServerSelect",
  props: {
    modelValue: {
      type: String,
      default: "",
    },
  },
  emits: ["update:modelValue"],
  setup(props, {emit}) {
    return {
      value: computed({
        get: () => props.modelValue,
        set: value => emit("update:modelValue", value),
      }),
    }
  },
})
</script>