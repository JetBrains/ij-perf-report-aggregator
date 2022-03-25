<template>
  <Dropdown
    v-model="value"
    :options="timeRanges"
    option-label="label"
    option-value="value"
  />
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