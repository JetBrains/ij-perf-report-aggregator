<template>
  <ElFormItem label="Period">
    <ElSelect
      v-model="value"
      filterable
    >
      <ElOption
        v-for="item in timeRanges"
        :key="item.value"
        :label="item.label"
        :value="item.value"
      />
    </ElSelect>
  </ElFormItem>
</template>
<script lang="ts">
import { defineComponent, ref } from "vue"
import { TimeRangeConfigurator } from "../configurators/TimeRangeConfigurator"

export default defineComponent({
  name: "TimeRangeSelect",
  props: {
    configurator: {
      type: TimeRangeConfigurator,
      required: true,
    },
  },
  setup(props) {
    return {
      value: props.configurator.value,
      timeRanges: ref(TimeRangeConfigurator.timeRanges),
    }
  },
})
</script>