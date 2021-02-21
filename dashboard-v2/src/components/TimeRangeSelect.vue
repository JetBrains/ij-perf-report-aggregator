<template>
  <el-form-item label="Period">
    <el-select
      v-model="value"
      filterable
    >
      <el-option
        v-for="item in timeRanges"
        :key="item.value"
        :label="item.label"
        :value="item.value"
      />
    </el-select>
  </el-form-item>
</template>
<script lang="ts">
import { defineComponent, ref } from "vue"
import { TimeRangeConfigurator } from "../configurators/TimeRangeConfigurator"

export default defineComponent({
  name: "TimeRangeSelect",
  props: {
    configurator: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    return {
      value: (props.configurator as TimeRangeConfigurator).value,
      timeRanges: ref(TimeRangeConfigurator.timeRanges),
    }
  },
})
</script>